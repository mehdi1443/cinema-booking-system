package migration

import "testing"

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
