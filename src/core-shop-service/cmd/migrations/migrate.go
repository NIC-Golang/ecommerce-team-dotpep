package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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
	db, err := sql.Open("postgres", postgres_str)
	if err != nil {
		return fmt.Errorf("error with connecting to postgres")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error occured:%s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("error pinging postgres: %v", err)
	}
	log.Println("Connected to PostgreSQL successfully")
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"project",
		driver)
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
