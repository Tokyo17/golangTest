package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	connStr := "postgresql://Tokyo17:pnm2fY6awAjE@ep-royal-sun-104233.us-east-2.aws.neon.tech/neondb?sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM \"user\" ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert users slice to JSON
	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJSON)
}
