package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a series of queries inside a database transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	qtx := New(tx)

	err = fn(qtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %w", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferMoneyParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type TransferMoneyResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferMoneyTx transfers money from one account to another.
// It creates a transfer record, add account entries, and update accounts balance within a single database transaction.
func (store *Store) TransferMoneyTx(ctx context.Context, arg TransferMoneyParams) (TransferMoneyResult, error) {
	var result TransferMoneyResult

	err := store.execTx(ctx, func(qtx *Queries) error {
		// create transfer record
		transfer, err := qtx.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.Transfer = transfer

		// create account entries
		fromEntry, err := qtx.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry = fromEntry

		toEntry, err := qtx.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry = toEntry

		// Update account balances
		fromAccount, err := qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromAccount = fromAccount

		toAccount, err := qtx.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToAccount = toAccount

		return nil
	})

	return result, err
}
