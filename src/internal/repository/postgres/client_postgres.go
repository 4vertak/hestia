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

type clientRepo struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) repository.ClientRepository {
	return &clientRepo{db: db}
}

func (r *clientRepo) Create(ctx context.Context, c *models.Client) error {
	query := `INSERT INTO clients (id, name, surname, birthday, gender, registration_date, address_id)
	          VALUES (:id, :name, :surname, :birthday, :gender, :registration_date, :address_id)`
	_, err := r.db.NamedExecContext(ctx, query, c)
	return err
}

func (r *clientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM clients WHERE id = $1`
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

func (r *clientRepo) FindByNameAndSurname(ctx context.Context, name, surname string) ([]models.ClientWithAddress, error) {
	var result []models.ClientWithAddress
	query := `SELECT c.id, c.name, c.surname, c.birthday, c.gender, c.registration_date, c.address_id,
	          a.country, a.city, a.street
	          FROM clients c
	          JOIN addresses a ON c.address_id = a.id
	          WHERE c.name = $1 AND c.surname = $2`
	err := r.db.SelectContext(ctx, &result, query, name, surname)
	return result, err
}

func (r *clientRepo) FindAll(ctx context.Context, limit, offset int) ([]models.ClientWithAddress, error) {
	var result []models.ClientWithAddress
	query := `SELECT c.id, c.name, c.surname, c.birthday, c.gender, c.registration_date, c.address_id,
	          a.country, a.city, a.street
	          FROM clients c
	          JOIN addresses a ON c.address_id = a.id
	          ORDER BY c.registration_date DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &result, query, limit, offset)
	return result, err
}

func (r *clientRepo) UpdateAddress(ctx context.Context, id uuid.UUID, newAddressID uuid.UUID) error {
	query := `UPDATE clients SET address_id = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newAddressID, id)
	return err
}

func (r *clientRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	var c models.Client
	query := `SELECT id, name, surname, birthday, gender, registration_date, address_id
	          FROM clients WHERE id = $1`
	err := r.db.GetContext(ctx, &c, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}