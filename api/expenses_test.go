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

func TestCreateExpenseAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	fellowTraveler := randomTraveler(trip)
	expense := randomExpanse(fellowTraveler)

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
				"trip_id":           expense.TripID,
				"payer_traveler_id": expense.PayerTravelerID,
				"amount":            expense.Amount,
				"description":       expense.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateExpenseParams{
					TripID:          expense.TripID,
					PayerTravelerID: expense.PayerTravelerID,
					Amount:          expense.Amount,
					Description:     expense.Description,
				}
				store.EXPECT().
					CreateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(expense, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchExpense(t, recorder.Body, expense)
			},
		},
		{
			name: "BadRequest",
			//Dont sending tripId
			body: gin.H{
				"payer_traveler_id": expense.PayerTravelerID,
				"amount":            expense.Amount,
				"description":       expense.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateExpenseParams{
					PayerTravelerID: expense.PayerTravelerID,
					Amount:          expense.Amount,
					Description:     expense.Description,
				}
				store.EXPECT().
					CreateExpense(gomock.Any(), gomock.Eq(arg)).
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
				"trip_id":           expense.TripID,
				"payer_traveler_id": expense.PayerTravelerID,
				"amount":            expense.Amount,
				"description":       expense.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateExpenseParams{
					TripID:          expense.TripID,
					PayerTravelerID: expense.PayerTravelerID,
					Amount:          expense.Amount,
					Description:     expense.Description,
				}
				store.EXPECT().
					CreateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Expenses{}, sql.ErrConnDone)
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

			url := "/expenses"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestGetExpense(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	fellowTraveler := randomTraveler(trip)
	expense := randomExpanse(fellowTraveler)

	testCases := []struct {
		name          string
		expenseID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "BadRequest",
			// expenseID: expense.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetExpenseParams{
					ID:     expense.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetExpense(gomock.Any(), gomock.Eq(arg)).
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
			name:      "Ok",
			expenseID: expense.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetExpenseParams{
					ID:     expense.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetExpense(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(expense, nil)

			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			expenseID: expense.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetExpenseParams{
					ID:     expense.ID,
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetExpense(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Expenses{}, sql.ErrConnDone)

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

			url := fmt.Sprintf("/expenses/%d", tc.expenseID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})

	}
}

func TestGetTripExpensesAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	fellowTraveler1 := randomTraveler(trip)
	fellowTraveler2 := randomTraveler(trip)
	expanse1 := randomExpanse(fellowTraveler1)
	expanse2 := randomExpanse(fellowTraveler2)
	expanse3 := randomExpanse(fellowTraveler2)

	expanses := []db.Expenses{expanse1, expanse2, expanse3}

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
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.GetTripExpensesParams{
					TripID: trip.ID,
					UserID: trip.UserID,
				}

				store.EXPECT().
					GetTripExpenses(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(expanses, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)

			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			//Dont sending tripId
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTripExpensesParams{
					UserID: trip.UserID,
				}
				store.EXPECT().
					GetTripExpenses(gomock.Any(), gomock.Eq(arg)).
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
			name:   "InteranlError",
			tripID: trip.ID,
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.GetTripExpensesParams{
					TripID: trip.ID,
					UserID: trip.UserID,
				}

				store.EXPECT().
					GetTripExpenses(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Expenses{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/trip/expenses/%d", tc.tripID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})

	}
}

func TestUpdateExpenseAPI(t *testing.T) {
	user, _ := randomUser(t)
	trip := randomTrip(user)
	fellowTraveler := randomTraveler(trip)
	expense := randomExpanse(fellowTraveler)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{

		{
			name: "BadRequest",
			//Dont sending id
			body: gin.H{
				// "id":                expense.ID,
				"payer_traveler_id": expense.PayerTravelerID,
				"amount":            expense.Amount,
				"description":       expense.Description,
				"trip_id":           expense.TripID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateExpenseParams{
					// ID:              expense.ID,
					PayerTravelerID: expense.PayerTravelerID,
					Amount:          expense.Amount,
					Description:     expense.Description,
					TripID:          expense.TripID,
				}
				store.EXPECT().
					UpdateExpense(gomock.Any(), gomock.Eq(arg)).
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
			name: "OK",
			body: gin.H{
				"id":                expense.ID,
				"payer_traveler_id": expense.PayerTravelerID,
				"amount":            expense.Amount,
				"description":       expense.Description,
				"trip_id":           expense.TripID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateExpenseParams{
					ID:              expense.ID,
					PayerTravelerID: expense.PayerTravelerID,
					Amount:          expense.Amount,
					Description:     expense.Description,
					TripID:          expense.TripID,
				}
				store.EXPECT().
					UpdateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(expense, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, trip.UserID, time.Minute)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchExpense(t, recorder.Body, expense)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id":                expense.ID,
				"payer_traveler_id": expense.PayerTravelerID,
				"amount":            expense.Amount,
				"description":       expense.Description,
				"trip_id":           expense.TripID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateExpenseParams{
					ID:              expense.ID,
					PayerTravelerID: expense.PayerTravelerID,
					Amount:          expense.Amount,
					Description:     expense.Description,
					TripID:          expense.TripID,
				}
				store.EXPECT().
					UpdateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Expenses{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/expenses")
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})

	}
}

func randomExpanse(traveler db.FellowTravelers) db.Expenses {
	return db.Expenses{
		ID:              int64(util.RandomInt(0, 1000)),
		TripID:          traveler.TripID,
		PayerTravelerID: traveler.ID,
		Amount:          util.RandomMoney(),
		Description:     util.RandomString(20),
	}
}

func requireBodyMatchExpense(t *testing.T, body *bytes.Buffer, expense db.Expenses) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotExpense db.Expenses
	err = json.Unmarshal(data, &gotExpense)
	require.NoError(t, err)
	require.Equal(t, expense, gotExpense)
}
