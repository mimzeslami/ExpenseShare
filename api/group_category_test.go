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

func TestCreateGroupCategoryAPI(t *testing.T) {

	user, _ := randomUser(t)

	groupCategory := randomGroupCategory()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{"name": groupCategory.Name},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateGroupCategory(gomock.Any(), gomock.Eq(groupCategory.Name)).
					Times(1).
					Return(groupCategory, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGroupCategory(t, recorder.Body, groupCategory)
			},
		},
		{
			name: "Unauthrized",
			body: gin.H{"name": groupCategory.Name},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateGroupCategory(gomock.Any(), gomock.Eq(groupCategory.Name)).
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
			body: gin.H{"name": ""},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateGroupCategory(gomock.Any(), gomock.Any()).
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
			body: gin.H{"name": groupCategory.Name},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateGroupCategory(gomock.Any(), gomock.Eq(groupCategory.Name)).
					Times(1).
					Return(db.GroupCategories{}, sql.ErrConnDone)
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

			url := "/group_categories"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func TestListGroupCategoriesAPI(t *testing.T) {

	user, _ := randomUser(t)

	n := 5
	groupCategories := make([]db.GroupCategories, n)

	for i := 0; i < n; i++ {
		groupCategories[i] = randomGroupCategory()
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

				arg := db.ListGroupCategoriesParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListGroupCategories(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(groupCategories, nil)
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
				requireBodyMatchGroupCategories(t, recorder.Body, groupCategories)
			},
		},
		{
			name: "Unauthrized",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListGroupCategories(gomock.Any(), gomock.Any()).
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

				store.EXPECT().
					ListGroupCategories(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			query: Query{
				Limit:  0,
				Offset: 0,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListGroupCategories(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.GroupCategories{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			query: Query{
				Limit:  int32(n),
				Offset: 1,
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

			url := "/group_categories"
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

func TestGetGroupCategoryAPI(t *testing.T) {

	user, _ := randomUser(t)

	groupCategory := randomGroupCategory()

	testCases := []struct {
		name            string
		buildStubs      func(store *mockdb.MockStore)
		setupAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		groupCategoryID int64
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name:            "OK",
			groupCategoryID: groupCategory.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetGroupCategory(gomock.Any(), gomock.Eq(groupCategory.ID)).
					Times(1).
					Return(groupCategory, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGroupCategory(t, recorder.Body, groupCategory)
			},
		},
		{
			name: "Unauthrized",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetGroupCategory(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)

			},
			groupCategoryID: groupCategory.ID,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
		{
			name: "BadRequest",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetGroupCategory(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)

			},
			groupCategoryID: 0,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetGroupCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GroupCategories{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			groupCategoryID: groupCategory.ID,
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

			url := fmt.Sprintf("/group_categories/%d", tc.groupCategoryID)

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}

func TestUpdateGroupCategoryAPI(t *testing.T) {

	user, _ := randomUser(t)

	groupCategory := randomGroupCategory()

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		body          gin.H
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name: "OK",
			body: gin.H{"id": groupCategory.ID, "name": groupCategory.Name},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateGroupCategoryParams{
					ID:   groupCategory.ID,
					Name: groupCategory.Name,
				}

				store.EXPECT().
					UpdateGroupCategory(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(groupCategory, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGroupCategory(t, recorder.Body, groupCategory)
			},
		},
		{
			name: "Unauthrized",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateGroupCategoryParams{
					ID:   groupCategory.ID,
					Name: groupCategory.Name,
				}

				store.EXPECT().
					UpdateGroupCategory(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			body: gin.H{"id": groupCategory.ID, "name": groupCategory.Name},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
		{
			name: "BadRequest",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateGroupCategory(gomock.Any(), gomock.Any()).
					Times(0)
			},
			body: gin.H{"id": "", "name": groupCategory.Name},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateGroupCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GroupCategories{}, sql.ErrConnDone)
			},
			body: gin.H{"id": groupCategory.ID, "name": groupCategory.Name},
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

			url := "/group_categories"

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}

func TestDeleteGroupCategoryAPI(t *testing.T) {

	user, _ := randomUser(t)

	groupCategory := randomGroupCategory()

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		groupCategory int64
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name:          "OK",
			groupCategory: groupCategory.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					DeleteGroupCategory(gomock.Any(), gomock.Eq(groupCategory.ID)).
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
			name:          "Unauthrized",
			groupCategory: groupCategory.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					DeleteGroupCategory(gomock.Any(), gomock.Any()).
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
			name:          "BadRequest",
			groupCategory: 0,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					DeleteGroupCategory(gomock.Any(), gomock.Any()).
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
			name:          "InternalError",
			groupCategory: groupCategory.ID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					DeleteGroupCategory(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/group_categories/%d", tc.groupCategory)

			request, err := http.NewRequest(http.MethodDelete, url, nil)

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}

}

func randomGroupCategory() db.GroupCategories {
	return db.GroupCategories{
		ID:        int64(util.RandomInt(1, 1000)),
		Name:      util.RandomString(6),
		CreatedAt: util.RandomDatetime(),
	}
}

func requireBodyMatchGroupCategory(t *testing.T, body *bytes.Buffer, groupCategory db.GroupCategories) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGroupCategory db.GroupCategories
	err = json.Unmarshal(data, &gotGroupCategory)
	require.NoError(t, err)
	require.Equal(t, groupCategory.Name, gotGroupCategory.Name)
}

func requireBodyMatchGroupCategories(t *testing.T, body *bytes.Buffer, groupCategories []db.GroupCategories) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGroupCategories []db.GroupCategories
	err = json.Unmarshal(data, &gotGroupCategories)
	require.NoError(t, err)
	require.Equal(t, groupCategories, gotGroupCategories)
}
