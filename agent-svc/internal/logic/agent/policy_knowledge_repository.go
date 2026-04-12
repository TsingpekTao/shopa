package agent

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/TsingpekTao/shopa/agent-svc/internal/dao"
	agentruntime "github.com/TsingpekTao/shopa/agent-svc/internal/logic/agent/runtime"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gerror"
)

type policyKnowledgeRepository struct{}

func newPolicyKnowledgeRepository() *policyKnowledgeRepository {
	return &policyKnowledgeRepository{}
}

func (r *policyKnowledgeRepository) SearchPolicyKnowledge(ctx context.Context, query string, limit uint32) ([]agentruntime.PolicyKnowledgeHit, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	if limit == 0 {
		limit = 3
	}
	if limit > 10 {
		limit = 10
	}

	var rows []entity.AgentKnowledgeChunk
	err := dao.AgentKnowledgeChunk.Ctx(ctx).
		Where(dao.AgentKnowledgeChunk.Columns().IsDeleted, 0).
		WhereLike(dao.AgentKnowledgeChunk.Columns().ChunkText, "%"+query+"%").
		OrderAsc(dao.AgentKnowledgeChunk.Columns().ChunkIndex).
		Limit(int(limit) * 4).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "search policy knowledge failed")
	}
	rows = filterPolicyKnowledgeChunksForBuyer(rows)
	if len(rows) > int(limit) {
		rows = rows[:limit]
	}

	hits := make([]agentruntime.PolicyKnowledgeHit, 0, len(rows))
	for i := range rows {
		snippet := strings.TrimSpace(rows[i].ChunkTextPreview)
		if snippet == "" {
			snippet = strings.TrimSpace(rows[i].ChunkText)
		}
		hits = append(hits, agentruntime.PolicyKnowledgeHit{
			SourceTypeCode: strings.TrimSpace(rows[i].SourceTypeCode),
			SourceID:       strings.TrimSpace(rows[i].SourceId),
			SourceVersion:  rows[i].SourceVersion,
			Title:          firstNonEmpty(strings.TrimSpace(rows[i].KnowledgeDocNo), strings.TrimSpace(rows[i].SourceId), strings.TrimSpace(rows[i].ChunkNo)),
			Snippet:        snippet,
		})
	}
	return hits, nil
}

type policyKnowledgeChunkMetadata struct {
	Visibility string `json:"visibility,omitempty"`
	Audience   string `json:"audience,omitempty"`
}

func filterPolicyKnowledgeChunksForBuyer(rows []entity.AgentKnowledgeChunk) []entity.AgentKnowledgeChunk {
	filtered := make([]entity.AgentKnowledgeChunk, 0, len(rows))
	for _, row := range rows {
		if isPolicyKnowledgeChunkVisibleToBuyer(row) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func isPolicyKnowledgeChunkVisibleToBuyer(row entity.AgentKnowledgeChunk) bool {
	metadata := parsePolicyKnowledgeChunkMetadata(row.MetadataJson)

	visibility := strings.ToLower(strings.TrimSpace(metadata.Visibility))
	if visibility == "" {
		visibility = inferPolicyKnowledgeVisibility(row.SourceTypeCode)
	}
	if visibility != "public" {
		return false
	}

	audience := strings.ToLower(strings.TrimSpace(metadata.Audience))
	if audience == "" {
		audience = inferPolicyKnowledgeAudience(row.SourceTypeCode)
	}

	switch audience {
	case "", "buyer", "universal":
		return true
	default:
		return false
	}
}

func parsePolicyKnowledgeChunkMetadata(raw string) policyKnowledgeChunkMetadata {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return policyKnowledgeChunkMetadata{}
	}

	var metadata policyKnowledgeChunkMetadata
	if err := json.Unmarshal([]byte(raw), &metadata); err != nil {
		return policyKnowledgeChunkMetadata{}
	}
	return metadata
}

func inferPolicyKnowledgeVisibility(sourceTypeCode string) string {
	normalized := normalizePolicyKnowledgeSourceType(sourceTypeCode)
	if normalized == "" {
		return "internal_only"
	}
	if strings.Contains(normalized, "SOP") || strings.Contains(normalized, "INTERNAL") || strings.Contains(normalized, "OPERATOR") {
		return "internal_only"
	}
	switch normalized {
	case "PLATFORM_RULE", "HELP_CENTER", "FAQ", "POLICY", "POLICY_DOC", "PRODUCT_GUIDE", "HELP_DOC":
		return "public"
	default:
		return "internal_only"
	}
}

func inferPolicyKnowledgeAudience(sourceTypeCode string) string {
	normalized := normalizePolicyKnowledgeSourceType(sourceTypeCode)
	if normalized == "" {
		return "operator"
	}
	if strings.Contains(normalized, "SELLER") {
		return "seller"
	}
	if strings.Contains(normalized, "OPERATOR") || strings.Contains(normalized, "SOP") || strings.Contains(normalized, "INTERNAL") {
		return "operator"
	}
	switch normalized {
	case "PLATFORM_RULE", "HELP_CENTER", "FAQ", "POLICY", "POLICY_DOC", "PRODUCT_GUIDE", "HELP_DOC":
		return "universal"
	default:
		return "operator"
	}
}

func normalizePolicyKnowledgeSourceType(sourceTypeCode string) string {
	return strings.ToUpper(strings.TrimSpace(sourceTypeCode))
}
