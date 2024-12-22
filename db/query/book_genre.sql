-- name: CreateBookGenre :exec
INSERT INTO book_genres (book_id, genre_name)
VALUES ($1, $2);

-- name: GetGenresByBookID :many
SELECT genre_name
FROM book_genres
WHERE book_id = $1
ORDER BY book_id LIMIT $2
OFFSET $3;

-- name: GetBooksByGenreName :many
SELECT book_id
FROM book_genres
WHERE genre_name = $1
ORDER BY book_id LIMIT $2
OFFSET $3;

-- name: UpdateGenreForBook :exec
UPDATE book_genres
SET genre_name = $3
WHERE book_id = $1
  AND genre_name = $2;

-- name: DeleteGenreForBook :exec
DELETE
FROM book_genres
WHERE book_id = $1
  AND genre_name = $2;