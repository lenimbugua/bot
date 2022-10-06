package db

import (
	"context"
	"testing"

	"github.com/lenimbugua/bot/util"
	"github.com/stretchr/testify/require"
)

func createRandomCompany(t *testing.T) Company {
	arg := CreateCompanyParams{
		Phone: util.RandomPhoneNumber(),
		Email: util.RandomEmail(),
		Name:  util.RandomString(6),
	}

	company, err := testQueries.CreateCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, company)

	require.Equal(t, arg.Phone, company.Phone)
	require.Equal(t, arg.Email, company.Email)
	require.Equal(t, arg.Name, company.Name)

	require.NotZero(t, company.ID)
	require.NotZero(t, company.CreatedAt)
	require.NotZero(t, company.UpdatedAt)

	return company
}

func TestCreateCompany(t *testing.T) {
	createRandomCompany(t)
}
