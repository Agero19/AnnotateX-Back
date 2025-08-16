package repository

import (
	"database/sql"

	_ "github.com/lib/pq" // Import the pq driver for PostgreSQL
)

// Repository is a struct that holds the database connection and repositories for different entities.
type Repository struct {
	Users       Users
	Images      Images
	Annotations Annotations
}

type Users interface {
	Create(user *User) error
	GetAll() ([]*User, error)
	GetByID(id string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}

type Images interface {
	Create(image *Image) error
	GetAll() ([]*Image, error)
	GetByID(id string) (*Image, error)
	Update(image *Image) error
	Delete(id string) error
}

type Annotations interface {
	Create(annotation *Annotation) error
	GetAll() ([]*Annotation, error)
	GetByID(id string) (*Annotation, error)
	Update(annotation *Annotation) error
	Delete(id string) error
}

// New creates a new Repository instance from a PostgreSQL connection pool.
func NewRepository(db *sql.DB) Repository {
	return Repository{
		Users:       &UserRepository{db: db},
		Images:      &ImageRepository{db: db},
		Annotations: &AnnotationRepository{db: db},
	}
}
