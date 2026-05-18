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

type imageRepo struct {
	db *sqlx.DB
}

func NewImageRepository(db *sqlx.DB) repository.ImageRepository {
	return &imageRepo{db: db}
}

func (r *imageRepo) Create(ctx context.Context, img *models.Image) error {
	query := `INSERT INTO images (id, image) VALUES (:id, :image)`
	_, err := r.db.NamedExecContext(ctx, query, img)
	return err
}

func (r *imageRepo) Replace(ctx context.Context, id uuid.UUID, newData []byte) error {
	query := `UPDATE images SET image = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, newData, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *imageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM images WHERE id = $1`
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

func (r *imageRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Image, error) {
	var img models.Image
	query := `SELECT id, image FROM images WHERE id = $1`
	err := r.db.GetContext(ctx, &img, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &img, nil
}

func (r *imageRepo) GetByProductID(ctx context.Context, productID uuid.UUID) (*models.Image, error) {
	var img models.Image
	query := `SELECT i.id, i.image FROM images i
	          INNER JOIN products p ON p.image_id = i.id
	          WHERE p.id = $1`
	err := r.db.GetContext(ctx, &img, query, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &img, nil
}