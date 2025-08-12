package repository

import (
	"database/sql"
)

// Image represents an image in the database - Model
type Image struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Visibility  bool   `json:"visibility"`
	CreatedAt   string `json:"created_at"`
}

// ImageRepository is a struct that provides methods to interact with the image database table. Implements the Images interface.
type ImageRepository struct {
	db *sql.DB
}

// Create inserts a new image into the database. It returns an error if the insertion fails.
func (r *ImageRepository) Create(image *Image) error {
	query := "INSERT INTO images (user_id, url, title, description, visibility) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at"
	err := r.db.QueryRow(
		query,
		image.UserID,
		image.URL,
		image.Title,
		image.Description,
		image.Visibility,
	).Scan(&image.ID, &image.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetAll retrieves all images from the database.
func (r *ImageRepository) GetAll() ([]*Image, error) {
	query := "SELECT id, user_id, url, title, description, visibility, created_at FROM images"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*Image
	for rows.Next() {
		var image Image
		if err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.URL,
			&image.Title,
			&image.Description,
			&image.Visibility,
			&image.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, &image)
	}
	return images, nil
}

// GetByID retrieves an image by its ID from the database. Does not return an error if the image is not found.
func (r *ImageRepository) GetByID(id string) (*Image, error) {
	query := "SELECT id, user_id, url, title, description, visibility, created_at FROM images WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var image Image
	if err := row.Scan(
		&image.ID,
		&image.UserID,
		&image.URL,
		&image.Title,
		&image.Description,
		&image.Visibility,
		&image.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &image, nil
}

// Update modifies an existing image in the database. It returns an error if the update fails.
func (r *ImageRepository) Update(image *Image) error {
	query := "UPDATE images SET url = $1, title = $2, description = $3, visibility = $4 WHERE id = $5"
	_, err := r.db.Exec(
		query,
		image.URL,
		image.Title,
		image.Description,
		image.Visibility,
		image.ID,
	)
	return err
}

// Delete removes an image from the database by its ID. It returns an error if the deletion fails.
func (r *ImageRepository) Delete(id string) error {
	query := "DELETE FROM images WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
