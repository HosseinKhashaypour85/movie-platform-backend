package migrations

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed *.up.sql
var files embed.FS

func Up(db *gorm.DB) error {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) NOT NULL PRIMARY KEY,
		applied_at DATETIME(3) NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`).Error; err != nil {
		return err
	}
	names, err := fs.Glob(files, "*.up.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)
	for _, name := range names {
		var count int64
		if err := db.Table("schema_migrations").Where("version = ?", name).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		raw, err := files.ReadFile(name)
		if err != nil {
			return err
		}
		for _, statement := range strings.Split(string(raw), ";") {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}
			if err := db.Exec(statement).Error; err != nil {
				return fmt.Errorf("migration %s: %w", name, err)
			}
		}
		if err := db.Exec("INSERT INTO schema_migrations (version,applied_at) VALUES (?,UTC_TIMESTAMP(3))", name).Error; err != nil {
			return fmt.Errorf("record migration %s: %w", name, err)
		}
	}
	return nil
}
