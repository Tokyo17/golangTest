package api

import (
	"fmt"
	"net/http"
)

func Handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go 2")
}
