package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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

func TestCreateGroupMemberAPI(t *testing.T) {
	user, _ := randomUser(t)
	user2, _ := randomUser(t)
	category := randomGroupCategory()
	group := randomGroup(t, user, category)

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
				"group_id":   group.ID,
				"first_name": user2.FirstName,
				"last_name":  user2.LastName,
				"phone":      user2.Phone,
				"email":      user2.Email,
				"time_zone":  user2.TimeZone,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddUserToGroupParams{
					GroupID:      group.ID,
					FirstName:    user2.FirstName,
					LastName:     user2.LastName,
					Phone:        user2.Phone,
					Email:        user2.Email,
					TimeZone:     user2.TimeZone,
					GroupOwnerID: user.ID,
				}
				store.EXPECT().
					AddUserToGroupTx(gomock.Any(), arg).
					Times(1).
					Return(db.AddUserToGroupResults{
						GroupMembers: db.GroupMembers{},
						User:         user2,
						Invitations:  db.Invitations{},
					}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"group_id":   group.ID,
				"first_name": user2.FirstName,
				"last_name":  user2.LastName,
				"phone":      user2.Phone,
				"email":      user2.Email,
				"time_zone":  user2.TimeZone,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddUserToGroupParams{
					GroupID:      group.ID,
					FirstName:    user2.FirstName,
					LastName:     user2.LastName,
					Phone:        user2.Phone,
					Email:        user2.Email,
					TimeZone:     user2.TimeZone,
					GroupOwnerID: user.ID,
				}
				store.EXPECT().
					AddUserToGroupTx(gomock.Any(), arg).
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
				// "group_id":   group.ID,
				"first_name": user2.FirstName,
				"last_name":  user2.LastName,
				"phone":      user2.Phone,
				"email":      user2.Email,
				"time_zone":  user2.TimeZone,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddUserToGroupParams{
					// GroupID:      group.ID,
					FirstName:    user2.FirstName,
					LastName:     user2.LastName,
					Phone:        user2.Phone,
					Email:        user2.Email,
					TimeZone:     user2.TimeZone,
					GroupOwnerID: user.ID,
				}
				store.EXPECT().
					AddUserToGroupTx(gomock.Any(), arg).
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
				"group_id":   group.ID,
				"first_name": user2.FirstName,
				"last_name":  user2.LastName,
				"phone":      user2.Phone,
				"email":      user2.Email,
				"time_zone":  user2.TimeZone,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddUserToGroupParams{
					GroupID:      group.ID,
					FirstName:    user2.FirstName,
					LastName:     user2.LastName,
					Phone:        user2.Phone,
					Email:        user2.Email,
					TimeZone:     user2.TimeZone,
					GroupOwnerID: user.ID,
				}
				store.EXPECT().
					AddUserToGroupTx(gomock.Any(), arg).
					Times(1).
					Return(db.AddUserToGroupResults{}, sql.ErrConnDone)
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

			url := "/group_members"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}

func TestListGroupMembersAPI(t *testing.T) {
	user, _ := randomUser(t)
	category := randomGroupCategory()
	group := randomGroup(t, user, category)

	n := 5

	members := []db.ListGroupMembersWithDetailsRow{}

	type Query struct {
		Limit   int32
		Offset  int32
		GroupID int64
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGroupMembersWithDetailsParams{
					GroupID: group.ID,
					Limit:   5,
					Offset:  0,
				}
				store.EXPECT().
					ListGroupMembersWithDetails(gomock.Any(), arg).
					Times(1).
					Return(members, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			query: Query{
				Limit:   int32(n),
				Offset:  1,
				GroupID: group.ID,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGroupMembersWithDetailsParams{
					GroupID: group.ID,
					Limit:   5,
					Offset:  0,
				}
				store.EXPECT().
					ListGroupMembersWithDetails(gomock.Any(), arg).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)
			},
			query: Query{
				Limit:   int32(n),
				Offset:  1,
				GroupID: group.ID,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGroupMembersWithDetailsParams{
					GroupID: group.ID,
					Limit:   5,
					Offset:  0,
				}
				store.EXPECT().
					ListGroupMembersWithDetails(gomock.Any(), arg).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			query: Query{
				Limit:   0,
				Offset:  1,
				GroupID: group.ID,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGroupMembersWithDetailsParams{
					GroupID: group.ID,
					Limit:   5,
					Offset:  0,
				}
				store.EXPECT().
					ListGroupMembersWithDetails(gomock.Any(), arg).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			query: Query{
				Limit:   int32(n),
				Offset:  1,
				GroupID: group.ID,
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

			url := "/group_members"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("offset", fmt.Sprintf("%d", tc.query.Offset))
			q.Add("limit", fmt.Sprintf("%d", tc.query.Limit))
			q.Add("group_id", fmt.Sprintf("%d", tc.query.GroupID))
			request.URL.RawQuery = q.Encode()

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func TestDeleteGroupMemberAPI(t *testing.T) {

	user, _ := randomUser(t)

	groupCategory := randomGroupCategory()

	group := randomGroup(t, user, groupCategory)

	member := randomGroupMember(t, user, group)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		groupID       int64
		memberID      int64
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name:     "OK",
			groupID:  group.ID,
			memberID: member.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetGroupMemberByID(gomock.Any(), gomock.Eq(member.ID)).
					Times(1).
					Return(member, nil)

				store.EXPECT().
					DeleteGroupMember(gomock.Any(), gomock.Eq(member.ID)).
					Times(1).
					Return(nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

			},
		},
		{
			name:     "Unauthorized",
			groupID:  group.ID,
			memberID: member.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetGroupMemberByID(gomock.Any(), gomock.Eq(member.ID)).
					Times(0)

				store.EXPECT().
					DeleteGroupMember(gomock.Any(), gomock.Eq(member.ID)).
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
			name:    "BadRequest",
			groupID: group.ID,
			// memberID: member.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetGroupMemberByID(gomock.Any(), gomock.Eq(member.ID)).
					Times(0)

				store.EXPECT().
					DeleteGroupMember(gomock.Any(), gomock.Eq(member.ID)).
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
			name:     "InternalError",
			groupID:  group.ID,
			memberID: member.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetGroupMemberByID(gomock.Any(), gomock.Eq(member.ID)).
					Times(1).
					Return(db.GroupMembers{}, sql.ErrConnDone)

				store.EXPECT().
					DeleteGroupMember(gomock.Any(), gomock.Eq(member.ID)).
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

			url := fmt.Sprintf("/group_members/%d/%d", tc.groupID, tc.memberID)

			request, err := http.NewRequest(http.MethodDelete, url, nil)

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func randomGroupMember(t *testing.T, user db.Users, group db.Groups) db.GroupMembers {

	member := db.GroupMembers{
		ID:        int64(util.RandomInt(1, 1000)),
		GroupID:   group.ID,
		UserID:    user.ID,
		CreatedAt: util.RandomDatetime(),
	}

	return member
}
