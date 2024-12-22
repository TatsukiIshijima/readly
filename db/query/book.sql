-- name: ListBooksByTitle :many
SELECT *
FROM books
WHERE title LIKE $1;

-- name: ListBooksByIsbn :many
SELECT *
FROM books
WHERE isbn = $1;

-- name: ListBooksByAuthorName :many
SELECT *
FROM books
WHERE author_name LIKE $1;

-- name: InsertBook :one
INSERT INTO books (title,
                   description,
                   cover_image_url,
                   url,
                   author_name,
                   publisher_name,
                   published_date,
                   isbn)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateBook :one
UPDATE books
SET title            = $2,
    description      = $3,
    cover_image_url  = $4,
    url              = $5,
    author_name      = $6,
    publisher_name   = $7,
    published_date   = $8,
    isbn             = $9,
    updated_at       = $10
WHERE id = $1 RETURNING *;