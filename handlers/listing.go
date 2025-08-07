package handlers

import (
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "strconv"
    "strings"
    "time"
    "ebtek-frontend/models"
    "ebtek-frontend/auth"
)

var listings = []models.Listing{
    {
        ID:          1,
        UserID:      1,
        Title:       "Lüks Daire",
        Description: "3+1 Deniz Manzaralı",
        Price:       1500000,
        Type:        models.PropertyType,
        Location:    "İstanbul, Kadıköy",
        Images:      []string{"/static/images/default-home.jpg"},
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    },
    {
        ID:          2,
        UserID:      1,
        Title:       "Sıfır BMW",
        Description: "2023 Model, Sıfır KM",
        Price:       2500000,
        Type:        models.VehicleType,
        Location:    "İstanbul, Beşiktaş",
        Images:      []string{"/static/images/default-car.jpg"},
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    },
    {
        ID:          3,
        UserID:      1,
        Title:       "Kiralık Ofis",
        Description: "200m² Plaza Ofis",
        Price:       15000,
        Type:        models.PropertyType,
        Location:    "İstanbul, Maslak",
        Images:      []string{"/static/images/default-home.jpg"},
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    },
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := auth.Store.Get(r, "session-name")
    username := session.Values["username"].(string)

    data := struct {
        Title    string
        Username string
        Listings []models.Listing
    }{
        Title:    "Ana Sayfa",
        Username: username,
        Listings: listings,
    }

    tmpl, err := template.ParseFiles(
        "templates/layout.html",
        "templates/index.html",
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func ListingHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // Tüm ilanları listele veya ID'ye göre tek ilan getir
        if id := r.URL.Query().Get("id"); id != "" {
            getSingleListing(w, r, id)
            return
        }
        getAllListings(w, r)
    
    case http.MethodPost:
        // Yeni ilan oluştur
        createListing(w, r)
    
    case http.MethodPut:
        // İlan güncelle
        updateListing(w, r)
    
    case http.MethodDelete:
        // İlan sil
        deleteListing(w, r)
    
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getAllListings(w http.ResponseWriter, r *http.Request) {
    // Filtreleme parametrelerini al
    listingType := r.URL.Query().Get("type")
    minPriceStr := r.URL.Query().Get("min_price")
    maxPriceStr := r.URL.Query().Get("max_price")
    location := r.URL.Query().Get("location")

    // Fiyat filtrelerini dönüştür
    var minPrice, maxPrice float64
    if minPriceStr != "" {
        if p, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
            minPrice = p
        }
    }
    if maxPriceStr != "" {
        if p, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
            maxPrice = p
        }
    }

    // Filtreleme işlemleri
    filteredListings := make([]models.Listing, 0)
    for _, l := range listings {
        if listingType != "" && string(l.Type) != listingType {
            continue
        }
        if minPrice > 0 && l.Price < minPrice {
            continue
        }
        if maxPrice > 0 && l.Price > maxPrice {
            continue
        }
        if location != "" && !strings.Contains(strings.ToLower(l.Location), strings.ToLower(location)) {
            continue
        }
        filteredListings = append(filteredListings, l)
    }

    json.NewEncoder(w).Encode(filteredListings)
}

func getSingleListing(w http.ResponseWriter, r *http.Request, idStr string) {
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Gerçek uygulamada veritabanından çekilecek
    for _, listing := range listings {
        if listing.ID == id {
            json.NewEncoder(w).Encode(listing)
            return
        }
    }

    http.Error(w, "Listing not found", http.StatusNotFound)
}

func createListing(w http.ResponseWriter, r *http.Request) {
    var listing models.Listing
    if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Gerçek uygulamada veritabanına kaydedilecek
    listings = append(listings, listing)
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(listing)
}

func updateListing(w http.ResponseWriter, r *http.Request) {
    var updatedListing models.Listing
    if err := json.NewDecoder(r.Body).Decode(&updatedListing); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Gerçek uygulamada veritabanında güncellenecek
    for i, listing := range listings {
        if listing.ID == updatedListing.ID {
            listings[i] = updatedListing
            json.NewEncoder(w).Encode(updatedListing)
            return
        }
    }

    http.Error(w, "Listing not found", http.StatusNotFound)
}

func deleteListing(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Gerçek uygulamada veritabanından silinecek
    for i, listing := range listings {
        if listing.ID == id {
            listings = append(listings[:i], listings[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }

    http.Error(w, "Listing not found", http.StatusNotFound)
}

func NewListingHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles(
        "templates/layout.html",
        "templates/new-listing.html",
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Get username from session for layout
    session, _ := auth.Store.Get(r, "session-name")
    username, _ := session.Values["username"].(string)

    if r.Method == http.MethodGet {
        data := struct {
            Title    string
            Username string
        }{
            Title:    "Yeni İlan",
            Username: username,
        }
        tmpl.Execute(w, data)
        return
    }
}

func ListingDetailHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles(
        "templates/layout.html",
        "templates/listing-detail.html",
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // URL'den ID'yi al
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 3 {
        http.Error(w, "Invalid listing ID", http.StatusBadRequest)
        return
    }
    
    idStr := pathParts[2]
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid listing ID", http.StatusBadRequest)
        return
    }

    // İlanı bul
    var listing models.Listing
    found := false
    for _, l := range listings {
        if l.ID == id {
            listing = l
            found = true
            break
        }
    }

    if !found {
        http.Error(w, "Listing not found", http.StatusNotFound)
        return
    }

    // Session'dan kullanıcı adını al
    session, _ := auth.Store.Get(r, "session-name")
    username, _ := session.Values["username"].(string)

    // Template verisi
    data := struct {
        models.Listing
        UserName string
    }{
        Listing: listing,
        UserName: username,
    }

    tmpl.Execute(w, data)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
    // Şimdilik basit bir yanıt döndürelim
    fmt.Fprintf(w, "Profil sayfası yakında eklenecek")
}
