package db

import (
	"context"
	"github.com/katatrina/my-simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	expectedUser := createRandomUser(t)

	actualUser, err := testQueries.GetUser(context.Background(), expectedUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, actualUser)

	require.Equal(t, expectedUser.Username, actualUser.Username)
	require.Equal(t, expectedUser.HashedPassword, actualUser.HashedPassword)
	require.Equal(t, expectedUser.FullName, actualUser.FullName)
	require.Equal(t, expectedUser.Email, actualUser.Email)
	require.WithinDuration(t, expectedUser.PasswordChangedAt, actualUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, expectedUser.CreatedAt, actualUser.CreatedAt, time.Second)
}
