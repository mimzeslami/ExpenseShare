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
	mockdb "github.com/mimzeslami/expense_share/db/mock"
	db "github.com/mimzeslami/expense_share/db/sqlc"
	"github.com/mimzeslami/expense_share/token"
	"github.com/mimzeslami/expense_share/util"
	"github.com/stretchr/testify/require"
)

func TestCreateExpenseAPI(t *testing.T) {

	user, _ := randomUser(t)
	groupCategory := randomGroupCategory()
	group := randomGroup(t, user, groupCategory)
	group_member := randomGroupMember(t, user, group)
	expense := randomExpense(group_member)

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
				"group_id":    group_member.GroupID,
				"paid_by_id":  group_member.UserID,
				"amount":      expense.Amount,
				"description": expense.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateExpenseParams{
					GroupID:     group_member.GroupID,
					PaidByID:    group_member.UserID,
					Amount:      expense.Amount,
					Description: expense.Description,
				}

				arg1 := db.GetGroupMemberByGroupIDAndUserIDParams{
					GroupID: group_member.GroupID,
					UserID:  group_member.UserID,
				}

				store.EXPECT().
					GetGroupMemberByGroupIDAndUserID(gomock.Any(), gomock.Eq(arg1)).
					Times(1).
					Return(group_member, nil)

				store.EXPECT().
					CreateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(expense, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchExpense(t, recorder.Body, expense)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"group_id":    group_member.GroupID,
				"paid_by_id":  group_member.UserID,
				"amount":      expense.Amount,
				"description": expense.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateExpenseParams{
					GroupID:     group_member.GroupID,
					PaidByID:    group_member.UserID,
					Amount:      expense.Amount,
					Description: expense.Description,
				}

				arg1 := db.GetGroupMemberByGroupIDAndUserIDParams{
					GroupID: expense.GroupID,
					UserID:  user.ID,
				}

				store.EXPECT().
					GetGroupMemberByGroupIDAndUserID(gomock.Any(), gomock.Eq(arg1)).
					Times(0)
				store.EXPECT().
					CreateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				// "group_id":    expense.GroupID,
				"paid_by_id":  group_member.UserID,
				"amount":      expense.Amount,
				"description": expense.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateExpenseParams{
					// GroupID:     expense.GroupID,
					PaidByID:    group_member.UserID,
					Amount:      expense.Amount,
					Description: expense.Description,
				}

				arg1 := db.GetGroupMemberByGroupIDAndUserIDParams{
					// GroupID: expense.GroupID,
					UserID: user.ID,
				}

				store.EXPECT().
					GetGroupMemberByGroupIDAndUserID(gomock.Any(), gomock.Eq(arg1)).
					Times(0)
				store.EXPECT().
					CreateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},

		{
			name: "InternalError",
			body: gin.H{
				"group_id":    group_member.GroupID,
				"paid_by_id":  group_member.UserID,
				"amount":      expense.Amount,
				"description": expense.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateExpenseParams{
					GroupID:     group_member.GroupID,
					PaidByID:    group_member.UserID,
					Amount:      expense.Amount,
					Description: expense.Description,
				}

				arg1 := db.GetGroupMemberByGroupIDAndUserIDParams{
					GroupID: group_member.GroupID,
					UserID:  group_member.UserID,
				}

				store.EXPECT().
					GetGroupMemberByGroupIDAndUserID(gomock.Any(), gomock.Eq(arg1)).
					Times(1).
					Return(group_member, nil)
				store.EXPECT().
					CreateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Expenses{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},

		{
			name: "InternalErrorInGroupMemberByGroupIDAndUserID",
			body: gin.H{
				"group_id":    group_member.GroupID,
				"paid_by_id":  group_member.UserID,
				"amount":      expense.Amount,
				"description": expense.Description,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateExpenseParams{
					GroupID:     group_member.GroupID,
					PaidByID:    group_member.UserID,
					Amount:      expense.Amount,
					Description: expense.Description,
				}

				arg1 := db.GetGroupMemberByGroupIDAndUserIDParams{
					GroupID: group_member.GroupID,
					UserID:  group_member.UserID,
				}

				store.EXPECT().
					GetGroupMemberByGroupIDAndUserID(gomock.Any(), gomock.Eq(arg1)).
					Times(1).
					Return(db.GroupMembers{}, sql.ErrConnDone)
				store.EXPECT().
					CreateExpense(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/expenses"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func randomExpense(group_member db.GroupMembers) db.Expenses {
	return db.Expenses{
		ID:          int64(util.RandomInt(1, 1000)),
		GroupID:     group_member.ID,
		PaidByID:    group_member.ID,
		Amount:      util.RandomMoney(),
		Description: util.RandomString(10),
		CreatedAt:   util.RandomDatetime(),
	}
}

func requireBodyMatchExpense(t *testing.T, body *bytes.Buffer, expense db.Expenses) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var got db.Expenses
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	require.Equal(t, expense, got)
	require.Equal(t, expense.ID, got.ID)
	require.Equal(t, expense.GroupID, got.GroupID)
	require.Equal(t, expense.PaidByID, got.PaidByID)
	require.Equal(t, expense.Amount, got.Amount)

}
