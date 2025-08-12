package repository

import "database/sql"

// Annotation represents an annotation in the database - Model
type Annotation struct {
	ID        string `json:"id"`
	ImageID   string `json:"image_id"`
	UserID    string `json:"user_id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

// AnnotationRepository is a struct that provides methods to interact with the annotation database table. Implements the Annotations interface.
type AnnotationRepository struct {
	db *sql.DB
}

// Create inserts a new annotation into the database. It returns an error if the insertion fails.
func (r *AnnotationRepository) Create(annotation *Annotation) error {
	query := `INSERT INTO annotations (image_id, user_id, x, y, width, height, comment) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`
	err := r.db.QueryRow(
		query,
		annotation.ImageID,
		annotation.UserID,
		annotation.X,
		annotation.Y,
		annotation.Width,
		annotation.Height,
		annotation.Comment,
	).Scan(&annotation.ID, &annotation.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetAll retrieves all annotations from the database.
func (r *AnnotationRepository) GetAll() ([]*Annotation, error) {
	query := `SELECT id, image_id, user_id, x, y, width, height, comment, created_at FROM annotations`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var annotations []*Annotation
	for rows.Next() {
		var annotation Annotation
		if err := rows.Scan(
			&annotation.ID,
			&annotation.ImageID,
			&annotation.UserID,
			&annotation.X,
			&annotation.Y,
			&annotation.Width,
			&annotation.Height,
			&annotation.Comment,
			&annotation.CreatedAt); err != nil {
			return nil, err
		}
		annotations = append(annotations, &annotation)
	}
	return annotations, nil
}

// GetByID retrieves an annotation by its ID from the database. Does not return an error if the annotation is not found.
func (r *AnnotationRepository) GetByID(id string) (*Annotation, error) {
	query := `SELECT id, image_id, user_id, x, y, width, height, comment, created_at FROM annotations WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var annotation Annotation
	if err := row.Scan(
		&annotation.ID,
		&annotation.ImageID,
		&annotation.UserID,
		&annotation.X,
		&annotation.Y,
		&annotation.Width,
		&annotation.Height,
		&annotation.Comment,
		&annotation.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &annotation, nil
}

// Update modifies an existing annotation in the database. It returns an error if the update fails.
func (r *AnnotationRepository) Update(annotation *Annotation) error {
	query := `UPDATE annotations SET  x = $1, y = $2, width = $3, height = $4, comment = $5 WHERE id = $6`
	_, err := r.db.Exec(query,
		annotation.X,
		annotation.Y,
		annotation.Width,
		annotation.Height,
		annotation.Comment,
		annotation.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes an annotation from the database by its ID. It returns an error if the deletion fails.
func (r *AnnotationRepository) Delete(id string) error {
	query := `DELETE FROM annotations WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
