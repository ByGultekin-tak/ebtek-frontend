package main

import (
	"fmt"
	"log"
	"net/http"
	"ebtek-frontend/auth"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hoş geldiniz! Bu korumalı dashboard sayfasıdır.")
}

func main() {
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/dashboard", auth.RequireAuth(dashboard))
	
	fmt.Println("Sunucu http://localhost:8080 adresinde başlatılıyor...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
}
