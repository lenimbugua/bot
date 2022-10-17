package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/lenimbugua/bot/db/mock"
	db "github.com/lenimbugua/bot/db/sqlc"
	"github.com/lenimbugua/bot/token"
	"github.com/stretchr/testify/require"
)

func TestUpdateCompanyAPI(t *testing.T) {
	user, _ := randomUser(t)
	company := randomCompany()
	updatedAt := time.Now()
	testcases := []struct {
		name          string
		id            int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			id:   company.ID,
			body: gin.H{
				"name":       company.Name,
				"phone":      company.Phone,
				"email":      company.Email,
				"updated_at": updatedAt,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.UpdateCompanyParams{
					Name:      sql.NullString{String: company.Name, Valid: true},
					Email:     sql.NullString{String: company.Email, Valid: true},
					Phone:     sql.NullString{String: company.Phone, Valid: true},
					UpdatedAt: sql.NullTime{Time: updatedAt, Valid: false},
					ID:        company.ID,
				}
				store.EXPECT().
					UpdateCompany(gomock.Any(), EqUpdateCompanyParams(arg)).
					Times(1).
					Return(company, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCompany(t, recorder.Body, company)

			},
		},
		{
			name: "Bad Request",
			id:   company.ID,
			body: gin.H{
				"name":       company.Name,
				"phone":      company.Phone,
				"email":      company.Email,
				"updated_at": updatedAt,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request (Bad ID)",
			id:   -1,
			body: gin.H{
				"name":       company.Name,
				"phone":      company.Phone,
				"email":      company.Email,
				"updated_at": updatedAt,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "Not Found",
			id:   1,
			body: gin.H{
				"name":       company.Name,
				"phone":      company.Phone,
				"email":      company.Email,
				"updated_at": updatedAt,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Company{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: gin.H{
				"name":       company.Name,
				"phone":      company.Phone,
				"email":      company.Email,
				"updated_at": updatedAt,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Company{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testcases {
		testCase := testcases[i]

		t.Run(testCase.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStub(store)
			//start server and send requests
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/companies/%d", testCase.id)

			require := require.New(t)
			argJSON, err := json.Marshal(testCase.body)
			require.NoError(err)

			bytesArg := bytes.NewBuffer(argJSON)

			request, err := http.NewRequest(http.MethodPut, url, bytesArg)
			require.NoError(err)
			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

type eqUpdateCompanyParamMatcher struct {
	arg db.UpdateCompanyParams
}

func (e eqUpdateCompanyParamMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateCompanyParams)
	if !ok {
		return false
	}
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateCompanyParamMatcher) String() string {
	return fmt.Sprintf("Matches arg %v", e.arg)
}

func EqUpdateCompanyParams(arg db.UpdateCompanyParams) gomock.Matcher {
	return eqUpdateCompanyParamMatcher{arg}
}


