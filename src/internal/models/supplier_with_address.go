package models

import (
	"github.com/google/uuid"
)

// SupplierWithAddress – результат JOIN suppliers + addresses
type SupplierWithAddress struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	AddressID   uuid.UUID `db:"address_id"`
	PhoneNumber *string   `db:"phone_number"`
	Country     string    `db:"country"`
	City        string    `db:"city"`
	Street      string    `db:"street"`
}