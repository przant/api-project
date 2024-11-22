package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/przant/api-project/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
        INSERT INTO products (name, description, price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	now := time.Now()
	return r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		now,
		now,
	).Scan(&product.ID)
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	product := &models.Product{}
	query := `
        SELECT id, name, description, price, created_at, updated_at
        FROM products
        WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return product, err
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
        UPDATE products
        SET name = $1, description = $2, price = $3, updated_at = $4
        WHERE id = $5`

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		time.Now(),
		product.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {
	query := `
        SELECT id, name, description, price, created_at, updated_at
        FROM products
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
