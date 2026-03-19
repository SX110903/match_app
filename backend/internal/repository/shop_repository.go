package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type shopRepository struct {
	db *database.DB
}

func NewShopRepository(db *database.DB) IShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) CreateTransaction(ctx context.Context, tx *domain.ShopTransaction) error {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	if tx.ID == "" {
		tx.ID = uuid.New().String()
	}
	_, err := r.db.ExecContext(dbCtx,
		`INSERT INTO shop_transactions (id, user_id, item_type, item_value, cost) VALUES (?, ?, ?, ?, ?)`,
		tx.ID, tx.UserID, tx.ItemType, tx.ItemValue, tx.Cost,
	)
	return err
}

func (r *shopRepository) GetTransactionsByUser(ctx context.Context, userID string, limit, offset int) ([]domain.ShopTransaction, error) {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(dbCtx,
		`SELECT id, user_id, item_type, item_value, cost, created_at
		 FROM shop_transactions WHERE user_id = ?
		 ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}
	defer rows.Close()

	var txs []domain.ShopTransaction
	for rows.Next() {
		var t domain.ShopTransaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.ItemType, &t.ItemValue, &t.Cost, &t.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, rows.Err()
}

// PurchaseVIP performs an atomic purchase: deducts credits and upgrades vip_level in one transaction.
func (r *shopRepository) PurchaseVIP(ctx context.Context, userID string, itemValue, cost int) error {
	sqlDB := r.db.DB.DB // unwrap to *sql.DB

	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	sqlTx, err := sqlDB.BeginTx(dbCtx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = sqlTx.Rollback()
		}
	}()

	var credits, vipLevel int
	err = sqlTx.QueryRowContext(dbCtx,
		`SELECT credits, vip_level FROM users WHERE id = ? FOR UPDATE`,
		userID,
	).Scan(&credits, &vipLevel)
	if err != nil {
		return fmt.Errorf("lock user: %w", err)
	}

	if credits < cost {
		err = domain.ErrInvalidInput
		return err
	}
	if itemValue != vipLevel+1 {
		err = domain.ErrInvalidInput
		return err
	}
	if itemValue > 5 {
		err = domain.ErrInvalidInput
		return err
	}

	_, err = sqlTx.ExecContext(dbCtx,
		`UPDATE users SET credits = credits - ?, vip_level = ? WHERE id = ?`,
		cost, itemValue, userID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	txID := uuid.New().String()
	_, err = sqlTx.ExecContext(dbCtx,
		`INSERT INTO shop_transactions (id, user_id, item_type, item_value, cost) VALUES (?, ?, 'vip_upgrade', ?, ?)`,
		txID, userID, itemValue, cost,
	)
	if err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	return sqlTx.Commit()
}
