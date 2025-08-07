package database

import (
    "database/sql"
    "log"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    connStr := "postgres://username:password@localhost/ebtek?sslmode=disable"
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal(err)
    }

    // Tabloları oluştur
    createTables()
}

func createTables() {
    // Listings tablosu
    _, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS listings (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            price DECIMAL(10,2) NOT NULL,
            type VARCHAR(50) NOT NULL,
            location VARCHAR(255),
            images TEXT[],
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Properties tablosu
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS properties (
            listing_id INTEGER PRIMARY KEY REFERENCES listings(id),
            area DECIMAL(10,2),
            rooms INTEGER,
            floor INTEGER,
            total_floors INTEGER,
            age INTEGER,
            heating VARCHAR(100),
            features TEXT[]
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Vehicles tablosu
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS vehicles (
            listing_id INTEGER PRIMARY KEY REFERENCES listings(id),
            brand VARCHAR(100),
            model VARCHAR(100),
            year INTEGER,
            kilometer INTEGER,
            fuel_type VARCHAR(50),
            transmission VARCHAR(50),
            color VARCHAR(50),
            features TEXT[]
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Users tablosu
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(100) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            phone VARCHAR(20),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Messages tablosu
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id SERIAL PRIMARY KEY,
            sender_id INTEGER REFERENCES users(id),
            receiver_id INTEGER REFERENCES users(id),
            listing_id INTEGER REFERENCES listings(id),
            message TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            read_at TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Favorites tablosu
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS favorites (
            user_id INTEGER REFERENCES users(id),
            listing_id INTEGER REFERENCES listings(id),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (user_id, listing_id)
        )
    `)
    if err != nil {
        log.Fatal(err)
    }
}
