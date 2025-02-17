package product

import (
	"context"
	"database/sql"
	"ecommerce/pkg/helper"
)

func (r *productRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {

	tx, err := r.Database.BeginTx(ctx, nil)
	if err != nil {
		return nil, helper.HandleError(err)
	}
	return tx, nil
}

func (r *productRepository) RollbackTx(ctx context.Context, tx *sql.Tx) error {

	if err := tx.Rollback(); err != nil {
		return helper.HandleError(err)
	}
	return nil
}

func (r *productRepository) CommitTx(ctx context.Context, tx *sql.Tx) error {

	if err := tx.Commit(); err != nil {
		return helper.HandleError(err)
	}
	return nil
}
