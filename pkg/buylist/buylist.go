package buylist

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/regalias/scry/pkg/models"
)

func (m *Manager) NewBuylist(ctx context.Context, name string) (*models.Buylist, error) {
	buylist := &models.Buylist{
		Name:      name,
		CreatedAt: time.Now().UnixMilli(),
	}

	sql, args, err := squirrel.
		Insert(buylistTableName).
		Columns("name", "created_at").
		Values(buylist.Name, buylist.CreatedAt).ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed inserting buylist: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed getting buylist ID after insertion: %w", err)
	}
	buylist.ID = id
	return buylist, nil
}

func (m *Manager) GetBuylist(ctx context.Context, id int64) (*models.Buylist, error) {

	sql, args, err := squirrel.
		Select("*").
		From(buylistTableName).
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	buylist := &models.Buylist{}
	if err := m.db.GetContext(ctx, buylist, sql, args...); err != nil {
		return nil, fmt.Errorf("failed getting buylist: %w", err)
	}

	// Query cards
	cards, err := m.GetCardsForBuylist(ctx, buylist.ID)
	if err != nil {
		return nil, fmt.Errorf("failed getting buylist: %w", err)
	}
	buylist.Cards = cards
	for _, card := range cards {
		buylist.TotalCards += card.Quantity
		buylist.TotalPrice += card.TotalSelectionPrice
		buylist.TotalSelections += card.TotalSelections
	}

	return buylist, nil
}

func (m *Manager) ListBuylists(ctx context.Context) ([]models.Buylist, error) {

	sql, args, err := squirrel.Select("*").From(buylistTableName).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	buylists := []models.Buylist{}

	if err := m.db.Select(&buylists, sql, args...); err != nil {
		return nil, fmt.Errorf("failed listing buylists: %w", err)
	}

	for i := range buylists {
		cards, err := m.GetCardsForBuylist(ctx, buylists[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed listing buylists on %v: %w", buylists[i], err)
		}
		buylists[i].Cards = cards
		for _, card := range cards {
			buylists[i].TotalCards += card.Quantity
			buylists[i].TotalPrice += card.TotalSelectionPrice
			buylists[i].TotalSelections += card.TotalSelections
		}
	}

	return buylists, nil
}

func (m *Manager) DeleteBuylist(ctx context.Context, id int64) error {

	sql, args, err := squirrel.Delete(buylistTableName).Where("id = ?", id).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := m.db.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed deleting buylist: %w", err)
	}

	return nil
}

func (m *Manager) UpdateBuylistName(ctx context.Context, id int64, name string) error {
	sql, args, err := squirrel.Update(buylistTableName).
		Set("name", name).
		Where(squirrel.Eq{"id": id}).ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed updating buylist: %w", err)
	}

	return nil
}
