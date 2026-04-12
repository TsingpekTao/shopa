package agent

import (
	agentruntime "github.com/TsingpekTao/shopa/agent-svc/internal/logic/agent/runtime"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
)

func resolveConversationTaskSession(runs []entity.AgentRun) agentruntime.TaskSessionState {
	var (
		bestSession agentruntime.TaskSessionState
		bestRunID   uint64
		hasBest     bool
	)

	for _, run := range runs {
		session := agentruntime.ParseGraphTaskSession(run.GraphStateJson)
		if session == nil {
			continue
		}
		if !hasBest || shouldPreferTaskSession(bestSession, bestRunID, *session, run.Id) {
			bestSession = *session
			bestRunID = run.Id
			hasBest = true
		}
	}

	return bestSession
}

func shouldPreferTaskSession(current agentruntime.TaskSessionState, currentRunID uint64, candidate agentruntime.TaskSessionState, candidateRunID uint64) bool {
	if candidate.SessionVersion != current.SessionVersion {
		return candidate.SessionVersion > current.SessionVersion
	}

	currentScore := taskSessionPriorityScore(current)
	candidateScore := taskSessionPriorityScore(candidate)
	if candidateScore != currentScore {
		return candidateScore > currentScore
	}

	return candidateRunID > currentRunID
}

func taskSessionPriorityScore(session agentruntime.TaskSessionState) int {
	score := 0
	if session.SelectedOrderNo != "" {
		score += 100
	}
	if session.SelectedSubOrderNo != "" {
		score += 80
	}
	if len(session.SlotValues) > 0 {
		score += 40
	}
	if len(session.MissingSlots) > 0 {
		score += 20
	}
	if session.ActiveTaskCode != "" {
		score += 10
	}
	if session.PendingSelectionQuery != "" {
		score += 5
	}
	return score
}
