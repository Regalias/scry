package buylist

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/regalias/scry/pkg/models"
)

func (m *Manager) GetCardsForBuylist(ctx context.Context, buylistId int64) ([]models.Card, error) {

	sql, args, err := squirrel.Select("*").From(cardsTableName).Where("buylist_id = ?", buylistId).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	cards := []models.Card{}

	rows, err := m.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed getting cards for buylist: %w", err)
	}
	for rows.Next() {
		card := models.Card{}
		if err := rows.StructScan(&card); err != nil {
			return nil, fmt.Errorf("failed getting cards for buylist: %w", err)
		}
		cards = append(cards, card)
	}

	// Fetch all selections
	for i := range cards {
		selections, err := m.GetSelections(ctx, cards[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed getting cards for buylist: %w", err)
		}
		cards[i].Selections = selections

		// Calc current price/qty for selections
		for _, sel := range selections {
			cards[i].TotalSelections += sel.Quantity
			cards[i].TotalSelectionPrice += sel.Offering.Price
		}
	}

	return cards, nil
}

type AddCardsRequest struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}

func (m *Manager) AddCardsToBuylist(ctx context.Context, buylistId int64, cards []*AddCardsRequest) error {

	q := squirrel.Insert(cardsTableName).Columns("name", "quantity", "buylist_id")

	for i := range cards {
		q = q.Values(cards[i].Name, cards[i].Quantity, buylistId)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err = m.db.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed adding cards for buylist: %w", err)
	}
	return nil
}

func (m *Manager) DeleteCards(ctx context.Context, cardIds []int64) error {

	sql, args, err := squirrel.Delete(cardsTableName).Where(squirrel.Eq{"id": cardIds}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err = m.db.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed deleting cards: %w", err)
	}
	return nil
}

func (m *Manager) DeleteCardsForBuylist(ctx context.Context, tx *sql.Tx, buylistId int64) error {

	sql, args, err := squirrel.Delete(cardsTableName).Where(squirrel.Eq{"buylist_id": buylistId}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if tx != nil {
		if _, err = tx.ExecContext(ctx, sql, args...); err != nil {
			return fmt.Errorf("failed getting cards for buylist: %w", err)
		}
	} else {
		if _, err = m.db.ExecContext(ctx, sql, args...); err != nil {
			return fmt.Errorf("failed getting cards for buylist: %w", err)
		}
	}

	return nil
}

func (m *Manager) UpdateCardQty(ctx context.Context, cardId int64, quantity int64) error {

	sql, args, err := squirrel.Update(cardsTableName).
		SetMap(map[string]interface{}{
			"quantity": quantity,
		}).
		Where(squirrel.Eq{"id": cardId}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err = m.db.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed updating card: %w", err)
	}
	return nil
}
