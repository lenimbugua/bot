package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/lenimbugua/bot/db/mock"
	db "github.com/lenimbugua/bot/db/sqlc"
	"github.com/lenimbugua/bot/token"
	"github.com/lenimbugua/bot/util"

	"github.com/stretchr/testify/require"
)

func randomCompany() db.Company {
	return db.Company{
		Name:  util.RandomString(6),
		Email: util.RandomEmail(),
		Phone: util.RandomPhoneNumber(),
		ID:    util.RandInt(1, 1000),
	}
}

func requireBodyMatchCompany(t *testing.T, body *bytes.Buffer, company db.Company) {
	require := require.New(t)
	data, err := ioutil.ReadAll(body)
	require.NoError(err)

	var gotCompany db.Company
	err = json.Unmarshal(data, &gotCompany)
	require.NoError(err)
	require.Equal(gotCompany, company)
}

func requireBodyMatchListCompanies(t *testing.T, body *bytes.Buffer, companies []db.Company) {
	require := require.New(t)
	data, err := ioutil.ReadAll(body)
	require.NoError(err)

	var gotCompanies []db.Company
	err = json.Unmarshal(data, &gotCompanies)
	require.NoError(err)

	require.Equal(companies, gotCompanies)
}

func TestCreateCompanyAPI(t *testing.T) {
	company := randomCompany()
	user, _ := randomUser(t, company.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"phone": company.Phone,
				"name":  company.Name,
				"email": company.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCompanyParams{
					Phone: company.Phone,
					Name:  company.Name,
					Email: company.Email,
				}
				store.EXPECT().CreateCompany(gomock.Any(), gomock.Eq(arg)).Times(1).Return(company, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCompany(t, recorder.Body, company)

			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"name":  company.Name,
				"email": company.Email,
				"Phone": company.Phone,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name":  company.Name,
				"email": company.Email,
				"Phone": company.Phone,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCompany(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Company{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"email": "invalid",
				"name":  company.Name,
				"Phone": company.Phone,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPhone",
			body: gin.H{
				"email": company.Email,
				"name":  company.Name,
				"Phone": "invalid",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/companies"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
