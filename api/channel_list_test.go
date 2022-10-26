package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/lenimbugua/bot/db/mock"
	db "github.com/lenimbugua/bot/db/sqlc"
	"github.com/lenimbugua/bot/token"
	"github.com/lenimbugua/bot/util"
	"github.com/stretchr/testify/require"
)

func requireBodyMatchListChannels(t *testing.T, body *bytes.Buffer, channels []db.Channel) {
	require := require.New(t)
	data, err := ioutil.ReadAll(body)
	require.NoError(err)

	var gotChannels []db.Channel
	err = json.Unmarshal(data, &gotChannels)
	require.NoError(err)

	require.Equal(channels, gotChannels)
}

func TestListChannelsAPI(t *testing.T) {
	user, _ := randomUser(t, util.RandInt(1, 100))
	n := 5
	channels := make([]db.Channel, n)
	for i := 0; i < n; i++ {
		channels[i] = randomChannel()
	}
	type Query struct {
		pageID   int
		pageSize int
	}

	testcases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.ListChannelsParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListChannels(gomock.Any(), arg).
					Times(1).
					Return(channels, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchListChannels(t, recorder.Body, channels)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListChannels(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Channel{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid Page Size",
			query: Query{
				pageID:   1,
				pageSize: 1000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListChannels(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Page ID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Phone, user.ID, user.Name, user.CompanyID, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListChannels(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoAuthorizationHeader",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListChannels(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			require := require.New(t)

			url := "/list/channels"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", testcase.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", testcase.query.pageSize))

			request.URL.RawQuery = q.Encode()

			testcase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testcase.checkResponse(recorder)
		})
	}
}
