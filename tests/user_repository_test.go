package tests

import (
	"log"
	"os"
	"testing"

	"github.com/Agero19/AnnotateX-api/internal/config"
	"github.com/Agero19/AnnotateX-api/internal/db"
	"github.com/Agero19/AnnotateX-api/internal/repository"
	"github.com/Agero19/AnnotateX-api/internal/testutils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var repo repository.Repository

func TestMain(m *testing.M) {
	// connect to usual db
	_ = godotenv.Load()
	cfg := config.LoadConfig()
	test_cfg := config.LoadTestConfig()
	adminDB, err := db.New(
		cfg.DB.URL+"?sslmode=disable",
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)

	if err != nil {
		log.Fatalf("failed to connect to admin db: %v", err)
	}

	// create test db
	_, err = adminDB.Exec("DROP DATABASE IF EXISTS " + test_cfg.DB.Name)
	if err != nil {
		log.Fatalf("failed to drop test db: %v", err)
	}
	_, err = adminDB.Exec("CREATE DATABASE " + test_cfg.DB.Name)
	if err != nil {
		log.Fatalf("failed to create test db: %v", err)
	}
	adminDB.Close()

	// connect to test db
	testDB, err := db.New(
		test_cfg.DB.URL,
		test_cfg.DB.MaxIdleConns,
		test_cfg.DB.MaxOpenConns,
		test_cfg.DB.MaxIdleTime,
	)
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}
	repo = repository.NewRepository(testDB)

	// Run migrations
	if err := testutils.RunMigrations(
		test_cfg.DB.URL,
		"./../cmd/migrate/migrations",
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// run tests
	code := m.Run()
	if err := testDB.Close(); err != nil {
		log.Printf("failed to close test db connection: %v", err)
	}

	//connect to admin db to clean up test db
	adminDB, err = db.New(
		cfg.DB.URL+"?sslmode=disable",
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		log.Fatalf("failed to connect to admin db: %v", err)
	}
	defer adminDB.Close()

	if _, err := adminDB.Exec(
		`SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1 AND pid <> pg_backend_pid()`,
		test_cfg.DB.Name,
	); err != nil {
		log.Printf("failed to terminate connections: %v", err)
	}

	if _, err := adminDB.Exec("DROP DATABASE IF EXISTS " + test_cfg.DB.Name); err != nil {
		log.Printf("failed to drop test db: %v", err)
	}

	os.Exit(code)
}

func TestUserRepository_CRUD(t *testing.T) {
	var createdUser repository.User

	t.Run("Create", func(t *testing.T) {
		user := &repository.User{
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "secretpassword",
		}
		err := repo.Users.Create(user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		if user.ID == "" {
			t.Error("expected user ID to be set")
		}
		createdUser = *user
	})

	t.Run("GetAll", func(t *testing.T) {
		users, err := repo.Users.GetAll()
		if err != nil {
			t.Fatalf("failed to get all users: %v", err)
		}
		if len(users) == 0 {
			t.Error("expected at least one user")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		user, err := repo.Users.GetByID(createdUser.ID)
		if err != nil {
			t.Fatalf("failed to get user by id: %v", err)
		}
		if user == nil || user.ID != createdUser.ID {
			t.Errorf("expected to find user %s, got %+v", createdUser.ID, user)
		}
	})

	t.Run("Update", func(t *testing.T) {
		createdUser.Username = "updatedname"
		createdUser.Email = "updateduser@example.com"
		createdUser.Password = "newsecret"

		err := repo.Users.Update(&createdUser)
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}

		updated, _ := repo.Users.GetByID(createdUser.ID)
		if updated.Username != "updatedname" {
			t.Errorf("expected username to be updated, got %s", updated.Username)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Users.Delete(createdUser.ID)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}

		deleted, _ := repo.Users.GetByID(createdUser.ID)
		if deleted != nil {
			t.Errorf("expected user to be deleted")
		}
	})
}
