package repository

import (
	"context"
	"database/sql"
	"errors"
	db "readly/db/sqlc"
	"readly/domain"
	"time"
)

type BookRepository interface {
	Register(ctx context.Context, req RegisterRequest) (domain.Book, error)
	Get(ctx context.Context, id int64) (*domain.Book, error)
	List(ctx context.Context, req ListRequest) ([]*domain.Book, error)
	Delete(ctx context.Context, req DeleteRequest) error
}

type BookRepositoryImpl struct {
	store *Store
}

func NewBookRepository(store *Store) BookRepository {
	return BookRepositoryImpl{store: store}
}

type RegisterRequest struct {
	UserID        int64
	Title         string
	Genres        []string
	Description   string
	CoverImageURL string
	URL           string
	AuthorName    string
	PublisherName string
	PublishDate   time.Time
	ISBN          string
}

func (r BookRepositoryImpl) Register(ctx context.Context, args RegisterRequest) (domain.Book, error) {
	var result domain.Book

	err := r.store.execTx(ctx, func(q *db.Queries) error {
		if err := r.registerAuthorIfNotExist(ctx, q, args.AuthorName); err != nil {
			return err
		}
		if err := r.registerPublisherIfNotExist(ctx, q, args.PublisherName); err != nil {
			return err
		}
		for _, genre := range args.Genres {
			if err := r.registerGenreIfNotExist(ctx, q, genre); err != nil {
				return err
			}
		}
		book, err := q.CreateBook(ctx, db.CreateBookParams{
			Title:         sql.NullString{String: args.Title, Valid: true},
			Description:   sql.NullString{String: args.Description, Valid: true},
			CoverImageUrl: sql.NullString{String: args.CoverImageURL, Valid: true},
			Url:           sql.NullString{String: args.URL, Valid: true},
			AuthorName:    args.AuthorName,
			PublisherName: args.PublisherName,
			PublishedDate: sql.NullTime{Time: args.PublishDate, Valid: true},
			Isbn:          sql.NullString{String: args.ISBN, Valid: true},
		})
		if err != nil {
			return err
		}
		for _, genre := range args.Genres {
			if _, err := q.CreateBookGenre(ctx, db.CreateBookGenreParams{
				BookID:    book.ID,
				GenreName: genre,
			}); err != nil {
				return err
			}
		}
		if _, err := q.CreateReadingHistory(ctx, db.CreateReadingHistoryParams{
			UserID:    args.UserID,
			BookID:    book.ID,
			Status:    db.ReadingStatusUnread,
			StartDate: sql.NullTime{Time: time.Time{}, Valid: true},
			EndDate:   sql.NullTime{Time: time.Time{}, Valid: false},
		}); err != nil {
			return err
		}
		genres, err := q.GetGenresByBookID(ctx, book.ID)
		if err != nil {
			return err
		}
		result = domain.Book{
			ID:            book.ID,
			Title:         book.Title.String,
			Genres:        genres,
			Description:   book.Description.String,
			CoverImageURL: book.CoverImageUrl.String,
			URL:           book.Url.String,
			AuthorName:    book.AuthorName,
			PublisherName: book.PublisherName,
			PublishDate:   book.PublishedDate.Time,
			ISBN:          book.Isbn.String,
		}

		return nil
	})
	return result, err
}

func (r BookRepositoryImpl) registerAuthorIfNotExist(ctx context.Context, q *db.Queries, name string) error {
	var err error
	_, err = q.GetAuthorByName(ctx, name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		_, err = q.CreateAuthor(ctx, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r BookRepositoryImpl) registerPublisherIfNotExist(ctx context.Context, q *db.Queries, name string) error {
	_, err := q.GetPublisherByName(ctx, name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		_, err = q.CreatePublisher(ctx, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r BookRepositoryImpl) registerGenreIfNotExist(ctx context.Context, q *db.Queries, name string) error {
	_, err := q.GetGenreByName(ctx, name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		_, err = q.CreateGenre(ctx, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r BookRepositoryImpl) Get(ctx context.Context, id int64) (*domain.Book, error) {
	book, err := r.store.Queries.GetBookById(ctx, id)
	if err != nil {
		return nil, err
	}
	genres, err := r.store.Queries.GetGenresByBookID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.Book{
		ID:            book.ID,
		Title:         book.Title.String,
		Genres:        genres,
		Description:   book.Description.String,
		CoverImageURL: book.CoverImageUrl.String,
		URL:           book.Url.String,
		AuthorName:    book.AuthorName,
		PublisherName: book.PublisherName,
		PublishDate:   book.PublishedDate.Time,
		ISBN:          book.Isbn.String,
	}, nil
}

type ListRequest struct {
	UserID int64
	Limit  int32
	Offset int32
}

func (r BookRepositoryImpl) List(ctx context.Context, req ListRequest) ([]*domain.Book, error) {
	historyParams := db.GetReadingHistoryByUserIDParams{
		UserID: req.UserID,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	histories, err := r.store.Queries.GetReadingHistoryByUserID(ctx, historyParams)
	if err != nil {
		return nil, err
	}
	res := make([]*domain.Book, 0, len(histories))
	for _, history := range histories {
		book, err := r.Get(ctx, history.BookID)
		if err != nil {
			return nil, err
		}
		res = append(res, book)
	}
	return res, nil
}

type DeleteRequest struct {
	UserID int64
	BookID int64
}

func (r BookRepositoryImpl) Delete(ctx context.Context, req DeleteRequest) error {
	err := r.store.execTx(ctx, func(q *db.Queries) error {
		deleteHistoryParam := db.DeleteReadingHistoryParams{
			UserID: req.UserID,
			BookID: req.BookID,
		}
		if err := r.store.Queries.DeleteReadingHistory(ctx, deleteHistoryParam); err != nil {
			return err
		}
		if err := r.deleteBookGenres(ctx, req.BookID); err != nil {
			return err
		}
		if err := r.store.Queries.DeleteBook(ctx, req.BookID); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r BookRepositoryImpl) deleteBookGenres(ctx context.Context, bookID int64) error {
	bookGenres, err := r.store.Queries.GetGenresByBookID(ctx, bookID)
	if err != nil {
		return err
	}
	for _, genre := range bookGenres {
		param := db.DeleteGenreForBookParams{
			BookID:    bookID,
			GenreName: genre,
		}
		err := r.store.Queries.DeleteGenreForBook(ctx, param)
		if err != nil {
			return err
		}
	}
	return nil
}
