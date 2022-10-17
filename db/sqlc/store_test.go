package db

// import (
// 	"context"
// 	"testing"

// 	"github.com/lenimbugua/bot/util"
// 	"github.com/stretchr/testify/require"
// )

// func TestUserCompanyTx(t *testing.T) {
// 	store := NewSQLStore(testDB)

// 	user := createRandomUser(t)

// 	//run n concurrent company user transactions
// 	n := 5
// 	errs := make(chan error)
// 	results := make(chan UserCompanyTxResult)

// 	for i := 0; i < n; i++ {
// 		go func() {
// 			result, err := store.UserCompanyTx(context.Background(), UserCompanyTxParams{
// 				UserID:       user.ID,
// 				CompanyPhone: util.RandomPhoneNumber(),
// 				CompanyName:  util.RandomString(6),
// 				CompanyEmail: util.RandomEmail(),
// 			})

// 			errs <- err
// 			results <- result
// 		}()
// 	}

// 	for i := 0; i < n; i++ {
// 		require := require.New(t)
// 		err := <-errs
// 		require.NoError(err)

// 		result := <-results
// 		require.NotEmpty(result)

// 		userCompany := result.UserCompany
// 		company := result.Company

// 		require.Equal(userCompany.UserID, user.ID)
// 		require.Equal(userCompany.CompanyID, company.ID)
// 	}
// }
