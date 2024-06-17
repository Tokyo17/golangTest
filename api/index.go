package api

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	connStr := "postgresql://Tokyo17:pnm2fY6awAjE@ep-royal-sun-104233.us-east-2.aws.neon.tech/neondb?sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var version string
	if err := db.QueryRow("select version()").Scan(&version); err != nil {
		panic(err)
	}

	// fmt.Printf(w,"version=%s\n", version)

	fmt.Fprintf(w, "Hello from Go 1 %s", version)
}
