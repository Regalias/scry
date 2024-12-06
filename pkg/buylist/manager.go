package buylist

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Manager struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func NewManager(logger *slog.Logger) (*Manager, error) {

	db, err := sqlx.Connect("sqlite3", "file:buylists.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(buylistSchema); err != nil {
		return nil, fmt.Errorf("failed creating buylist table: %w", err)
	}

	if _, err := db.Exec(cardSchema); err != nil {
		return nil, fmt.Errorf("failed creating card table: %w", err)
	}

	if _, err := db.Exec(selectionsSchema); err != nil {
		return nil, fmt.Errorf("failed creating selections table: %w", err)
	}

	return &Manager{
		logger: logger,
		db:     db,
	}, nil
}

func (m *Manager) Shutdown() error {
	return m.db.Close()
}
