package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

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

	require.NotZero(t, company.CreatedAt)
	require.NotZero(t, company.UpdatedAt)

	return company
}

func TestCreateCompany(t *testing.T) {
	createRandomCompany(t)
}

func TestGetCompanyByEmail(t *testing.T) {
	company := createRandomCompany(t)

	require := require.New(t)

	company1, err := testQueries.GetCompanyByEmail(context.Background(), company.Email)
	require.NoError(err)
	require.NotEmpty(company1)

	require.Equal(company.Phone, company1.Phone)
	require.Equal(company.Email, company1.Email)
	require.Equal(company.Name, company1.Name)

	require.Equal(company.CreatedAt, company1.CreatedAt)
	require.Equal(company.CreatedAt, company.UpdatedAt)
}

func TestListCompanies(t *testing.T) {
	require := require.New(t)
	for i := 0; i < 10; i++ {
		createRandomCompany(t)
	}
	arg := ListCompaniesParams{
		Limit:  5,
		Offset: 5,
	}

	companies, err := testQueries.ListCompanies(context.Background(), arg)
	require.NoError(err)
	require.Equal(len(companies), 5)

	for _, company := range companies {
		require.NotEmpty(company)
	}
}

func TestUpdateCompanyAllFields(t *testing.T) {
	require := require.New(t)

	company := createRandomCompany(t)

	name := util.RandomString(6)
	email := util.RandomString(6)
	phone := util.RandomPhoneNumber()
	updatedAt := time.Now()
	arg := UpdateCompanyParams{
		Name:      sql.NullString{String: name, Valid: true},
		Email:     sql.NullString{String: email, Valid: true},
		Phone:     sql.NullString{String: phone, Valid: true},
		ID:        company.ID,
	}
	updatedCompany, err := testQueries.UpdateCompany(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(updatedCompany)
	require.NotEqual(company.Name, updatedCompany.Name)
	require.Equal(name, updatedCompany.Name)
	require.NotEqual(company.Email, updatedCompany.Email)
	require.Equal(email, updatedCompany.Email)
	require.NotEqual(company.Phone, updatedCompany.Phone)
	require.Equal(phone, updatedCompany.Phone)
	require.NotEqual(company.UpdatedAt, updatedCompany.UpdatedAt)
	require.WithinDuration(updatedAt, updatedCompany.UpdatedAt, 30*time.Second)
	require.WithinDuration(company.CreatedAt, updatedCompany.CreatedAt, time.Second)
}

func TestUpdateCompanyOnlyName(t *testing.T) {
	require := require.New(t)

	company := createRandomCompany(t)

	name := util.RandomString(6)
	arg := UpdateCompanyParams{
		Name:      sql.NullString{String: name, Valid: true},
		ID:        company.ID,
	}
	updatedCompany, err := testQueries.UpdateCompany(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(updatedCompany)
	require.NotEqual(company.Name, updatedCompany.Name)
	require.Equal(name, updatedCompany.Name)
	require.Equal(company.Email, updatedCompany.Email)
	require.Equal(company.Phone, updatedCompany.Phone)
	require.WithinDuration(time.Now(), updatedCompany.UpdatedAt, time.Minute)
	require.WithinDuration(company.CreatedAt, updatedCompany.CreatedAt, time.Second)
}

func TestUpdateCompanyOnlyEmail(t *testing.T) {
	require := require.New(t)

	company := createRandomCompany(t)

	email := util.RandomEmail()
	arg := UpdateCompanyParams{
		Email:     sql.NullString{String: email, Valid: true},
		ID:        company.ID,
	}
	updatedCompany, err := testQueries.UpdateCompany(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(updatedCompany)
	require.NotEqual(company.Email, updatedCompany.Email)
	require.Equal(email, updatedCompany.Email)
	require.Equal(company.Name, updatedCompany.Name)
	require.Equal(company.Phone, updatedCompany.Phone)
	require.WithinDuration(time.Now(), updatedCompany.UpdatedAt, time.Minute)
	require.WithinDuration(company.CreatedAt, updatedCompany.CreatedAt, time.Second)
}

func TestUpdateCompanyOnlyPhone(t *testing.T) {
	require := require.New(t)

	company := createRandomCompany(t)

	phone := util.RandomPhoneNumber()
	arg := UpdateCompanyParams{
		Phone:     sql.NullString{String: phone, Valid: true},
		ID:        company.ID,
	}
	updatedCompany, err := testQueries.UpdateCompany(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(updatedCompany)
	require.NotEqual(company.Phone, updatedCompany.Phone)
	require.Equal(phone, updatedCompany.Phone)
	require.Equal(company.Name, updatedCompany.Name)
	require.Equal(company.Email, updatedCompany.Email)
	require.WithinDuration(time.Now(), updatedCompany.UpdatedAt, time.Minute)
	require.WithinDuration(company.CreatedAt, updatedCompany.CreatedAt, time.Second)
}

func TestDeleteCompany(t *testing.T) {
	require := require.New(t)

	company := createRandomCompany(t)
	err := testQueries.DeleteCompany(context.Background(), company.ID)
	require.NoError(err)
	company1, err := testQueries.GetCompanyByEmail(context.Background(), company.Email)
	require.Error(err)
	require.EqualError(err, sql.ErrNoRows.Error())
	require.Empty(company1)
}
