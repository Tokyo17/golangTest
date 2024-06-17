package main

import (
	"golang/api"
	"net/http"
)

func main() {

	http.HandleFunc("/", api.Handler)
	http.HandleFunc("/2", api.Handler2) // Gunakan handler yang didefinisikan di index.go
	http.ListenAndServe(":8080", nil)   // Jalankan server HTTP pada port 8080
}
