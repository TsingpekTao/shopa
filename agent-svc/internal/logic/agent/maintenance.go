package agent

import (
	"context"
	"strings"
	"time"

	"github.com/TsingpekTao/shopa/agent-svc/internal/dao"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	runtimeMaintenanceInterval = 30 * time.Second
	processingRunTimeout       = 2 * time.Minute
	generatingRunTimeout       = 45 * time.Second
)

func (s *sAgent) startRuntimeMaintenance() {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	go func() {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		ticker := time.NewTicker(runtimeMaintenanceInterval)
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		defer ticker.Stop()
		// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
		for range ticker.C {
			// 统一收口超时或遗留的 Run，把无法恢复的执行明确标记为失败。
			if err := s.reapZombieRuns(context.Background()); err != nil {
				// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
				g.Log().Warning(context.Background(), "agent runtime maintenance failed:", err)
			}
		}
	}()
}

func (s *sAgent) reapZombieRuns(ctx context.Context) error {
	// 扫描并终结超时执行，确保等待工具或生成中的 Run 不会无限挂起。
	if err := s.failExpiredToolRuns(ctx); err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return err
	}
	// 扫描并终结超时执行，确保等待工具或生成中的 Run 不会无限挂起。
	if err := s.failStaleRuns(ctx, runtimeStatusGenerating, gtime.Now().Add(-generatingRunTimeout), "GENERATING_TIMEOUT", "run stayed in generating for too long"); err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return err
	}
	// 扫描并终结超时执行，确保等待工具或生成中的 Run 不会无限挂起。
	if err := s.failStaleRuns(ctx, runtimeStatusProcessing, gtime.Now().Add(-processingRunTimeout), "PROCESSING_STALE", "run was orphaned during processing"); err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return err
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return nil
}

func (s *sAgent) failExpiredToolRuns(ctx context.Context) error {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var rows []entity.AgentRun
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	err := dao.AgentRun.Ctx(ctx).
		Where(dao.AgentRun.Columns().RunStatusCode, runtimeStatusWaitingToolCall).
		WhereLTE(dao.AgentRun.Columns().ToolWaitTimeoutAt, gtime.Now()).
		Scan(&rows)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return gerror.Wrap(err, "query expired tool wait runs failed")
	}
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for i := range rows {
		// 扫描并终结超时执行，确保等待工具或生成中的 Run 不会无限挂起。
		if err = s.failRun(ctx, &rows[i], "TOOL_WAIT_TIMEOUT", "tool call wait timed out"); err != nil {
			// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
			return err
		}
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return nil
}

func (s *sAgent) failStaleRuns(ctx context.Context, runStatus string, cutoff *gtime.Time, errorCode, errorMessage string) error {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var rows []entity.AgentRun
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	err := dao.AgentRun.Ctx(ctx).
		Where(dao.AgentRun.Columns().RunStatusCode, runStatus).
		WhereLTE(dao.AgentRun.Columns().UpdatedAt, cutoff).
		Scan(&rows)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return gerror.Wrap(err, "query stale runs failed")
	}
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for i := range rows {
		// 扫描并终结超时执行，确保等待工具或生成中的 Run 不会无限挂起。
		if err = s.failRun(ctx, &rows[i], errorCode, errorMessage); err != nil {
			// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
			return err
		}
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return nil
}

func (s *sAgent) failRun(ctx context.Context, run *entity.AgentRun, errorCode, errorMessage string) error {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if run == nil || run.Id == 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err := dao.AgentRun.Ctx(ctx).
		Where(dao.AgentRun.Columns().RunNo, run.RunNo).
		Where(dao.AgentRun.Columns().RunStatusCode, run.RunStatusCode).
		Data(do.AgentRun{
			RunStatusCode:     runtimeStatusFailed,
			CurrentNodeCode:   firstNonEmpty(run.CurrentNodeCode, "RuntimeMaintenance"),
			ToolResultStatus:  normalizeToolResultStatus(run.ToolResultStatus),
			ErrorCode:         strings.TrimSpace(errorCode),
			ErrorMessage:      strings.TrimSpace(errorMessage),
			CheckpointVersion: gdb.Raw("checkpoint_version + 1"),
			FinishedAt:        now,
			UpdatedAt:         now,
		}).
		Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return gerror.Wrap(err, "fail stale run failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, _ = dao.AgentConversation.Ctx(ctx).
		Where(dao.AgentConversation.Columns().ConversationNo, run.ConversationNo).
		Data(do.AgentConversation{
			LastRunStatusCode: runtimeStatusFailed,
			UpdatedAt:         now,
		}).
		Update()
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return nil
}
