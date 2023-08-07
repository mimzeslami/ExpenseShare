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
	mockdb "github.com/mimzeslami/expense_share/db/mock"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func TestGetTripAPI(t *testing.T) {

	trip := randomTrip()

	testCases := []struct {
		name          string
		tripID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Eq(trip.ID)).
					Times(1).
					Return(trip, nil)
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
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Eq(trip.ID)).
					Times(1).
					Return(db.Trips{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTrip(gomock.Any(), gomock.Eq(trip.ID)).
					Times(1).
					Return(db.Trips{}, sql.ErrConnDone)
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
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func randomTrip() db.Trips {

	return db.Trips{
		ID:        int64(util.RandomInt(1, 1000)),
		Title:     util.RandomString(6),
		StartDate: util.RandomDatetime(),
		EndDate:   util.RandomDatetime(),
		UserID:    int64(util.RandomInt(1, 1000)),
		CreatedAt: util.RandomDatetime(),
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
	require.WithinDuration(t, trip.StartDate, gotTrip.StartDate, time.Second)
	require.WithinDuration(t, trip.EndDate, gotTrip.EndDate, time.Second)

}
