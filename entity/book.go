package entity

import "time"

type Book struct {
	ID            int64         `json:"id"`
	Title         string        `json:"title"`
	Genres        []string      `json:"genres"`
	Description   *string       `json:"description"`
	CoverImageURL *string       `json:"cover_image_url"`
	URL           *string       `json:"url"`
	AuthorName    *string       `json:"author_name"`
	PublisherName *string       `json:"publisher_name"`
	PublishDate   *time.Time    `json:"publish_date"`
	ISBN          *string       `json:"isbn"`
	Status        ReadingStatus `json:"status"`
	StartDate     *time.Time    `json:"start_date"`
	EndDate       *time.Time    `json:"end_date"`
}
