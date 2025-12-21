package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	dsn := "postgres://mockpay:mockpay@localhost:5433/mockpay?sslmode=disable"

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := db.Exec(
			context.Background(),
			"INSERT INTO users (username, email) VALUES ($1, $2)",
			u.Username,
			u.Email,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	log.Println("Users service started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
