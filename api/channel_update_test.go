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

type eqUpdateChannelParamMatcher struct {
	arg db.UpdateChannelParams
}

func (e eqUpdateChannelParamMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateChannelParams)
	if !ok {
		return false
	}
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateChannelParamMatcher) String() string {
	return fmt.Sprintf("Matches arg %v", e.arg)
}

func EqUpdateChannelParams(arg db.UpdateChannelParams) gomock.Matcher {
	return eqUpdateChannelParamMatcher{arg}
}

func TestUpdateChannel(t *testing.T) {
	company := randomCompany()
	user, _ := randomUser(t, company.ID)
	channel := randomChannel()
	testCases := []struct {
		name          string
		channelID     int32
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{

		//OK
		{
			name:      "OK",
			channelID: channel.ID,
			body: gin.H{
				"name": channel.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateChannelParams{
					Name: sql.NullString{String: channel.Name, Valid: true},
					ID:   channel.ID,
				}

				store.EXPECT().
					UpdateChannel(gomock.Any(), EqUpdateChannelParams(arg)).
					Times(1).
					Return(channel, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchChannel(t, recorder.Body, channel)
			},
		},

		//No Authorisation
		{
			name:      "NoAuthorization",
			channelID: channel.ID,
			body: gin.H{
				"name": channel.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateChannel(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		//NotFound
		{
			name:      "NotFound",
			channelID: channel.ID,
			body: gin.H{
				"name": channel.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateChannel(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Channel{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		//InvalidID
		{
			name:      "InvalidID",
			channelID: -1,
			body: gin.H{
				"name": channel.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateChannel(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//Emptyname
		{
			name:      "EmptyName",
			channelID: channel.ID,
			body: gin.H{
				"name": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateChannel(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//InternalServerError
		{
			name:      "InternalServerError",
			channelID: channel.ID,
			body: gin.H{
				"name": channel.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateChannel(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Channel{}, sql.ErrConnDone)
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
			url := fmt.Sprintf("/channels/%d", tc.channelID)

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
