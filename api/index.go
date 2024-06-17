package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"` // Password will not be included in JSON response
}

type Response struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

const connStr = "postgresql://Tokyo17:pnm2fY6awAjE@ep-royal-sun-104233.us-east-2.aws.neon.tech/neondb?sslmode=require"

func Handler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	db, err := connectDB()
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodGet:
		handleGetUsers(w, db)
	case http.MethodPost:
		handleCreateUser(w, r, db)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
}

func connectDB() (*sql.DB, error) {
	return sql.Open("postgres", connStr)
}

func handleGetUsers(w http.ResponseWriter, db *sql.DB) {
	users, err := fetchUsers(db)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err := writeJSONResponse(w, users); err != nil {
		httpError(w, err, http.StatusInternalServerError)
	}
}

func fetchUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM \"user\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}
	fmt.Println(user.Name, user)
	if user.Name == "" || user.Email == "" || user.Password == "" {
		httpError(w, errors.New("name, email, and password cannot be empty"), http.StatusBadRequest)
		return
	}

	if err := createUser(db, &user); err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "User successfully created",
		User:    user,
	}

	w.WriteHeader(http.StatusCreated)
	if err := writeJSONResponse(w, response); err != nil {
		httpError(w, err, http.StatusInternalServerError)
	}
}

func createUser(db *sql.DB, user *User) error {
	err := db.QueryRow(
		"INSERT INTO \"user\" (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password,
	).Scan(&user.ID)
	return err
}

func httpError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
