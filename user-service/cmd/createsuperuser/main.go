package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	login := flag.String("login", "", "superuser login")
	password := flag.String("password", "", "superuser password")
	flag.Parse()

	if *login == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "usage: createsuperuser -login <login> -password <password>")
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close(ctx)

	query := `
		INSERT INTO users (login, password_hash, is_super_user)
		VALUES ($1, $2, TRUE)
		ON CONFLICT (login) DO UPDATE SET is_super_user = TRUE
		RETURNING id`

	var id string
	err = conn.QueryRow(ctx, query, *login, string(hash)).Scan(&id)
	if err != nil {
		log.Fatalf("failed to create superuser: %v", err)
	}

	fmt.Printf("superuser created: id=%s login=%s\n", id, *login)
}
