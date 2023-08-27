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

func TestCreateGroupAPI(t *testing.T) {
	user, _ := randomUser(t)
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
				"name":          group.Name,
				"category_id":   group.CategoryID,
				"created_by_id": group.CreatedByID,
				"image_path":    group.ImagePath,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateGroupParams{
					Name:        group.Name,
					CategoryID:  group.CategoryID,
					CreatedByID: group.CreatedByID,
					ImagePath:   group.ImagePath,
				}
				store.EXPECT().
					CreateGroupTx(gomock.Any(), arg).
					Times(1).
					Return(db.CreateGroupTxResults{
						Group:        group,
						GroupMembers: db.GroupMembers{},
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
				"name":          group.Name,
				"category_id":   group.CategoryID,
				"created_by_id": group.CreatedByID,
				"image_path":    group.ImagePath,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateGroupParams{
					Name:        group.Name,
					CategoryID:  group.CategoryID,
					CreatedByID: group.CreatedByID,
					ImagePath:   group.ImagePath,
				}
				store.EXPECT().
					CreateGroupTx(gomock.Any(), arg).
					Times(0)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"name": group.Name,
				// "category_id":   group.CategoryID,
				"created_by_id": group.CreatedByID,
				"image_path":    group.ImagePath,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateGroupParams{
					Name: group.Name,
					// CategoryID:  group.CategoryID,
					CreatedByID: group.CreatedByID,
					ImagePath:   group.ImagePath,
				}
				store.EXPECT().
					CreateGroupTx(gomock.Any(), arg).
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
				"name":          group.Name,
				"category_id":   group.CategoryID,
				"created_by_id": group.CreatedByID,
				"image_path":    group.ImagePath,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateGroupParams{
					Name:        group.Name,
					CategoryID:  group.CategoryID,
					CreatedByID: group.CreatedByID,
					ImagePath:   group.ImagePath,
				}
				store.EXPECT().
					CreateGroupTx(gomock.Any(), arg).
					Times(1).
					Return(db.CreateGroupTxResults{}, sql.ErrConnDone)
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

			url := "/groups"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func TestListGroupsAPI(t *testing.T) {

	user, _ := randomUser(t)
	category := randomGroupCategory()

	n := 5
	groups := make([]db.Groups, n)

	for i := 0; i < n; i++ {
		groups[i] = randomGroup(t, user, category)
	}

	type Query struct {
		Limit  int32
		Offset int32
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		query         Query
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGroupsParams{
					Limit:       int32(n),
					Offset:      0,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					ListGroups(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(groups, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			query: Query{
				Limit:  int32(n),
				Offset: 1,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGroups(t, recorder.Body, groups)
			},
		},
		{
			name: "Unauthorized",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGroupsParams{
					Limit:       int32(n),
					Offset:      0,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					ListGroups(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)

			},
			query: Query{
				Limit:  int32(n),
				Offset: 1,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
		{
			name: "BadRequest",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGroupsParams{
					Limit:       int32(n),
					Offset:      0,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					ListGroups(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			query: Query{
				Limit:  0,
				Offset: 1,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGroupsParams{
					Limit:       int32(n),
					Offset:      0,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					ListGroups(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			query: Query{
				Limit:  int32(n),
				Offset: 1,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

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

			url := "/groups"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("offset", fmt.Sprintf("%d", tc.query.Offset))
			q.Add("limit", fmt.Sprintf("%d", tc.query.Limit))
			request.URL.RawQuery = q.Encode()

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}

func TestGetGroupAPI(t *testing.T) {

	user, _ := randomUser(t)
	category := randomGroupCategory()
	group := randomGroup(t, user, category)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		groupID       int64
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetGroupByIDParams{
					ID:          group.ID,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					GetGroupByID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(group, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, group.CreatedByID, time.Minute)

			},
			groupID: group.ID,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGroup(t, recorder.Body, group)
			},
		},
		{
			name: "Unauthorized",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetGroupByIDParams{
					ID:          group.ID,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					GetGroupByID(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			groupID: group.ID,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, group.CreatedByID, -time.Minute)

			},
		},
		{
			name: "BadRequest",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetGroupByIDParams{
					ID:          group.ID,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					GetGroupByID(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			groupID: 0,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, group.CreatedByID, time.Minute)
			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetGroupByIDParams{
					ID:          group.ID,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					GetGroupByID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Groups{}, sql.ErrConnDone)
			},
			groupID: group.ID,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, group.CreatedByID, time.Minute)
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

			url := fmt.Sprintf("/groups/%d", tc.groupID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func TestUpdateGroupAPI(t *testing.T) {

	user, _ := randomUser(t)

	groupCategory := randomGroupCategory()

	group := randomGroup(t, user, groupCategory)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		body          gin.H
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name: "OK",
			body: gin.H{"id": group.ID, "name": group.Name, "category_id": group.CategoryID, "image_path": group.ImagePath},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateGroupParams{
					ID:          group.ID,
					Name:        group.Name,
					CategoryID:  group.CategoryID,
					ImagePath:   group.ImagePath,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					UpdateGroup(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(group, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGroup(t, recorder.Body, group)

			},
		},
		{
			name: "Unauthorized",
			body: gin.H{"id": group.ID, "name": group.Name, "category_id": group.CategoryID, "image_path": group.ImagePath},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateGroupParams{
					ID:          group.ID,
					Name:        group.Name,
					CategoryID:  group.CategoryID,
					ImagePath:   group.ImagePath,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					UpdateGroup(gomock.Any(), gomock.Eq(arg)).
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
			body: gin.H{"name": group.Name, "category_id": group.CategoryID, "image_path": group.ImagePath},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateGroupParams{
					// ID:          group.ID,
					Name:        group.Name,
					CategoryID:  group.CategoryID,
					ImagePath:   group.ImagePath,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					UpdateGroup(gomock.Any(), gomock.Eq(arg)).
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
			body: gin.H{"id": group.ID, "name": group.Name, "category_id": group.CategoryID, "image_path": group.ImagePath},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateGroupParams{
					ID:          group.ID,
					Name:        group.Name,
					CategoryID:  group.CategoryID,
					ImagePath:   group.ImagePath,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					UpdateGroup(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Groups{}, sql.ErrConnDone)
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

			url := "/groups"

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}

func TestDeleteGroupAPI(t *testing.T) {

	user, _ := randomUser(t)

	groupCategory := randomGroupCategory()

	group := randomGroup(t, user, groupCategory)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		groupID       int64
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name:    "OK",
			groupID: group.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.DeleteGroupParams{
					ID:          group.ID,
					CreatedByID: user.ID,
				}

				store.EXPECT().
					DeleteGroupTx(gomock.Any(), gomock.Eq(arg)).
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
			name:    "Unauthrized",
			groupID: group.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					DeleteGroupTx(gomock.Any(), gomock.Any()).
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
			groupID: 0,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					DeleteGroupTx(gomock.Any(), gomock.Any()).
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
			name:    "InternalError",
			groupID: group.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					DeleteGroupTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
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

			url := fmt.Sprintf("/groups/%d", tc.groupID)

			request, err := http.NewRequest(http.MethodDelete, url, nil)

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func randomGroup(t *testing.T, user db.Users, category db.GroupCategories) (group db.Groups) {

	g := db.Groups{
		ID:          int64(util.RandomInt(1, 1000)),
		Name:        util.RandomString(6),
		CategoryID:  category.ID,
		ImagePath:   util.RandomString(6),
		CreatedByID: user.ID,
		CreatedAt:   util.RandomDatetime(),
	}

	return g

}

func requireBodyMatchGroups(t *testing.T, body *bytes.Buffer, groups []db.Groups) {

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGroups []db.Groups
	err = json.Unmarshal(data, &gotGroups)
	require.NoError(t, err)

	require.Equal(t, groups, gotGroups)

}

func requireBodyMatchGroup(t *testing.T, body *bytes.Buffer, group db.Groups) {

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGroup db.Groups
	err = json.Unmarshal(data, &gotGroup)
	require.NoError(t, err)

	require.Equal(t, group, gotGroup)

}
