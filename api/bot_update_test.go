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

type eqUpdateBotParamMatcher struct {
	arg db.UpdateBotParams
}

func (e eqUpdateBotParamMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateBotParams)
	if !ok {
		return false
	}
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateBotParamMatcher) String() string {
	return fmt.Sprintf("Matches arg %v", e.arg)
}

func EqUpdateBotParams(arg db.UpdateBotParams) gomock.Matcher {
	return eqUpdateBotParamMatcher{arg}
}

func TestUpdateBot(t *testing.T) {
	company := randomCompany()
	user, _ := randomUser(t, company.ID)
	bot := randomBot(t, company.ID)
	testCases := []struct {
		name          string
		botID         int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{

		//OK
		{
			name:  "OK",
			botID: bot.ID,
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateBotParams{
					Title:     sql.NullString{String: bot.Title, Valid: true},
					CompanyID: sql.NullInt64{Int64: bot.CompanyID, Valid: true},
					ID:        bot.ID,
				}

				store.EXPECT().
					UpdateBot(gomock.Any(), EqUpdateBotParams(arg)).
					Times(1).
					Return(bot, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBot(t, recorder.Body, bot)
			},
		},
		//Unauthorized user
		{
			name:  "UnauthorizedUser",
			botID: bot.ID,
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID+1, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBot(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		//No Authorisation
		{
			name:  "NoAuthorization",
			botID: bot.ID,
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateBot(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		//NotFound
		{
			name:  "NotFound",
			botID: bot.ID,
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateBot(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Bot{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		//InvalidID
		{
			name:  "InvalidID",
			botID: -1,
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateBot(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//EmptyTitle
		{
			name:  "EmptyTitle",
			botID: bot.ID,
			body: gin.H{
				"title":      "",
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateBot(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//InvalidCompanyID
		{
			name:  "InvalidCompanyID",
			botID: bot.ID,
			body: gin.H{
				"title":      bot.Title,
				"company_id": -1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateBot(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//InternalServerError
		{
			name:  "InternalServerError",
			botID: bot.ID,
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateBot(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Bot{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
			url := fmt.Sprintf("/bots/%d", tc.botID)

			require := require.New(t)
			data, err := json.Marshal(tc.body)
			require.NoError(err)

			byteData := bytes.NewBuffer(data)

			request, err := http.NewRequest(http.MethodPut, url, byteData)
			require.NoError(err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
