package models

import (
	"time"
	"github.com/google/uuid"
)

// ClientWithAddress – результат JOIN clients + addresses
type ClientWithAddress struct {
	ID               uuid.UUID  `db:"id"`
	Name             string     `db:"name"`
	Surname          string     `db:"surname"`
	Birthday         *time.Time `db:"birthday"`
	Gender           *string    `db:"gender"`
	RegistrationDate time.Time  `db:"registration_date"`
	AddressID        uuid.UUID  `db:"address_id"`
	Country          string     `db:"country"`
	City             string     `db:"city"`
	Street           string     `db:"street"`
}