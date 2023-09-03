package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	fmt.Println(">> Before:", fromAccount.Balance, toAccount.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// run each transfer money transaction in a separate goroutine
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			// send the result or error to the corresponding channel
			errors <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccountResult := result.FromAccount
		require.NotEmpty(t, fromAccountResult)
		require.Equal(t, fromAccount.ID, fromAccountResult.ID)

		toAccountResult := result.ToAccount
		require.NotEmpty(t, toAccountResult)
		require.Equal(t, toAccount.ID, toAccountResult.ID)

		// check accounts balance
		fmt.Println(">> Tx:", fromAccountResult.Balance, toAccountResult.Balance)
		diff1 := fromAccount.Balance - fromAccountResult.Balance // amount of money that's going out of fromAccount
		diff2 := toAccountResult.Balance - toAccount.Balance     // amount of money that's going into toAccount
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		K := int(diff1 / amount)
		require.True(t, K >= 1 && K <= 5)
		require.NotContains(t, existed, K)
		existed[K] = true
	}

	// check the final updated balance
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> After:", updatedFromAccount.Balance, updatedToAccount.Balance)
	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}
