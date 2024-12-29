package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAuthor(t *testing.T) Author {
	arg := randomString(6)
	author, err := testQueries.CreateAuthor(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, author)
	return author
}

func TestCreateAuthor(t *testing.T) {
	createRandomAuthor(t)
}

func TestGetAuthorByName(t *testing.T) {
	author1 := createRandomAuthor(t)
	author2, err := testQueries.GetAuthorByName(context.Background(), author1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, author2)
	require.Equal(t, author1.Name, author2.Name)
	require.WithinDuration(t, author1.CreatedAt, author2.CreatedAt, time.Second)
}

func TestDeleteAuthor(t *testing.T) {
	author1 := createRandomAuthor(t)
	err := testQueries.DeleteAuthor(context.Background(), author1.Name)
	require.NoError(t, err)

	author2, err := testQueries.GetAuthorByName(context.Background(), author1.Name)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, author2)
}

func TestGetAllAuthors(t *testing.T) {
	for i := 0; i < 4; i++ {
		createRandomAuthor(t)
	}

	arg := GetAllAuthorsParams{
		Limit:  2,
		Offset: 0,
	}

	authors, err := testQueries.GetAllAuthors(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, authors, 2)

	for _, author := range authors {
		require.NotEmpty(t, author)
	}
}