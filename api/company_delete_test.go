package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/lenimbugua/bot/db/mock"
	"github.com/lenimbugua/bot/token"
	"github.com/stretchr/testify/require"
)

func TestDeleteCompanyAPI(t *testing.T) {
	company := randomCompany()
	user, _ := randomUser(t, company.ID)

	testcases := []struct {
		name           string
		companyID      int64
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub      func(store *mockdb.MockStore)
		checkResponses func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			companyID: company.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Invalid ID",
			companyID: -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCompany(gomock.Any(), gomock.Eq(-1)).
					Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Not found",
			companyID: company.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCompany(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrNoRows)

			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "internal server error",
			companyID: company.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCompany(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrConnDone)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			companyID: company.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCompany(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testcases {
		testcase := testcases[i]
		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testcase.buildStub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/companies/%d", testcase.companyID)

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require := require.New(t)
			require.NoError(err)
			testcase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testcase.checkResponses(recorder)
		})
	}
}
