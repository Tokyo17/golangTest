package main

import (
	"golang/api" // Ubah `my-module` menjadi nama modul yang Anda tentukan
	"net/http"
)

func main() {
	http.HandleFunc("/", api.Handler) // Gunakan handler yang didefinisikan di index.go
	http.ListenAndServe(":8080", nil) // Jalankan server HTTP pada port 8080
}
