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
	db "github.com/lenimbugua/bot/db/sqlc"
	"github.com/lenimbugua/bot/token"
	"github.com/stretchr/testify/require"
)

func TestDeleteBotAPI(t *testing.T) {
	company := randomCompany()
	bot := randomBot(t, company.ID)
	user, _ := randomUser(t, company.ID)

	testcases := []struct {
		name           string
		botID          int64
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub      func(store *mockdb.MockStore)
		checkResponses func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			botID: bot.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetBot(gomock.Any(), gomock.Eq(bot.ID)).Return(bot, nil)
				store.EXPECT().
					DeleteBot(gomock.Any(), gomock.Eq(bot.ID)).
					Times(1)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Invalid ID",
			botID: -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetBot(gomock.Any(), gomock.Eq(bot.ID)).Times(0)
				store.EXPECT().
					DeleteBot(gomock.Any(), gomock.Eq(-1)).
					Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Not found",
			botID: bot.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetBot(gomock.Any(), gomock.Eq(bot.ID)).Times(1).Return(db.Bot{}, sql.ErrNoRows)
				store.EXPECT().DeleteBot(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "internal server error",
			botID: bot.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetBot(gomock.Any(), gomock.Eq(bot.ID)).Times(1).Return(bot, nil)
				store.EXPECT().DeleteBot(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrConnDone)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "GetBotInternalServerError",
			botID: bot.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetBot(gomock.Any(), gomock.Eq(bot.ID)).Times(1).Return(db.Bot{}, sql.ErrConnDone)
				store.EXPECT().DeleteBot(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "NotAuthorizedUser",
			botID: bot.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID+1, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetBot(gomock.Any(), gomock.Eq(bot.ID)).Times(1).Return(bot, nil)
				store.EXPECT().DeleteBot(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "NoAuthorization",
			botID: bot.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetBot(gomock.Any(), gomock.Eq(bot.ID)).Times(0)
				store.EXPECT().DeleteBot(gomock.Any(), gomock.Any()).Times(0)
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

			url := fmt.Sprintf("/bots/%d", testcase.botID)

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require := require.New(t)
			require.NoError(err)
			testcase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testcase.checkResponses(recorder)
		})
	}
}
