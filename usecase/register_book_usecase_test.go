package usecase

import (
	"context"
	"github.com/stretchr/testify/require"
	"readly/entity"
	"readly/testdata"
	"testing"
	"time"
)

func TestRegisterBook(t *testing.T) {
	signUpUseCase := newTestSignUpUseCase(t)
	registerBookUseCase := newTestRegisterBookUseCase(t)

	signUpReq := SignUpRequest{
		Name:     testdata.RandomString(10),
		Email:    testdata.RandomEmail(),
		Password: testdata.RandomString(16),
	}
	signUpRes, err := signUpUseCase.SignUp(context.Background(), signUpReq)
	require.NoError(t, err)

	testCases := []struct {
		name  string
		setup func(t *testing.T) RegisterBookRequest
		check func(t *testing.T, req RegisterBookRequest, res *entity.Book, err error)
	}{
		{
			name: "New unread book with required fields register success",
			setup: func(t *testing.T) RegisterBookRequest {
				return RegisterBookRequest{
					UserID: signUpRes.UserID,
					Title:  testdata.RandomString(10),
					Genres: []string{testdata.RandomString(6)},
					Status: 0,
				}
			},
			check: func(t *testing.T, req RegisterBookRequest, res *entity.Book, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.ID)
				require.Equal(t, req.Title, res.Title)
				require.Equal(t, req.Genres, res.Genres)
				require.Equal(t, req.Description, res.Description)
				require.Equal(t, req.CoverImageURL, res.CoverImageURL)
				require.Equal(t, req.URL, res.URL)
				require.Equal(t, req.AuthorName, res.AuthorName)
				require.Equal(t, req.PublisherName, res.PublisherName)
				require.Equal(t, req.PublishDate, res.PublishDate)
				require.Equal(t, req.ISBN, res.ISBN)
				require.Equal(t, req.Status, res.Status)
				require.True(t, isSameDate(req.StartDate, res.StartDate))
				require.True(t, isSameDate(req.EndDate, res.EndDate))
			},
		},
		{
			name: "New read book with all fields register success",
			setup: func(t *testing.T) RegisterBookRequest {
				desc := testdata.RandomString(100)
				coverImgURL := testdata.RandomString(255)
				url := testdata.RandomString(255)
				author := testdata.RandomString(10)
				publisher := testdata.RandomString(10)
				publishDate, err := testdata.TimeFrom("1970-01-01 00:00:00")
				require.NoError(t, err)
				pb := publishDate.UTC()
				ISBN := testdata.RandomString(13)
				startDate := time.Now().UTC()
				endDate := time.Now().Add(time.Duration(60*60*24) * time.Second).UTC()

				return RegisterBookRequest{
					UserID:        signUpRes.UserID,
					Title:         testdata.RandomString(10),
					Genres:        []string{testdata.RandomString(6)},
					Description:   &desc,
					CoverImageURL: &coverImgURL,
					URL:           &url,
					AuthorName:    &author,
					PublisherName: &publisher,
					PublishDate:   &pb,
					ISBN:          &ISBN,
					Status:        2,
					StartDate:     &startDate,
					EndDate:       &endDate,
				}
			},
			check: func(t *testing.T, req RegisterBookRequest, res *entity.Book, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.ID)
				require.Equal(t, req.Title, res.Title)
				require.Equal(t, req.Genres, res.Genres)
				require.Equal(t, req.Description, res.Description)
				require.Equal(t, req.CoverImageURL, res.CoverImageURL)
				require.Equal(t, req.URL, res.URL)
				require.Equal(t, req.AuthorName, res.AuthorName)
				require.Equal(t, req.PublisherName, res.PublisherName)
				require.Equal(t, req.PublishDate, res.PublishDate)
				require.Equal(t, req.ISBN, res.ISBN)
				require.Equal(t, req.Status, res.Status)
				require.True(t, isSameDate(req.StartDate, res.StartDate))
				require.True(t, isSameDate(req.EndDate, res.EndDate))
			},
		},
		{
			name: "New unread book register success when genres are already exist.",
			setup: func(t *testing.T) RegisterBookRequest {
				genre := testdata.RandomString(6)

				req := RegisterBookRequest{
					UserID: signUpRes.UserID,
					Title:  testdata.RandomString(10),
					Genres: []string{genre},
					Status: 0,
				}
				_, err := registerBookUseCase.RegisterBook(context.Background(), req)
				require.NoError(t, err)

				return RegisterBookRequest{
					UserID: signUpRes.UserID,
					Title:  testdata.RandomString(10),
					Genres: []string{genre},
					Status: 0,
				}
			},
			check: func(t *testing.T, req RegisterBookRequest, res *entity.Book, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.ID)
				require.Equal(t, req.Title, res.Title)
				require.Equal(t, req.Genres, res.Genres)
				require.Equal(t, req.Description, res.Description)
				require.Equal(t, req.CoverImageURL, res.CoverImageURL)
				require.Equal(t, req.URL, res.URL)
				require.Equal(t, req.AuthorName, res.AuthorName)
				require.Equal(t, req.PublisherName, res.PublisherName)
				require.Equal(t, req.PublishDate, res.PublishDate)
				require.Equal(t, req.ISBN, res.ISBN)
				require.Equal(t, req.Status, res.Status)
				require.True(t, isSameDate(req.StartDate, res.StartDate))
				require.True(t, isSameDate(req.EndDate, res.EndDate))
			},
		},
		{
			name: "New reading book with author & publisher register success when author & publisher are already exist.",
			setup: func(t *testing.T) RegisterBookRequest {
				author := testdata.RandomString(10)
				publisher := testdata.RandomString(10)
				startDate := time.Now().UTC()

				req := RegisterBookRequest{
					UserID:        signUpRes.UserID,
					Title:         testdata.RandomString(10),
					Genres:        []string{testdata.RandomString(6)},
					AuthorName:    &author,
					PublisherName: &publisher,
					Status:        0,
				}
				_, err := registerBookUseCase.RegisterBook(context.Background(), req)
				require.NoError(t, err)

				return RegisterBookRequest{
					UserID:        signUpRes.UserID,
					Title:         testdata.RandomString(10),
					Genres:        []string{testdata.RandomString(6)},
					AuthorName:    &author,
					PublisherName: &publisher,
					Status:        1,
					StartDate:     &startDate,
				}
			},
			check: func(t *testing.T, req RegisterBookRequest, res *entity.Book, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.ID)
				require.Equal(t, req.Title, res.Title)
				require.Equal(t, req.Genres, res.Genres)
				require.Equal(t, req.Description, res.Description)
				require.Equal(t, req.CoverImageURL, res.CoverImageURL)
				require.Equal(t, req.URL, res.URL)
				require.Equal(t, req.AuthorName, res.AuthorName)
				require.Equal(t, req.PublisherName, res.PublisherName)
				require.Equal(t, req.PublishDate, res.PublishDate)
				require.Equal(t, req.ISBN, res.ISBN)
				require.Equal(t, req.Status, res.Status)
				require.True(t, isSameDate(req.StartDate, res.StartDate))
				require.True(t, isSameDate(req.EndDate, res.EndDate))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := tc.setup(t)
			res, err := registerBookUseCase.RegisterBook(context.Background(), req)
			tc.check(t, req, res, err)
		})
	}
}

func isSameDate(t1, t2 *time.Time) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// TODO: Goroutineを使ったテストケースを追加する
