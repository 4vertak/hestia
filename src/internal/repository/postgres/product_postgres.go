package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) repository.ProductRepository {
	return &productRepo{db: db}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func (r *productRepo) Create(ctx context.Context, p *models.Product) error {
	query := `INSERT INTO products (id, name, category, price, available_stock, last_update_date, supplier_id)
	          VALUES (:id, :name, :category, :price, :available_stock, :last_update_date, :supplier_id)`
	p.LastUpdateDate = timePtr(time.Now())
	_, err := r.db.NamedExecContext(ctx, query, p)
	return err
}

func (r *productRepo) DecreaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `UPDATE products
	          SET available_stock = available_stock - $1,
	              last_update_date = now()
	          WHERE id = $2 AND available_stock >= $1`
	res, err := r.db.ExecContext(ctx, query, quantity, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Проверим, существует ли товар вообще
		var exists bool
		errCheck := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)`, id)
		if errCheck == nil && !exists {
			return sql.ErrNoRows // товар не найден
		}
		return ErrInsufficientStock // товар есть, но остатка недостаточно
	}
	return nil
}

func (r *productRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var p models.Product
	query := `SELECT id, name, category, price, available_stock, last_update_date, supplier_id, image_id
	          FROM products WHERE id = $1`
	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) FindAllAvailable(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	query := `SELECT id, name, category, price, available_stock, last_update_date, supplier_id, image_id
	          FROM products WHERE available_stock > 0 ORDER BY name`
	err := r.db.SelectContext(ctx, &products, query)
	return products, err
}

func (r *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
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

func (r *productRepo) UpdateLastUpdate(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET last_update_date = now() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *productRepo) UpdateImageID(ctx context.Context, productID uuid.UUID, imageID uuid.UUID) error {
	query := `UPDATE products SET image_id = $1, last_update_date = now() WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, imageID, productID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}