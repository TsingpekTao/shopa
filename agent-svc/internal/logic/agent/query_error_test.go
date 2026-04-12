package agent

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func TestNormalizeOptionalQueryError(t *testing.T) {
	t.Parallel()

	if err := normalizeOptionalQueryError(nil); err != nil {
		t.Fatalf("expected nil to stay nil, got %v", err)
	}
	if err := normalizeOptionalQueryError(sql.ErrNoRows); err != nil {
		t.Fatalf("expected sql.ErrNoRows to be ignored, got %v", err)
	}

	sentinel := errors.New("boom")
	if err := normalizeOptionalQueryError(sentinel); !errors.Is(err, sentinel) {
		t.Fatalf("expected non-empty error to pass through, got %v", err)
	}
}

func TestNormalizeRequiredQueryError(t *testing.T) {
	t.Parallel()

	if err := normalizeRequiredQueryError(nil, "missing"); err != nil {
		t.Fatalf("expected nil to stay nil, got %v", err)
	}

	err := normalizeRequiredQueryError(sql.ErrNoRows, "missing")
	if err == nil {
		t.Fatal("expected sql.ErrNoRows to become not found error")
	}
	if gerrorCode := gerror.Code(err); gerrorCode != gcode.CodeNotFound {
		t.Fatalf("expected not found code, got %v", gerrorCode)
	}
	if err.Error() != "missing" {
		t.Fatalf("expected not found message to be preserved, got %q", err.Error())
	}
}
