package warehouse

import (
	"context"
	"database/sql"
	"ecommerce/pkg/helper"
)

func (r *warehouseRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {

	tx, err := r.Database.BeginTx(ctx, nil)
	if err != nil {
		return nil, helper.HandleError(err)
	}
	return tx, nil
}

func (r *warehouseRepository) RollbackTx(ctx context.Context, tx *sql.Tx) error {

	if err := tx.Rollback(); err != nil {
		return helper.HandleError(err)
	}
	return nil
}

func (r *warehouseRepository) CommitTx(ctx context.Context, tx *sql.Tx) error {

	if err := tx.Commit(); err != nil {
		return helper.HandleError(err)
	}
	return nil
}
