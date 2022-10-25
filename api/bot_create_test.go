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

func randomBot(t *testing.T, companyID int64) (bot db.Bot) {

	bot = db.Bot{
		Title:     util.RandomString(6),
		CompanyID: companyID,
		ID:        util.RandInt(1, 100),
	}
	return
}

func requireBodyMatchBot(t *testing.T, body *bytes.Buffer, bot db.Bot) {
	require := require.New(t)
	data, err := ioutil.ReadAll(body)
	require.NoError(err)

	var gotBot db.Bot
	err = json.Unmarshal(data, &gotBot)
	require.NoError(err)
	require.Equal(gotBot, bot)
}

func TestCreateBotAPI(t *testing.T) {
	company := randomCompany()
	bot := randomBot(t, company.ID)

	user, _ := randomUser(t, company.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		//Ok (Happy Path)
		{
			name: "OK",
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateBotParams{
					Title:     bot.Title,
					CompanyID: bot.CompanyID,
				}
				store.EXPECT().GetCompanyByID(gomock.Any(), gomock.Eq(company.ID)).Times(1).Return(company, nil)
				store.EXPECT().CreateBot(gomock.Any(), gomock.Eq(arg)).Times(1).Return(bot, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBot(t, recorder.Body, bot)

			},
		},
		//No authorisation
		{
			name: "NoAuthorization",
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBot(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		//CompanyWithIDSpecifiedNotFound
		{
			name: "CompanyNotFound",
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			}, buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCompanyByID(gomock.Any(), gomock.Any()).Times(1).Return(db.Company{}, sql.ErrNoRows)
				store.EXPECT().
					CreateBot(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		//Internal Server Error
		{
			name: "InternalError",
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCompanyByID(gomock.Any(), gomock.Eq(company.ID)).Times(1).Return(company, nil)
				store.EXPECT().
					CreateBot(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Bot{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		//Internal Server Error When getting Company
		{
			name: "GetCompanyInternalError",
			body: gin.H{
				"title":      bot.Title,
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCompanyByID(gomock.Any(), gomock.Any()).Times(1).Return(db.Company{}, sql.ErrConnDone)
				store.EXPECT().
					CreateBot(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		//Invalid Company ID
		{
			name: "InvalidCompanyID",
			body: gin.H{
				"title":      bot.Title,
				"company_id": -1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBot(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//Invalid Title
		{
			name: "InvalidTitle",
			body: gin.H{
				"title":      "",
				"company_id": bot.CompanyID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBot(gomock.Any(), gomock.Any()).
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

			url := "/bots"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
