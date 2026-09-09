package migration

import (
	"strings"
	"testing"
)

func TestEmbeddedMigrationsAreOrderedAndParseable(t *testing.T) {
	files, err := loadFiles()
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no embedded migrations found")
	}
	for i := range files {
		if files[i].upSQL == "" {
			t.Fatalf("migration %s has empty Up SQL", files[i].name)
		}
		if i > 0 && files[i-1].version >= files[i].version {
			t.Fatalf("migrations are not strictly ordered: %s then %s", files[i-1].name, files[i].name)
		}
	}
}

func TestSplitSQLStatements(t *testing.T) {
	source := `
ALTER TYPE payment_status ADD VALUE IF NOT EXISTS 'REFUNDED';
-- A semicolon in a comment must not split; here.
DO $body$
BEGIN
    PERFORM 'value;with;semicolons';
END;
$body$ LANGUAGE plpgsql;
SELECT "semi;colon", 'it''s;safe';
/* Nested comments are valid in PostgreSQL: /* inner; */ outer; */
SELECT 1;
`
	statements, err := splitSQLStatements(source)
	if err != nil {
		t.Fatalf("split SQL: %v", err)
	}
	if len(statements) != 4 {
		t.Fatalf("expected 4 statements, got %d: %#v", len(statements), statements)
	}
	if !strings.Contains(statements[1], "PERFORM 'value;with;semicolons'") {
		t.Fatalf("dollar-quoted statement was split incorrectly: %q", statements[1])
	}
}

func TestSplitSQLStatementsRejectsUnterminatedInput(t *testing.T) {
	if _, err := splitSQLStatements("SELECT 'unterminated;"); err == nil {
		t.Fatal("expected unterminated quote error")
	}
}
