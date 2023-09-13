package db

import (
	"context"
	"database/sql"
	"github.com/katatrina/my-simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	expectedAccount := createRandomAccount(t)

	actualAccount, err := testQueries.GetAccount(context.Background(), expectedAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, actualAccount)

	require.Equal(t, expectedAccount.ID, actualAccount.ID)
	require.Equal(t, expectedAccount.Owner, actualAccount.Owner)
	require.Equal(t, expectedAccount.Balance, actualAccount.Balance)
	require.Equal(t, expectedAccount.Currency, actualAccount.Currency)
	require.WithinDuration(t, expectedAccount.CreatedAt, actualAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	expectedAccount := createRandomAccount(t)

	arg := UpdateAccountBalanceParams{
		ID:      expectedAccount.ID,
		Balance: util.RandomMoney(),
	}

	actualAccount, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, actualAccount)

	require.Equal(t, expectedAccount.ID, actualAccount.ID)
	require.Equal(t, expectedAccount.Owner, actualAccount.Owner)
	require.Equal(t, arg.Balance, actualAccount.Balance)
	require.Equal(t, expectedAccount.Currency, actualAccount.Currency)
	require.WithinDuration(t, expectedAccount.CreatedAt, actualAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	expectedAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), expectedAccount.ID)
	require.NoError(t, err)

	actualAccount, err := testQueries.GetAccount(context.Background(), expectedAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, actualAccount)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
