package models

import "time"

type ListingType string

const (
	PropertyType ListingType = "property"
	VehicleType  ListingType = "vehicle"
)

type Listing struct {
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Type        ListingType `json:"type"`
	Location    string      `json:"location"`
	Images      []string    `json:"images"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Property struct {
	ListingID   int64    `json:"listing_id"`
	Area        float64  `json:"area"`         // Metrekare
	Rooms       int      `json:"rooms"`        // Oda sayısı
	Floor       int      `json:"floor"`        // Bulunduğu kat
	TotalFloors int      `json:"total_floors"` // Toplam kat
	Age         int      `json:"age"`          // Bina yaşı
	Heating     string   `json:"heating"`      // Isıtma tipi
	Features    []string `json:"features"`     // Özellikler (Asansör, Otopark vb.)
}

type Vehicle struct {
	ListingID    int64    `json:"listing_id"`
	Brand        string   `json:"brand"`        // Marka
	Model        string   `json:"model"`        // Model
	Year         int      `json:"year"`         // Yıl
	Kilometer    int      `json:"kilometer"`    // Kilometre
	FuelType     string   `json:"fuel_type"`    // Yakıt tipi
	Transmission string   `json:"transmission"` // Vites
	Color        string   `json:"color"`        // Renk
	Features     []string `json:"features"`     // Özellikler
}
