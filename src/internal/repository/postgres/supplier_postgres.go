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

type supplierRepo struct {
	db *sqlx.DB
}

func NewSupplierRepository(db *sqlx.DB) repository.SupplierRepository {
	return &supplierRepo{db: db}
}

func (r *supplierRepo) Create(ctx context.Context, s *models.Supplier) error {
	query := `INSERT INTO suppliers (id, name, address_id, phone_number)
	          VALUES (:id, :name, :address_id, :phone_number)`
	_, err := r.db.NamedExecContext(ctx, query, s)
	return err
}

func (r *supplierRepo) UpdateAddress(ctx context.Context, id uuid.UUID, newAddressID uuid.UUID) error {
	query := `UPDATE suppliers SET address_id = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, newAddressID, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *supplierRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM suppliers WHERE id = $1`
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

func (r *supplierRepo) FindAll(ctx context.Context) ([]models.SupplierWithAddress, error) {
	var result []models.SupplierWithAddress
	query := `SELECT s.id, s.name, s.address_id, s.phone_number, a.country, a.city, a.street
	          FROM suppliers s
	          JOIN addresses a ON s.address_id = a.id
	          ORDER BY s.name`
	err := r.db.SelectContext(ctx, &result, query)
	return result, err
}

func (r *supplierRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.SupplierWithAddress, error) {
	var s models.SupplierWithAddress
	query := `SELECT s.id, s.name, s.address_id, s.phone_number, a.country, a.city, a.street
	          FROM suppliers s
	          JOIN addresses a ON s.address_id = a.id
	          WHERE s.id = $1`
	err := r.db.GetContext(ctx, &s, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}