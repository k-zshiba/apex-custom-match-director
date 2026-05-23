package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/koshi/apex-custom-match-director/internal/domain"
	_ "modernc.org/sqlite"
)

type SettingsStore struct {
	db *sql.DB
}

func NewSettingsStore(path string) (*SettingsStore, error) {
	if path == "" {
		path = "apex-custom-match-director.sqlite"
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return nil, fmt.Errorf("create settings directory: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open settings database: %w", err)
	}
	store := &SettingsStore{db: db}
	if err := store.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SettingsStore) Load(ctx context.Context) (domain.Settings, error) {
	row := s.db.QueryRowContext(ctx, `select value from settings where key = 'settings'`)
	var raw string
	err := row.Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.DefaultSettings(), nil
	}
	if err != nil {
		return domain.Settings{}, fmt.Errorf("load settings: %w", err)
	}
	settings := domain.DefaultSettings()
	if err := json.Unmarshal([]byte(raw), &settings); err != nil {
		return domain.Settings{}, fmt.Errorf("decode settings: %w", err)
	}
	return domain.NormalizeSettings(settings), nil
}

func (s *SettingsStore) Save(ctx context.Context, settings domain.Settings) error {
	settings = domain.NormalizeSettings(settings)
	raw, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("encode settings: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		insert into settings(key, value) values('settings', ?)
		on conflict(key) do update set value = excluded.value
	`, string(raw))
	if err != nil {
		return fmt.Errorf("save settings: %w", err)
	}
	return nil
}

func (s *SettingsStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SettingsStore) migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		create table if not exists settings (
			key text primary key,
			value text not null
		)
	`)
	if err != nil {
		return fmt.Errorf("migrate settings database: %w", err)
	}
	return nil
}
