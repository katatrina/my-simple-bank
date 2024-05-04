package db

import (
	"context"
	"testing"

	"github.com/katatrina/my-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)
	require.Zero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)

	require.NotEmpty(t, user2)
	require.Equal(t, user2.Username, user1.Username)
	require.Equal(t, user2.FullName, user1.FullName)
	require.Equal(t, user2.Email, user1.Email)
	require.Equal(t, user2.PasswordChangedAt, user1.PasswordChangedAt)
	require.Equal(t, user2.CreatedAt, user1.CreatedAt)
}

// func TestUpdateAccountBalance(t *testing.T) {
// 	account1 := createRandomAccount(t)
//
// 	arg := UpdateAccountBalanceParams{
// 		ID:      account1.ID,
// 		Balance: util.RandomMoney(),
// 	}
// 	account2, err := testQueries.UpdateAccountBalance(context.Background(), arg)
// 	require.NoError(t, err)
//
// 	require.NotEmpty(t, account2)
// 	require.Equal(t, account2.ID, account1.ID)
// 	require.Equal(t, account2.Owner, account1.Owner)
// 	require.Equal(t, account2.Balance, arg.Balance)
// 	require.Equal(t, account2.Currency, account2.Currency)
// 	require.Equal(t, account2.CreatedAt, account2.CreatedAt)
// }
//
// func TestDeleteAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)
//
// 	err := testQueries.DeleteAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)
//
// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, account2)
// }
//
// func TestListAccounts(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createRandomAccount(t)
// 	}
//
// 	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
// 		Limit:  5,
// 		Offset: 5,
// 	})
// 	require.NoError(t, err)
// 	require.Len(t, accounts, 5)
//
// 	for _, account := range accounts {
// 		require.NotEmpty(t, account)
// 	}
// }
