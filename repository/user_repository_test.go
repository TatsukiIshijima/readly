package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"readly/testdata"
	"testing"
)

func createRandomUser(t *testing.T) *CreateUserResponse {
	name := testdata.RandomString(10)
	email := testdata.RandomEmail()
	password := testdata.RandomString(16)

	req := CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}
	u, err := userRepo.CreateUser(context.Background(), req)
	require.NoError(t, err)
	return u
}

func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)
	gu, err := querier.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)

	require.Equal(t, user.ID, gu.ID)
	require.Equal(t, user.Name, gu.Name)
	require.Equal(t, user.Email, gu.Email)
}
