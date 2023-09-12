package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/drmanalo/simplebank/db/mock"
	db "github.com/drmanalo/simplebank/db/sqlc"
	"github.com/drmanalo/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  password,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:    user.Email,
					FullName: user.FullName,
					Username: user.Username,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
				assertBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  password,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:    user.Email,
					FullName: user.FullName,
					Username: user.Username,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(db.User{}, db.ErrUniqueVioloation)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  password,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:    user.Email,
					FullName: user.FullName,
					Username: user.Username,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"email":     "invalid-email",
				"full_name": user.FullName,
				"password":  password,
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:    user.Email,
					FullName: user.FullName,
					Username: user.Username,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  password,
				"username":  "invalid-user#1",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:    user.Email,
					FullName: user.FullName,
					Username: user.Username,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "VeryShortPassword",
			body: gin.H{
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  "123",
				"username":  user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:    user.Email,
					FullName: user.FullName,
					Username: user.Username,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			assert.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func assertBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, gotUser.Email)
	assert.Equal(t, user.FullName, gotUser.FullName)
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Empty(t, gotUser.HashedPassword)
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	assert.NoError(t, err)

	user = db.User{
		Email:          util.RandomEmail(),
		FullName:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Username:       util.RandomOwner(),
	}
	return
}
