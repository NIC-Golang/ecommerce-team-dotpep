package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func DBconn() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	password := os.Getenv("SQL_PASS")
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	connStr := fmt.Sprintf("user=Fiveret password=%s dbname=project sslmode=disable", password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Print("Seccessfully connected to PostGRE !!!")

}

func GetDBConnection(ctx context.Context) (*pgx.Conn, error) {
	host := os.Getenv("IP_SQL")
	password := os.Getenv("SQL_PASS")
	port := os.Getenv("PORT_SQL")
	if port == "" {
		port = "5432"
	}
	connStr := fmt.Sprintf("postgres://fiveret:%s@%s:%s/project", password, host, port)
	return pgx.Connect(ctx, connStr)
}
