package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"
)

type User struct {
	ID   int
	Name string
}

func queryUsers(db *sql.DB) {
	rows, error := db.Query("SELECT * FROM users")
	if error != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", error)
		os.Exit(1)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if error := rows.Scan(&user.ID, &user.Name); error != nil {
			fmt.Println("Error scanning row:", error)
			return
		}

		users = append(users, user)
		fmt.Println(user.ID, user.Name)
	}

	if error := rows.Err(); error != nil {
		fmt.Println("Error during rows iteration:", error)
	}
}

func main() {
	error := godotenv.Load()
	if error != nil {
		fmt.Println("Error while loading .env file", error)
	}

	dbName := "test-db"
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	directory, error := os.MkdirTemp("", "libsql-*")
	if error != nil {
		fmt.Println("Error creating temporary directory:", error)
		os.Exit(1)
	}
	defer os.RemoveAll(directory)

	dbPath := filepath.Join(directory, dbName)

	connector, error := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl, libsql.WithAuthToken(authToken))
	if error != nil {
		fmt.Println("Error creating connector:", error)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)

	queryUsers(db)
	defer db.Close()
}
