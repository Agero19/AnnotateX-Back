package repository

import "database/sql"

// User represents a user in the database - Model
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

// UserRepository is a struct that provides methods to interact with the user database table. Implements the Users interface.
type UserRepository struct {
	db *sql.DB
}

// Create inserts a new user into the database. It returns an error if the insertion fails.
func (r *UserRepository) Create(user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) Returning id, created_at`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

// GetAll retrieves all users from the database.
func (r *UserRepository) GetAll() ([]*User, error) {
	query := `SELECT id, username, email, created_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// GetByID retrieves a user by their ID from the database. Does not return an error if the user is not found.
func (r *UserRepository) GetByID(id string) (*User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update modifies an existing user in the database. It returns an error if the update fails.
func (r *UserRepository) Update(user *User) error {
	query := `UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4`
	_, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.ID,
	)
	return err
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
