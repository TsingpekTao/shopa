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

// StartBackgroundWorkers 启动 IAM 后台 worker。
// 幂等保证：通过 workerOnce 确保即使被重复调用，也只会启动一组 goroutine。
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

// archiveLoop 周期性归档已发送事件，控制热表体量。
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

// metricsLoop 周期性打印 outbox 指标，便于观测积压与容量变化。
func (s *Service) metricsLoop() {
	ticker := time.NewTicker(s.outbox.MetricsInterval)
	defer ticker.Stop()

	ctx := context.Background()
	for {
		s.logOutboxMetrics(ctx)
		<-ticker.C
	}
}

// dispatchOnce 执行一轮“认领 + 投递”。
// 认领失败会整体返回错误；单条投递失败仅记录日志并继续后续条目。
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

// claimOutboxBatch 使用 FOR UPDATE SKIP LOCKED 认领一批事件。
// 并发安全：
// 1) 行锁 + SKIP LOCKED 保证多实例下同一行只会被一个 worker 领走。
// 2) 在同一事务内把状态改为 PROCESSING，避免“查到但未标记”导致重复消费。
func (s *Service) claimOutboxBatch(ctx context.Context, batchSize int) ([]entity.IamOutboxEvent, error) {
	items := make([]entity.IamOutboxEvent, 0, batchSize)
	if batchSize <= 0 {
		return items, nil
	}
	// ids 记录本次事务中已认领的主键，用于批量更新状态。
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

// dispatchOne 派发单条事件。
// 状态机：
// 1) 投递成功：PROCESSING -> SENT，并记录 sent_at。
// 2) 投递失败：PROCESSING -> FAILED/DLQ，并写 fail_count、last_error、下一次 available_at。
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

// computeBackoff 按失败次数计算指数退避，并封顶到最大延迟。
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

// archiveOnce 把已发送且超过归档阈值的数据搬迁到归档表。
// 原子性：在单事务里完成“插入归档表 + 删除热表”，避免数据重复或丢失。
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
// 采用 LIMIT 分批删除，降低长事务和锁冲突风险。
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

// logOutboxMetrics 采集并打印 outbox 关键指标。
// 指标包括：行数、未发送最老年龄、热表与归档表估算体积。
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

// truncateError 截断错误文本，避免超长错误撑爆数据库字段。
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
