package db

import (
	"context"
	"testing"
	"time"

	"github.com/drmanalo/simplebank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	assert.NoError(t, err)

	arg := CreateUserParams{
		Email:          util.RandomEmail(),
		FullName:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Username:       util.RandomOwner(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, arg.Email, user.Email)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.Username, user.Username)

	assert.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	assert.NoError(t, err)
	assert.NotEmpty(t, user2)

	assert.Equal(t, user1.Email, user2.Email)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
	assert.Equal(t, user1.Username, user2.Username)
	assert.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}
