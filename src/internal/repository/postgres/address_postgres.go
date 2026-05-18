package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository"
)

type addressRepo struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) repository.AddressRepository {
	return &addressRepo{db: db}
}

func (r *addressRepo) Create(ctx context.Context, addr *models.Address) error {
	query := `INSERT INTO addresses (id, country, city, street) VALUES (:id, :country, :city, :street)`
	_, err := r.db.NamedExecContext(ctx, query, addr)
	return err
}

func (r *addressRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	var addr models.Address
	query := `SELECT id, country, city, street FROM addresses WHERE id = $1`
	err := r.db.GetContext(ctx, &addr, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &addr, nil
}

func (r *addressRepo) Update(ctx context.Context, addr *models.Address) error {
	query := `UPDATE addresses SET country=:country, city=:city, street=:street WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, addr)
	return err
}

func (r *addressRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM addresses WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}