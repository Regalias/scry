package buylist

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/regalias/scry/pkg/models"
)

func (m *Manager) GetSelections(ctx context.Context, cardId int64) ([]models.ProductSelection, error) {

	sql, args, err := squirrel.Select("*").From(selectionsTableName).Where("card_id = ?", cardId).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	selections := []models.ProductSelection{}

	rows, err := m.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed getting selections for card: %w", err)
	}
	for rows.Next() {
		item := models.ProductSelection{}
		if err := rows.StructScan(&item); err != nil {
			return nil, fmt.Errorf("failed getting selections for card: %w", err)
		}
		selections = append(selections, item)
	}

	return selections, nil
}

func (m *Manager) AddSelection(ctx context.Context, cardId int64, offering *models.Offering, quantity int64) (*models.ProductSelection, error) {

	selection := &models.ProductSelection{
		Offering: *offering,
		Quantity: quantity,
	}

	serializedOffering, err := selection.Offering.Value()
	if err != nil {
		return nil, fmt.Errorf("failed serializing offering: %w", err)
	}

	sql, args, err := squirrel.Insert(selectionsTableName).
		Columns("quantity", "offering", "is_purchased", "is_flagged", "card_id").
		Values(selection.Quantity, serializedOffering, selection.IsPurchased, selection.IsFlagged, cardId).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last row ID: %w", err)
	}

	selection.ID = id

	return selection, nil
}

func (m *Manager) DeleteSelection(ctx context.Context, selectionId int64) error {
	sql, args, err := squirrel.Delete(selectionsTableName).Where(squirrel.Eq{"id": selectionId}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteSelectionsForCardId(ctx context.Context, cardId int64) error {
	sql, args, err := squirrel.Delete(selectionsTableName).Where(squirrel.Eq{"card_id": cardId}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

type UpdateSelectionRequest struct {
	SelectionId int64            `json:"selectionId"`
	Quantity    *int64           `json:"quantity,omitempty"`
	IsFlagged   *bool            `json:"isFlagged,omitempty"`
	IsPurchased *bool            `json:"isPurchased,omitempty"`
	Offering    *models.Offering `json:"offering,omitempty"`
}

func (m *Manager) UpdateSelection(ctx context.Context, req *UpdateSelectionRequest) error {

	if req.Quantity != nil && *req.Quantity == 0 {
		// Delete instead
		return m.DeleteSelection(ctx, req.SelectionId)
	}

	query := squirrel.Update(selectionsTableName)
	if req.Quantity != nil {
		query = query.Set("quantity", req.Quantity)
	}
	if req.IsFlagged != nil {
		query = query.Set("is_flagged", req.IsFlagged)
	}
	if req.IsPurchased != nil {
		query = query.Set("is_purchased", req.IsPurchased)
	}
	if req.Offering != nil {
		val, err := req.Offering.Value()
		if err != nil {
			return err
		}
		query = query.Set("offering", val)
	}

	query = query.Where(squirrel.Eq{"id": req.SelectionId})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
