package main

import (
	"ebtek-frontend/auth"
	"ebtek-frontend/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("Sunucu başlatılıyor...")

	// Static dizini oluştur
	staticDir := filepath.Join("static", "images")
	err := os.MkdirAll(staticDir, 0755)
	if err != nil {
		log.Printf("Static dizin oluşturulurken hata: %v\n", err)
	}

	// Varsayılan resimleri kopyala
	defaultImages := map[string]string{
		"default-home.jpg": "https://images.unsplash.com/photo-1570129477492-45c003edd2be",
		"default-car.jpg":  "https://images.unsplash.com/photo-1568605117036-5fe5e7bab0b7",
	}

	for name, url := range defaultImages {
		filePath := filepath.Join(staticDir, name)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Varsayılan resim indiriliyor: %s\n", name)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Resim indirme hatası: %v\n", err)
				continue
			}
			defer resp.Body.Close()

			file, err := os.Create(filePath)
			if err != nil {
				log.Printf("Dosya oluşturma hatası: %v\n", err)
				continue
			}
			defer file.Close()
		}
	}

	// Route handlers
	mux := http.NewServeMux()

	// Ana sayfa yönlendirmesi
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Sadece kök dizin için yönlendirme yap
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Diğer tüm tanımsız URL'ler için 404
		http.NotFound(w, r)
	})

	// Static dosyaları serve et
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public routes
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/logout", auth.LogoutHandler)

	// Protected routes
	mux.HandleFunc("/dashboard", auth.RequireAuth(handlers.DashboardHandler))
	mux.HandleFunc("/profile", auth.RequireAuth(handlers.ProfileHandler))
	mux.HandleFunc("/listing/", auth.RequireAuth(handlers.ListingDetailHandler))
	mux.HandleFunc("/new-listing", auth.RequireAuth(handlers.NewListingHandler))
	mux.HandleFunc("/api/listings", auth.RequireAuth(handlers.ListingHandler))

	port := ":8080"

	// Sunucu ayarları
	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Sunucu http://localhost%s adresinde başlatılıyor...\n", port)
	fmt.Printf("Ana sayfa: http://localhost%s/\n", port)
	fmt.Printf("Login sayfası: http://localhost%s/login\n", port)

	log.Fatal(server.ListenAndServe())
}
