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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/mimzeslami/expense_share/db/mock"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func TestGetTripAPI(t *testing.T) {
	user, _ := randomUser(t)

	trip := randomTrip(user)

	testCases := []struct {
		name          string
		tripID        int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(trip, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTrip(t, recorder.Body, trip)
			},
		},
		{
			name:   "NotFound",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Trips{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Trips{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			tripID: 0,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "UnauthorizedUser",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, -time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			crtl := gomock.NewController(t)
			defer crtl.Finish()

			store := mockdb.NewMockStore(crtl)

			tc.buildStubs(store)

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/trips/%d", tc.tripID)

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func TestCreateTripAPI(t *testing.T) {
	user, _ := randomUser(t)

	trip := randomTrip(user)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":      trip.Title,
				"start_date": trip.StartDate,
				"end_date":   trip.EndDate,
				"user_id":    trip.UserID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTripParams{
					Title:     trip.Title,
					StartDate: trip.StartDate,
					EndDate:   trip.EndDate,
					UserID:    trip.UserID,
				}

				store.EXPECT().
					CreateTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(trip, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTrip(t, recorder.Body, trip)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"title":      trip.Title,
				"start_date": trip.StartDate,
				"end_date":   trip.EndDate,
				"user_id":    trip.UserID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTripParams{
					Title:     trip.Title,
					StartDate: trip.StartDate,
					EndDate:   trip.EndDate,
					UserID:    trip.UserID,
				}

				store.EXPECT().
					CreateTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Trips{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequestBody",
			body: gin.H{
				"start_date": trip.StartDate,
				"end_date":   trip.EndDate,
				"user_id":    trip.UserID,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateTrip(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "UnauthorizedUser",
			body: gin.H{
				"title":      trip.Title,
				"start_date": trip.StartDate,
				"end_date":   trip.EndDate,
				"user_id":    trip.UserID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTripParams{
					Title:     trip.Title,
					StartDate: trip.StartDate,
					EndDate:   trip.EndDate,
					UserID:    trip.UserID,
				}

				store.EXPECT().
					CreateTrip(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/trips"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}

func TestListUserTripsAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	trips := make([]db.Trips, n)
	for i := 0; i < n; i++ {
		trips[i] = randomTrip(user)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTripParams{
					UserID: user.ID,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(trips, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTrip(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTrip(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Trips{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTrip(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTrip(gomock.Any(), gomock.Any()).
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

			url := "/trips"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("offset", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("limit", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateTripAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	testCases := []struct {
		name          string
		tripID        int64
		update        gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			tripID: trip.ID,
			update: gin.H{
				"id":         trip.ID,
				"title":      trip.Title,
				"start_date": trip.StartDate,
				"end_date":   trip.EndDate,
				"user_id":    trip.UserID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTripParams{
					ID:        trip.ID,
					Title:     trip.Title,
					StartDate: trip.StartDate,
					EndDate:   trip.EndDate,
					UserID:    trip.UserID,
				}
				store.EXPECT().
					UpdateTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(trip, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTrip(t, recorder.Body, trip)
			},
		},
		{
			name:   "InternalError",
			tripID: trip.ID,
			update: gin.H{
				"id":         trip.ID,
				"title":      trip.Title,
				"start_date": trip.StartDate,
				"end_date":   trip.EndDate,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTripParams{
					ID:        trip.ID,
					Title:     trip.Title,
					StartDate: trip.StartDate,
					EndDate:   trip.EndDate,
					UserID:    trip.UserID,
				}
				store.EXPECT().
					UpdateTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Trips{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "BadRequest",
			tripID: trip.ID,
			update: gin.H{
				"id":    trip.ID,
				"title": trip.Title,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateTrip(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

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

			data, err := json.Marshal(tc.update)
			require.NoError(t, err)

			url := "/trips"

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}

func TestDeleteTripAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	testCases := []struct {
		name          string
		tripID        int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					DeleteTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					DeleteTrip(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "UnauthorizedUser",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					DeleteTrip(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "BadRequest",
			tripID: 0,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					DeleteTrip(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

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

			url := fmt.Sprintf("/trips/%d", tc.tripID)

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}
func randomTrip(user db.Users) db.Trips {

	return db.Trips{
		ID:        int64(util.RandomInt(1, 1000)),
		Title:     util.RandomString(6),
		StartDate: util.RandomDatetime(),
		EndDate:   util.RandomDatetime(),
		UserID:    user.ID,
	}
}

func requireBodyMatchTrip(t *testing.T, body *bytes.Buffer, trip db.Trips) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTrip db.Trips

	err = json.Unmarshal(data, &gotTrip)
	require.NoError(t, err)

	require.Equal(t, trip.ID, gotTrip.ID)
	require.Equal(t, trip.Title, gotTrip.Title)
	require.Equal(t, trip.UserID, gotTrip.UserID)

}
