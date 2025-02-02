package migrations

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func RunMigrations() error {
	err := godotenv.Load("/app/.env")
	if err != nil {
		return fmt.Errorf("error downloading env. file")
	}
	pass, user := os.Getenv("SQL_PASS"), os.Getenv("SQL_USER")
	postgres_str := fmt.Sprintf("postgres://%s:%s@postgres:5432/project?sslmode=disable", user, pass)
	m, err := migrate.New(
		"file://./migrations",
		postgres_str)
	if err != nil {
		return fmt.Errorf("error creating migration object: %v", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("error applying migrations: %v", err)
	}
	fmt.Println("Migrations applied successfully.")
	return nil
}
