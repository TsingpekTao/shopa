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

// ApplyRegisterInitEvent 处理“注册初始化资料”事件，将数据投影到 user_profile 表。
// 为了保证一致性，遵守以下约束：
// 1) 事件幂等：同一 event_id 在同一 consumer 中只处理一次；
// 2) 并发串行：同一 user_id 的事件通过事务+行锁顺序执行；
// 3) 用户优先：昵称已被用户手动设置后，系统初始化事件不会覆盖；
// 4) 时间优先：仅接受发生时间晚于上一次初始化事件的记录，屏蔽乱序与重放；
// 5) 版本推进：实际覆盖时递增 profile_version，便于下游做乐观并发控制。
func (s *sUserProfile) ApplyRegisterInitEvent(ctx context.Context, event service.RegisterInitEvent) error {
	// 将字符串字段 trim 掉，避免空白字符导致重复匹配。
	event.EventID = strings.TrimSpace(event.EventID)
	// EventVersion 目前仅用于存储规范化值，便于未来按版本分支。
	event.EventVersion = strings.TrimSpace(event.EventVersion)
	event.InitDisplayName = strings.TrimSpace(event.InitDisplayName)
	// 基础字段校验：用户 ID 和事件 ID 必须存在。
	if event.UserID == 0 || event.EventID == "" {
		return gerror.New("invalid register init event")
	}
	// 若未提供发生时间，则回填当前 UTC 时间，便于后续时序比较。
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}

	// 通过事务包裹幂等记录与 profile 投影更新，保持原子性。
	return dao.UserProfile.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 幂等去重表：同一个 consumer + event_id 只允许首次写入成功。
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
			// InsertIgnore 语义：rows=1 表示首次消费，rows=0 表示已处理，可直接返回 nil。
			if rows, rowsErr := dedupResult.RowsAffected(); rowsErr == nil && rows == 0 {
				return nil
			}
		}

		cols := dao.UserProfile.Columns()
		// 先将对应 user_profile 行加 FOR UPDATE 锁：
		// - 保证同一 user_id 的初始化事件在事务中串行执行；
		// - 避免后写覆盖比自己早提交的变更，保持昵称来源确定性。
		record, err := tx.Model(dao.UserProfile.Table()).
			Where(cols.UserId, event.UserID).
			LockUpdate().
			One()
		if err != nil {
			return err
		}

		occurredAt := gtime.NewFromTime(event.OccurredAt.UTC())
		// 初始化昵称统一归一化，空值回滚到“用户”。
		initName := normalizeInitDisplayName(event.InitDisplayName)

		if record == nil || record.IsEmpty() {
			// 首次初始化：
			// - 写入 display_name 并标注来源为 system_init；
			// - 记录 last_init_event_* 以便后续时序剪裁；
			// - profile_version 从 1 开始。
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

		// 若用户已自行设置昵称，系统初始化事件不再覆盖，确保“用户优先”。
		if profile.DisplayNameSource == consts.DisplayNameSourceUserSet {
			return nil
		}

		// 已有初始化事件时，仅接受发生时间更晚的事件：
		// - 当前事件时间 <= 已记录时间：视为乱序/重放，忽略；
		// - 当前事件时间 > 已记录时间：继续覆盖。
		if profile.LastInitEventAt != nil && !profile.LastInitEventAt.IsZero() {
			lastAt := time.Unix(0, profile.LastInitEventAt.TimestampNano()).UTC()
			if !event.OccurredAt.UTC().After(lastAt) {
				return nil
			}
		}

		// 执行覆盖更新：
		// 1) 更新显示名、来源与最后初始化事件游标；
		// 2) 用 SQL 表达式原子递增 profile_version，避免读写竞争窗口。
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

// normalizeInitDisplayName 归一化初始化昵称。
// 约定：空昵称统一回落为“用户”，避免写入空展示名。
func normalizeInitDisplayName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "\u7528\u6237"
	}
	return name
}
