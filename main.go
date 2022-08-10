package main

import (
	"log"
	"net/http"
	"tugas_mandiri/task"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", task.Home)
	mux.HandleFunc("/task", task.Index)
	mux.HandleFunc("/form", task.Form)
	mux.HandleFunc("/proses", task.Proses)
	mux.HandleFunc("/ubah", task.Ubah)
	mux.HandleFunc("/proses_ubah", task.ProsesUbah)
	mux.HandleFunc("/hapus", task.Hapus)

	fileServer := http.FileServer(http.Dir("style"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("port 8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
