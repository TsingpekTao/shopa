package agent

import (
	"testing"

	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
)

func TestFilterPolicyKnowledgeChunksForBuyer(t *testing.T) {
	rows := []entity.AgentKnowledgeChunk{
		{
			ChunkNo:        "CK_PUBLIC",
			SourceTypeCode: "PLATFORM_RULE",
			MetadataJson:   `{"visibility":"public","audience":"buyer"}`,
		},
		{
			ChunkNo:        "CK_INTERNAL",
			SourceTypeCode: "PLATFORM_RULE",
			MetadataJson:   `{"visibility":"internal_only","audience":"operator"}`,
		},
		{
			ChunkNo:        "CK_SELLER",
			SourceTypeCode: "PLATFORM_RULE",
			MetadataJson:   `{"visibility":"public","audience":"seller"}`,
		},
		{
			ChunkNo:        "CK_DEFAULT_PUBLIC",
			SourceTypeCode: "HELP_CENTER",
			MetadataJson:   `{}`,
		},
		{
			ChunkNo:        "CK_DEFAULT_INTERNAL",
			SourceTypeCode: "CUSTOMER_SERVICE_SOP",
			MetadataJson:   `{}`,
		},
	}

	got := filterPolicyKnowledgeChunksForBuyer(rows)
	if len(got) != 2 {
		t.Fatalf("expected two buyer-visible chunks, got %d", len(got))
	}
	if got[0].ChunkNo != "CK_PUBLIC" {
		t.Fatalf("expected first visible chunk CK_PUBLIC, got %q", got[0].ChunkNo)
	}
	if got[1].ChunkNo != "CK_DEFAULT_PUBLIC" {
		t.Fatalf("expected second visible chunk CK_DEFAULT_PUBLIC, got %q", got[1].ChunkNo)
	}
}
