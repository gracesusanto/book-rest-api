package database

import "time"

type Book struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Author      string    `gorm:"not null" json:"author"`
	PublishedAt time.Time `gorm:"not null" json:"published_at"`
	Edition     string    `json:"edition"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Collection struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Books     []Book    `gorm:"many2many:collection_books;" json:"books"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CollectionBook struct {
	CollectionID uint `gorm:"primaryKey" json:"collection_id"`
	BookID       uint `gorm:"primaryKey" json:"book_id"`
}
