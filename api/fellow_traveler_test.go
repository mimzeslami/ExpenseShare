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

func TestCreateFellowTravelerAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	fellowTraveler := randomTraveler(trip)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "BadRequest",
			body: gin.H{
				// "trip_id":           trip.ID,
				"first_name": fellowTraveler.FellowFirstName,
				"last_name":  fellowTraveler.FellowLastName,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateFellowTravelersParams{
					// TripID:          trip.ID,
					FellowFirstName: fellowTraveler.FellowFirstName,
					FellowLastName:  fellowTraveler.FellowLastName,
				}
				store.EXPECT().
					CreateFellowTravelers(gomock.Any(), gomock.Eq(arg)).
					Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "OK",
			body: gin.H{
				"trip_id":    trip.ID,
				"first_name": fellowTraveler.FellowFirstName,
				"last_name":  fellowTraveler.FellowLastName,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateFellowTravelersParams{
					TripID:          trip.ID,
					FellowFirstName: fellowTraveler.FellowFirstName,
					FellowLastName:  fellowTraveler.FellowLastName,
				}
				store.EXPECT().
					CreateFellowTravelers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(fellowTraveler, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchFellowTraveler(t, recorder.Body, fellowTraveler)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"trip_id":    trip.ID,
				"first_name": fellowTraveler.FellowFirstName,
				"last_name":  fellowTraveler.FellowLastName,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateFellowTravelersParams{
					TripID:          trip.ID,
					FellowFirstName: fellowTraveler.FellowFirstName,
					FellowLastName:  fellowTraveler.FellowLastName,
				}
				store.EXPECT().
					CreateFellowTravelers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.FellowTravelers{}, sql.ErrConnDone)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/fellow_travelers"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)

		})
	}

}
func TestGetFellowTravelerAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	fellowTraveler := randomTraveler(trip)

	testCases := []struct {
		name             string
		fellowTravelerID int64
		setupAuth        func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs       func(store *mockdb.MockStore)
		checkResponse    func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:             "OK",
			fellowTravelerID: fellowTraveler.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetFellowTravelerParams{
					ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetFellowTraveler(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(fellowTraveler, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchFellowTraveler(t, recorder.Body, fellowTraveler)
			},
		},
		{
			name:             "UnauthorizedUser",
			fellowTravelerID: fellowTraveler.ID,

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetFellowTraveler(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, -time.Minute)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BadRequestBody",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetFellowTravelerParams{
					// ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetFellowTraveler(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:             "FellowTravelerNotFound",
			fellowTravelerID: fellowTraveler.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetFellowTravelerParams{
					ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetFellowTraveler(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.FellowTravelers{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:             "InternalError",
			fellowTravelerID: fellowTraveler.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetFellowTravelerParams{
					ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetFellowTraveler(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.FellowTravelers{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			url := fmt.Sprintf("/fellow_travelers/%d", tc.fellowTravelerID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestListFellowTravelerAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	n := 5
	fellowTravelers := make([]db.FellowTravelers, n)
	for i := 0; i < n; i++ {
		fellowTravelers[i] = randomTraveler(trip)
	}

	testCases := []struct {
		name          string
		tripID        int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			tripID: trip.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripFellowTravelersParams{
					TripID: trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTripFellowTravelers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(fellowTravelers, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchFellowTravelers(t, recorder.Body, fellowTravelers)
			},
		},
		{
			name:   "UnauthorizedUser",
			tripID: trip.ID,

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripFellowTravelersParams{
					TripID: trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTripFellowTravelers(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name:   "InternalError",
			tripID: trip.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripFellowTravelersParams{
					TripID: trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTripFellowTravelers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.FellowTravelers{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			// tripID: trip.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripFellowTravelersParams{
					// TripID: trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTripFellowTravelers(gomock.Any(), gomock.Eq(arg)).
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

			url := fmt.Sprintf("/trips/fellow_travelers/%d", tc.tripID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}
}

func TestUpdateFellowTravelerAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	fellowTraveler := randomTraveler(trip)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"id":         fellowTraveler.ID,
				"first_name": "NewFirstName",
				"last_name":  "NewLastName",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFellowTravelerParams{
					ID:              fellowTraveler.ID,
					FellowFirstName: "NewFirstName",
					FellowLastName:  "NewLastName",
				}
				store.EXPECT().
					UpdateFellowTraveler(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(fellowTraveler, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchFellowTraveler(t, recorder.Body, fellowTraveler)
			},
		},
		{
			name: "UnauthorizedUser",
			body: gin.H{
				"id":         fellowTraveler.ID,
				"first_name": "NewFirstName",
				"last_name":  "NewLastName",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFellowTravelerParams{
					ID:              fellowTraveler.ID,
					FellowFirstName: "NewFirstName",
					FellowLastName:  "NewLastName",
				}
				store.EXPECT().
					UpdateFellowTraveler(gomock.All(), gomock.Eq(arg)).
					Times(0).
					Return(db.FellowTravelers{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BadRequestBody",
			body: gin.H{
				// "id":         fellowTraveler.ID,
				"first_name": "NewFirstName",
				"last_name":  "NewLastName",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFellowTravelerParams{
					ID:              fellowTraveler.ID,
					FellowFirstName: "NewFirstName",
					FellowLastName:  "NewLastName",
				}
				store.EXPECT().
					UpdateFellowTraveler(gomock.All(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)

			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id":         fellowTraveler.ID,
				"first_name": "NewFirstName",
				"last_name":  "NewLastName",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFellowTravelerParams{
					ID:              fellowTraveler.ID,
					FellowFirstName: "NewFirstName",
					FellowLastName:  "NewLastName",
				}
				store.EXPECT().
					UpdateFellowTraveler(gomock.All(), gomock.Eq(arg)).
					Times(1).
					Return(db.FellowTravelers{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/fellow_travelers"

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)

		})
	}
}

func TestDeleteFellowTravelerAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	fellowTraveler := randomTraveler(trip)

	testCases := []struct {
		name             string
		fellowTravelerID int64
		setupAuth        func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs       func(store *mockdb.MockStore)
		checkResponse    func(recorder *httptest.ResponseRecorder)
	}{

		{
			name:             "OK",
			fellowTravelerID: fellowTraveler.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				getFellowArg := db.GetFellowTravelerParams{
					ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				getTripArg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().GetFellowTraveler(gomock.Any(), gomock.Eq(getFellowArg)).
					Times(1).
					Return(fellowTraveler, nil)

				store.EXPECT().GetTrip(gomock.Any(), gomock.Eq(getTripArg)).
					Times(1).
					Return(trip, nil)

				store.EXPECT().
					DeleteFellowTraveler(gomock.Any(), gomock.Eq(fellowTraveler.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Empty(t, recorder.Body.Bytes())
			},
		},
		{
			name: "BadRequestUrl",
			// fellowTravelerID: fellowTraveler.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				getFellowArg := db.GetFellowTravelerParams{
					ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				getTripArg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().GetFellowTraveler(gomock.Any(), gomock.Eq(getFellowArg)).
					Times(0)

				store.EXPECT().GetTrip(gomock.Any(), gomock.Eq(getTripArg)).
					Times(0)

				store.EXPECT().
					DeleteFellowTraveler(gomock.Any(), gomock.Eq(fellowTraveler.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:             "FrellowTravelerNotFound",
			fellowTravelerID: fellowTraveler.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				getFellowArg := db.GetFellowTravelerParams{
					ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				getTripArg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().GetFellowTraveler(gomock.Any(), gomock.Eq(getFellowArg)).
					Times(1).
					Return(db.FellowTravelers{}, sql.ErrNoRows)

				store.EXPECT().GetTrip(gomock.Any(), gomock.Eq(getTripArg)).
					Times(0)

				store.EXPECT().
					DeleteFellowTraveler(gomock.Any(), gomock.Eq(fellowTraveler.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:             "tripNotFound",
			fellowTravelerID: fellowTraveler.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				getFellowArg := db.GetFellowTravelerParams{
					ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				getTripArg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().GetFellowTraveler(gomock.Any(), gomock.Eq(getFellowArg)).
					Times(1).
					Return(fellowTraveler, nil)

				store.EXPECT().GetTrip(gomock.Any(), gomock.Eq(getTripArg)).
					Times(1).
					Return(db.Trips{}, sql.ErrNoRows)

				store.EXPECT().
					DeleteFellowTraveler(gomock.Any(), gomock.Eq(fellowTraveler.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:             "InternalError",
			fellowTravelerID: fellowTraveler.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, trip.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				getFellowArg := db.GetFellowTravelerParams{
					ID:     fellowTraveler.ID,
					UserID: trip.UserID,
				}
				getTripArg := db.GetTripParams{
					ID:     trip.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().GetFellowTraveler(gomock.Any(), gomock.Eq(getFellowArg)).
					Times(1).
					Return(fellowTraveler, nil)

				store.EXPECT().GetTrip(gomock.Any(), gomock.Eq(getTripArg)).
					Times(1).
					Return(trip, nil)

				store.EXPECT().
					DeleteFellowTraveler(gomock.Any(), gomock.Eq(fellowTraveler.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			url := fmt.Sprintf("/fellow_travelers/%d", tc.fellowTravelerID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)

		})
	}

}

func randomTraveler(trip db.Trips) db.FellowTravelers {
	return db.FellowTravelers{
		ID:              int64(util.RandomInt(1, 1000)),
		TripID:          trip.ID,
		FellowFirstName: util.RandomString(6),
		FellowLastName:  util.RandomString(6),
	}
}

func requireBodyMatchFellowTraveler(t *testing.T, body *bytes.Buffer, traveler db.FellowTravelers) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotFellowTraveler db.FellowTravelers
	err = json.Unmarshal(data, &gotFellowTraveler)
	require.NoError(t, err)
	require.Equal(t, traveler, gotFellowTraveler)

}

func requireBodyMatchFellowTravelers(t *testing.T, body *bytes.Buffer, travelers []db.FellowTravelers) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotFellowTravelers []db.FellowTravelers
	err = json.Unmarshal(data, &gotFellowTravelers)
	require.NoError(t, err)
	require.Equal(t, travelers, gotFellowTravelers)

}
