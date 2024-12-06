package models

import "gorm.io/gorm"

// Deletion hooks for recurisve deletes

func (m *Buylist) BeforeDelete(tx *gorm.DB) (err error) {
	if result := tx.Unscoped().Delete(&m.Cards); result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *Card) BeforeDelete(tx *gorm.DB) (err error) {
	if result := tx.Unscoped().Delete(&m.Selections); result.Error != nil {
		return result.Error
	}
	return nil
}
