package db

import (
	"context"
	"testing"
	"time"

	"github.com/lenimbugua/bot/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	company := createRandomCompany(t)
	arg := CreateUserParams{
		Name:         util.RandomString(6),
		PasswordHash: hashedPassword,
		Phone:        util.RandomPhoneNumber(),
		CompanyID:    company.ID,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.Phone, user.Phone)
	require.Equal(t, arg.CompanyID, user.CompanyID)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Phone)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.CompanyID, user2.CompanyID)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}
