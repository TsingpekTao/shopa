package userprofile

import (
	"context"
	"strings"
	"time"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/consts"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/dao"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// ApplyRegisterInitEvent 消费并应用领域事件。
func (s *sUserProfile) ApplyRegisterInitEvent(ctx context.Context, event service.RegisterInitEvent) error {
	event.EventID = strings.TrimSpace(event.EventID)
	event.EventVersion = strings.TrimSpace(event.EventVersion)
	event.InitDisplayName = strings.TrimSpace(event.InitDisplayName)
	if event.UserID == 0 || event.EventID == "" {
		return gerror.New("invalid register init event")
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}

	return dao.UserProfile.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		dedupResult, err := tx.Model(dao.ConsumerEventDedup.Table()).
			Data(do.ConsumerEventDedup{
				ConsumerName: consts.ConsumerRegisterInit,
				EventId:      event.EventID,
				EventType:    "UserRegisteredV1",
				ProcessedAt:  gtime.Now(),
			}).
			InsertIgnore()
		if err != nil {
			return err
		}
		if dedupResult != nil {
			if rows, rowsErr := dedupResult.RowsAffected(); rowsErr == nil && rows == 0 {
				return nil
			}
		}

		cols := dao.UserProfile.Columns()
		record, err := tx.Model(dao.UserProfile.Table()).
			Where(cols.UserId, event.UserID).
			LockUpdate().
			One()
		if err != nil {
			return err
		}

		occurredAt := gtime.NewFromTime(event.OccurredAt.UTC())
		initName := normalizeInitDisplayName(event.InitDisplayName)

		if record == nil || record.IsEmpty() {
			_, err = tx.Model(dao.UserProfile.Table()).
				Data(do.UserProfile{
					UserId:            event.UserID,
					DisplayName:       initName,
					ProfileVersion:    1,
					DisplayNameSource: consts.DisplayNameSourceSystemInit,
					LastInitEventAt:   occurredAt,
					LastInitEventId:   event.EventID,
				}).
				Insert()
			return err
		}

		var profile entity.UserProfile
		if err = record.Struct(&profile); err != nil {
			return err
		}

		// 用户手动改过昵称后，不允许延迟到达的初始化事件覆盖。
		if profile.DisplayNameSource == consts.DisplayNameSourceUserSet {
			return nil
		}

		if profile.LastInitEventAt != nil && !profile.LastInitEventAt.IsZero() {
			lastAt := time.Unix(0, profile.LastInitEventAt.TimestampNano()).UTC()
			if !event.OccurredAt.UTC().After(lastAt) {
				return nil
			}
		}

		_, err = tx.Model(dao.UserProfile.Table()).
			Where(cols.UserId, event.UserID).
			Data(do.UserProfile{
				DisplayName:       initName,
				DisplayNameSource: consts.DisplayNameSourceSystemInit,
				LastInitEventAt:   occurredAt,
				LastInitEventId:   event.EventID,
				ProfileVersion:    gdb.Raw(cols.ProfileVersion + " + 1"),
			}).
			Update()
		return err
	})
}

// normalizeInitDisplayName 归一化事件里的初始化昵称。
func normalizeInitDisplayName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "\u7528\u6237"
	}
	return name
}
