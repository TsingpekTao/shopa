package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/TsingpekTao/shopa/iam-svc/internal/dao"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// StartBackgroundWorkers 启动后台任务或服务流程。
func (s *Service) StartBackgroundWorkers(ctx context.Context) {
	s.workerOnce.Do(func() {
		go s.dispatchLoop()
		go s.archiveLoop()
		go s.cleanupLoop()
		go s.metricsLoop()
	})
}

// dispatchLoop 周期性扫描热表并派发可投递事件。
func (s *Service) dispatchLoop() {
	ticker := time.NewTicker(s.outbox.PollInterval)
	defer ticker.Stop()

	ctx := context.Background()
	for {
		if err := s.dispatchOnce(ctx); err != nil {
			g.Log().Errorf(ctx, "[iam-svc] outbox dispatch tick failed: %+v", err)
		}
		<-ticker.C
	}
}

// archiveLoop 归档历史数据，控制热表规模。
func (s *Service) archiveLoop() {
	ticker := time.NewTicker(s.outbox.ArchiveInterval)
	defer ticker.Stop()

	ctx := context.Background()
	for {
		if err := s.archiveOnce(ctx); err != nil {
			g.Log().Errorf(ctx, "[iam-svc] outbox archive tick failed: %+v", err)
		}
		<-ticker.C
	}
}

// cleanupLoop 周期性清理归档表中过期数据，控制磁盘占用。
func (s *Service) cleanupLoop() {
	ticker := time.NewTicker(s.outbox.CleanupInterval)
	defer ticker.Stop()

	ctx := context.Background()
	for {
		if err := s.cleanupArchiveOnce(ctx); err != nil {
			g.Log().Errorf(ctx, "[iam-svc] outbox cleanup tick failed: %+v", err)
		}
		<-ticker.C
	}
}

// metricsLoop 周期性输出 outbox 指标，便于观察热表规模和堆积时长。
func (s *Service) metricsLoop() {
	ticker := time.NewTicker(s.outbox.MetricsInterval)
	defer ticker.Stop()

	ctx := context.Background()
	for {
		s.logOutboxMetrics(ctx)
		<-ticker.C
	}
}

// dispatchOnce 执行一次批量认领与派发。
func (s *Service) dispatchOnce(ctx context.Context) error {
	items, err := s.claimOutboxBatch(ctx, s.outbox.BatchSize)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	for _, item := range items {
		if err = s.dispatchOne(ctx, item); err != nil {
			g.Log().Warningf(ctx, "[iam-svc] outbox dispatch event failed id=%d eventId=%s err=%+v", item.Id, item.EventId, err)
		}
	}
	return nil
}

// claimOutboxBatch 使用 SKIP LOCKED 认领一批事件，避免多副本重复消费。
func (s *Service) claimOutboxBatch(ctx context.Context, batchSize int) ([]entity.IamOutboxEvent, error) {
	items := make([]entity.IamOutboxEvent, 0, batchSize)
	if batchSize <= 0 {
		return items, nil
	}
	// ids 是本次被当前 worker 认领的行集合，后续统一改为 PROCESSING。
	ids := make([]any, 0, batchSize)

	err := dao.IamOutboxEvent.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		sql := fmt.Sprintf(
			"SELECT id,event_id,event_type,payload_json,status,available_at,sent_at,fail_count,last_error,retry_count,next_retry_at,created_at,updated_at "+
				"FROM %s "+
				"WHERE status IN (%d,%d) AND available_at<=NOW(3) "+
				"ORDER BY id LIMIT %d FOR UPDATE SKIP LOCKED",
			dao.IamOutboxEvent.Table(),
			consts.OutboxStatusNew,
			consts.OutboxStatusFailed,
			batchSize,
		)
		records, queryErr := tx.GetAll(sql)
		if queryErr != nil {
			return queryErr
		}
		if len(records) == 0 {
			return nil
		}

		for _, record := range records {
			var item entity.IamOutboxEvent
			if structErr := record.Struct(&item); structErr != nil {
				return structErr
			}
			items = append(items, item)
			ids = append(ids, item.Id)
		}

		cols := dao.IamOutboxEvent.Columns()
		_, updateErr := tx.Model(dao.IamOutboxEvent.Table()).
			WhereIn(cols.Id, ids).
			WhereIn(cols.Status, []int{consts.OutboxStatusNew, consts.OutboxStatusFailed}).
			Data(do.IamOutboxEvent{
				Status:    consts.OutboxStatusProcessing,
				LastError: "",
			}).
			Update()
		return updateErr
	})
	return items, err
}

// dispatchOne 派发单条事件，成功写 SENT，失败写 FAILED 或 DLQ 并设置退避时间。
func (s *Service) dispatchOne(ctx context.Context, item entity.IamOutboxEvent) error {
	var err error
	if item.EventType == "UserRegisteredV1" {
		err = s.mq.PublishUserRegistered(ctx, item.EventId, []byte(item.PayloadJson))
	} else {
		err = fmt.Errorf("unsupported event_type: %s", item.EventType)
	}

	cols := dao.IamOutboxEvent.Columns()
	if err == nil {
		_, updateErr := dao.IamOutboxEvent.Ctx(ctx).
			Where(cols.Id, item.Id).
			Where(cols.Status, consts.OutboxStatusProcessing).
			Data(do.IamOutboxEvent{
				Status:    consts.OutboxStatusSent,
				SentAt:    gtime.Now(),
				LastError: "",
			}).
			Update()
		return updateErr
	}

	nextFailCount := int(item.FailCount) + 1
	nextStatus := consts.OutboxStatusFailed
	if nextFailCount >= s.outbox.MaxFailCount {
		nextStatus = consts.OutboxStatusDLQ
	}
	backoff := s.computeBackoff(nextFailCount)
	lastError := truncateError(err, 480)

	_, updateErr := dao.IamOutboxEvent.Ctx(ctx).
		Where(cols.Id, item.Id).
		Where(cols.Status, consts.OutboxStatusProcessing).
		Data(do.IamOutboxEvent{
			Status:      nextStatus,
			FailCount:   nextFailCount,
			LastError:   lastError,
			AvailableAt: gtime.NewFromTime(time.Now().UTC().Add(backoff)),
		}).
		Update()
	if updateErr != nil {
		return updateErr
	}
	return err
}

// computeBackoff 计算策略参数并返回结果。
func (s *Service) computeBackoff(failCount int) time.Duration {
	if failCount <= 0 {
		return s.outbox.FailBaseDelay
	}
	backoff := s.outbox.FailBaseDelay
	for i := 1; i < failCount; i++ {
		backoff = backoff * 2
		if backoff >= s.outbox.FailMaxDelay {
			return s.outbox.FailMaxDelay
		}
	}
	if backoff > s.outbox.FailMaxDelay {
		return s.outbox.FailMaxDelay
	}
	return backoff
}

// archiveOnce 归档历史数据，控制热表规模。
func (s *Service) archiveOnce(ctx context.Context) error {
	if s.outbox.ArchiveBatch <= 0 {
		return nil
	}
	cutoff := time.Now().UTC().Add(-s.outbox.ArchiveAfter)
	return dao.IamOutboxEvent.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		sql := fmt.Sprintf(
			"SELECT id,event_id,event_type,payload_json,status,available_at,sent_at,fail_count,last_error,retry_count,next_retry_at,created_at,updated_at "+
				"FROM %s WHERE status=%d AND sent_at<? ORDER BY id LIMIT %d FOR UPDATE SKIP LOCKED",
			dao.IamOutboxEvent.Table(),
			consts.OutboxStatusSent,
			s.outbox.ArchiveBatch,
		)
		records, err := tx.GetAll(sql, gtime.NewFromTime(cutoff))
		if err != nil {
			return err
		}
		if len(records) == 0 {
			return nil
		}

		ids := make([]any, 0, len(records))
		for _, record := range records {
			var item entity.IamOutboxEvent
			if err = record.Struct(&item); err != nil {
				return err
			}
			ids = append(ids, item.Id)

			_, err = tx.Model(dao.IamOutboxEventArchive.Table()).Data(do.IamOutboxEventArchive{
				Id:          item.Id,
				EventId:     item.EventId,
				EventType:   item.EventType,
				PayloadJson: item.PayloadJson,
				Status:      item.Status,
				AvailableAt: item.AvailableAt,
				SentAt:      item.SentAt,
				FailCount:   item.FailCount,
				LastError:   item.LastError,
				RetryCount:  item.RetryCount,
				NextRetryAt: item.NextRetryAt,
				CreatedAt:   item.CreatedAt,
				UpdatedAt:   item.UpdatedAt,
				ArchivedAt:  gtime.Now(),
			}).Insert()
			if err != nil {
				return err
			}
		}

		cols := dao.IamOutboxEvent.Columns()
		_, err = tx.Model(dao.IamOutboxEvent.Table()).
			WhereIn(cols.Id, ids).
			Delete()
		return err
	})
}

// cleanupArchiveOnce 删除归档表中超过保留期的数据。
func (s *Service) cleanupArchiveOnce(ctx context.Context) error {
	if s.outbox.CleanupBatch <= 0 {
		return nil
	}
	cutoff := time.Now().UTC().Add(-s.outbox.CleanupAfter)
	sql := fmt.Sprintf(
		"DELETE FROM %s WHERE archived_at<? ORDER BY id LIMIT %d",
		dao.IamOutboxEventArchive.Table(),
		s.outbox.CleanupBatch,
	)
	_, err := dao.IamOutboxEventArchive.DB().Exec(ctx, sql, gtime.NewFromTime(cutoff))
	return err
}

// logOutboxMetrics 采集并打印 outbox 关键容量指标。
func (s *Service) logOutboxMetrics(ctx context.Context) {
	hotRows, err := dao.IamOutboxEvent.Ctx(ctx).Count()
	if err != nil {
		g.Log().Warningf(ctx, "[iam-svc] outbox metrics hot count failed: %+v", err)
		return
	}
	archiveRows, err := dao.IamOutboxEventArchive.Ctx(ctx).Count()
	if err != nil {
		g.Log().Warningf(ctx, "[iam-svc] outbox metrics archive count failed: %+v", err)
		return
	}

	unsentAge := int64(0)
	ageSQL := fmt.Sprintf(
		"SELECT TIMESTAMPDIFF(SECOND, MIN(created_at), NOW(3)) AS age_sec FROM %s WHERE status IN (%d,%d,%d)",
		dao.IamOutboxEvent.Table(),
		consts.OutboxStatusNew,
		consts.OutboxStatusProcessing,
		consts.OutboxStatusFailed,
	)
	if value, qErr := dao.IamOutboxEvent.DB().GetValue(ctx, ageSQL); qErr == nil && value != nil && !value.IsNil() {
		unsentAge = value.Int64()
	}

	hotDiskMB := float64(0)
	archiveDiskMB := float64(0)
	sizeSQL := "SELECT ROUND((DATA_LENGTH+INDEX_LENGTH)/1024/1024,2) AS mb FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=?"
	if value, qErr := dao.IamOutboxEvent.DB().GetValue(ctx, sizeSQL, dao.IamOutboxEvent.Table()); qErr == nil && value != nil && !value.IsNil() {
		hotDiskMB = value.Float64()
	}
	if value, qErr := dao.IamOutboxEventArchive.DB().GetValue(ctx, sizeSQL, dao.IamOutboxEventArchive.Table()); qErr == nil && value != nil && !value.IsNil() {
		archiveDiskMB = value.Float64()
	}

	g.Log().Infof(
		ctx,
		"[iam-svc] outbox metrics hot_rows=%d archive_rows=%d oldest_unsent_age_seconds=%d hot_disk_mb=%.2f archive_disk_mb=%.2f",
		hotRows,
		archiveRows,
		unsentAge,
		hotDiskMB,
		archiveDiskMB,
	)
}

// truncateError 截断错误文本，防止超长错误撑爆字段长度。
func truncateError(err error, maxLen int) string {
	if err == nil {
		return ""
	}
	msg := strings.TrimSpace(err.Error())
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen]
}
