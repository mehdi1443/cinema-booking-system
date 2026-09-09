package migration

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mehdi1443/cinema-booking-system/migrations"
)

const advisoryLockID int64 = 0x63696D61

type migrationFile struct {
	version       int64
	name          string
	upSQL         string
	noTransaction bool
}

// Up applies embedded migrations on one PostgreSQL session. A session-level
// advisory lock prevents concurrent application replicas from racing.
func Up(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("migration database is nil")
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquire migration connection: %w", err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", advisoryLockID); err != nil {
		return fmt.Errorf("lock migrations: %w", err)
	}
	defer conn.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", advisoryLockID)

	if _, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS app_schema_migrations (
version BIGINT PRIMARY KEY, filename TEXT NOT NULL, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`); err != nil {
		return fmt.Errorf("create migration ledger: %w", err)
	}
	files, err := loadFiles()
	if err != nil {
		return err
	}
	gooseVersion, err := currentGooseVersion(ctx, conn)
	if err != nil {
		return err
	}
	for _, file := range files {
		var exists bool
		if err := conn.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM app_schema_migrations WHERE version=$1)", file.version).Scan(&exists); err != nil {
			return fmt.Errorf("read migration ledger: %w", err)
		}
		if exists {
			continue
		}
		if file.version <= gooseVersion {
			if _, err := conn.ExecContext(ctx, "INSERT INTO app_schema_migrations(version, filename) VALUES($1,$2) ON CONFLICT DO NOTHING", file.version, file.name); err != nil {
				return fmt.Errorf("import Goose migration %s: %w", file.name, err)
			}
			continue
		}
		if err := applyFile(ctx, conn, file); err != nil {
			return err
		}
	}
	return nil
}

func currentGooseVersion(ctx context.Context, conn *sql.Conn) (int64, error) {
	var table *string
	if err := conn.QueryRowContext(ctx, "SELECT to_regclass('public.goose_db_version')::text").Scan(&table); err != nil {
		return 0, fmt.Errorf("detect Goose ledger: %w", err)
	}
	if table == nil {
		return 0, nil
	}
	var version int64
	if err := conn.QueryRowContext(ctx, "SELECT COALESCE(MAX(version_id),0) FROM goose_db_version WHERE is_applied = TRUE").Scan(&version); err != nil {
		return 0, fmt.Errorf("read Goose version: %w", err)
	}
	return version, nil
}

func applyFile(ctx context.Context, conn *sql.Conn, file migrationFile) error {
	insert := "INSERT INTO app_schema_migrations(version, filename, applied_at) VALUES($1,$2,$3)"
	if file.noTransaction {
		statements, err := splitSQLStatements(file.upSQL)
		if err != nil {
			return fmt.Errorf("parse non-transactional migration %s: %w", file.name, err)
		}
		for index, statement := range statements {
			// Each call must be a separate autocommit transaction. PostgreSQL enum
			// values added by one statement cannot safely be referenced until that
			// statement has committed.
			if _, err := conn.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("apply migration %s statement %d: %w", file.name, index+1, err)
			}
		}
		if _, err := conn.ExecContext(ctx, insert, file.version, file.name, time.Now()); err != nil {
			return fmt.Errorf("record migration %s: %w", file.name, err)
		}
		return nil
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", file.name, err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, file.upSQL); err != nil {
		return fmt.Errorf("apply migration %s: %w", file.name, err)
	}
	if _, err := tx.ExecContext(ctx, insert, file.version, file.name, time.Now()); err != nil {
		return fmt.Errorf("record migration %s: %w", file.name, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", file.name, err)
	}
	return nil
}

// splitSQLStatements separates top-level SQL statements while preserving
// semicolons inside quoted strings, comments, identifiers, and dollar-quoted
// PostgreSQL function bodies.
func splitSQLStatements(source string) ([]string, error) {
	statements := make([]string, 0)
	start := 0
	inSingleQuote := false
	inDoubleQuote := false
	inLineComment := false
	blockCommentDepth := 0
	dollarTag := ""

	for index := 0; index < len(source); index++ {
		if inLineComment {
			if source[index] == '\n' {
				inLineComment = false
			}
			continue
		}
		if blockCommentDepth > 0 {
			if index+1 < len(source) && source[index:index+2] == "/*" {
				blockCommentDepth++
				index++
			} else if index+1 < len(source) && source[index:index+2] == "*/" {
				blockCommentDepth--
				index++
			}
			continue
		}
		if dollarTag != "" {
			if strings.HasPrefix(source[index:], dollarTag) {
				index += len(dollarTag) - 1
				dollarTag = ""
			}
			continue
		}
		if inSingleQuote {
			if source[index] == '\'' {
				if index+1 < len(source) && source[index+1] == '\'' {
					index++
				} else {
					inSingleQuote = false
				}
			}
			continue
		}
		if inDoubleQuote {
			if source[index] == '"' {
				if index+1 < len(source) && source[index+1] == '"' {
					index++
				} else {
					inDoubleQuote = false
				}
			}
			continue
		}

		switch {
		case index+1 < len(source) && source[index:index+2] == "--":
			inLineComment = true
			index++
		case index+1 < len(source) && source[index:index+2] == "/*":
			blockCommentDepth = 1
			index++
		case source[index] == '\'':
			inSingleQuote = true
		case source[index] == '"':
			inDoubleQuote = true
		case source[index] == '$':
			if tag := readDollarTag(source[index:]); tag != "" {
				dollarTag = tag
				index += len(tag) - 1
			}
		case source[index] == ';':
			if statement := strings.TrimSpace(source[start:index]); statement != "" {
				statements = append(statements, statement)
			}
			start = index + 1
		}
	}

	if inSingleQuote || inDoubleQuote || dollarTag != "" || blockCommentDepth > 0 {
		return nil, fmt.Errorf("unterminated SQL quote or comment")
	}
	if statement := strings.TrimSpace(source[start:]); statement != "" {
		statements = append(statements, statement)
	}
	if len(statements) == 0 {
		return nil, fmt.Errorf("migration contains no executable statements")
	}
	return statements, nil
}

func readDollarTag(source string) string {
	if len(source) < 2 || source[0] != '$' {
		return ""
	}
	for index := 1; index < len(source); index++ {
		switch character := source[index]; {
		case character == '$':
			return source[:index+1]
		case character == '_' || character >= 'a' && character <= 'z' || character >= 'A' && character <= 'Z' || index > 1 && character >= '0' && character <= '9':
			continue
		default:
			return ""
		}
	}
	return ""
}

func loadFiles() ([]migrationFile, error) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("read embedded migrations: %w", err)
	}
	result := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		prefix, _, ok := strings.Cut(entry.Name(), "_")
		if !ok {
			return nil, fmt.Errorf("invalid migration filename %q", entry.Name())
		}
		version, err := strconv.ParseInt(prefix, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid migration version %q: %w", entry.Name(), err)
		}
		body, err := fs.ReadFile(migrations.FS, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", entry.Name(), err)
		}
		text := string(body)
		upMarker, downMarker := "-- +goose Up", "-- +goose Down"
		start := strings.Index(text, upMarker)
		if start < 0 {
			return nil, fmt.Errorf("migration %q has no Up section", entry.Name())
		}
		up := text[start+len(upMarker):]
		if end := strings.Index(up, downMarker); end >= 0 {
			up = up[:end]
		}
		up = strings.ReplaceAll(up, "-- +goose StatementBegin", "")
		up = strings.ReplaceAll(up, "-- +goose StatementEnd", "")
		result = append(result, migrationFile{version: version, name: entry.Name(), upSQL: strings.TrimSpace(up), noTransaction: strings.Contains(text[:start], "-- +goose NO TRANSACTION")})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].version < result[j].version })
	return result, nil
}
