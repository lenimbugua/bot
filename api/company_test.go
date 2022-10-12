package api

// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	mockdb "github.com/lenimbugua/bot/db/mock"
// 	db "github.com/lenimbugua/bot/db/sqlc"
// 	"github.com/lenimbugua/bot/token"
// 	"github.com/lenimbugua/bot/util"

// 	"github.com/stretchr/testify/require"
// )

// func randomCompany() db.Company {
// 	return db.Company{
// 		Name:  util.RandomString(6),
// 		Email: util.RandomEmail(),
// 		Phone: util.RandomPhoneNumber(),
// 	}
// }

// func requireBodyMatchCompany(t *testing.T, body *bytes.Buffer, company db.Company) {
// 	require := require.New(t)
// 	data, err := ioutil.ReadAll(body)
// 	require.NoError(err)

// 	var gotCompany db.Company
// 	err = json.Unmarshal(data, &gotCompany)
// 	require.NoError(err)
// 	require.Equal(gotCompany, company)
// }

// func TestCreateCompanyAPI(t *testing.T) {
// 	company := randomCompany()

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"currency": account.Currency,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.CreateAccountParams{
// 					Owner:    account.Owner,
// 					Currency: account.Currency,
// 					Balance:  0,
// 				}

// 				store.EXPECT().
// 					CreateAccount(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					Return(account, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchAccount(t, recorder.Body, account)
// 			},
// 		},
// 		{
// 			name: "NoAuthorization",
// 			body: gin.H{
// 				"currency": account.Currency,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAccount(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			body: gin.H{
// 				"currency": account.Currency,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAccount(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.Account{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidCurrency",
// 			body: gin.H{
// 				"currency": "invalid",
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAccount(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/accounts"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			tc.setupAuth(t, request, server.tokenMaker)
// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }
