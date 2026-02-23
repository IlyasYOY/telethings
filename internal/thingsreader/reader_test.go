package thingsreader

import (
	"strings"
	"testing"
)

func TestTasksByTagPageScript_MatchesDescendants(t *testing.T) {
	script := tasksByTagPageScript("Work/ClientA", 0, 11)
	if !strings.Contains(script, `if candidatePath is "Work/ClientA" or candidatePath starts with "Work/ClientA/" then`) {
		t.Fatalf("expected descendant-aware match condition, got script:\n%s", script)
	}
}

func TestTasksByTagPageScript_EscapesQuotes(t *testing.T) {
	script := tasksByTagPageScript(`Work/"ClientA"`, 0, 11)
	if !strings.Contains(script, `candidatePath is "Work/\"ClientA\""`) {
		t.Fatalf("expected escaped quoted tag path, got script:\n%s", script)
	}
}

func TestTasksByTagPageScript_AvoidsGlobalTodoScan(t *testing.T) {
	script := tasksByTagPageScript("Work/ClientA", 10, 11)
	if strings.Contains(script, "repeat with t in to dos\n") {
		t.Fatalf("expected no global todo scan, got script:\n%s", script)
	}
	if !strings.Contains(script, "repeat with t in to dos of tg") {
		t.Fatalf("expected per-tag todo iteration, got script:\n%s", script)
	}
	if !strings.Contains(script, "set neededCount to 10 + 11") {
		t.Fatalf("expected bounded page target count, got script:\n%s", script)
	}
}
