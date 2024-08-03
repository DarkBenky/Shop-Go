package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	StoreID     int    `json:"store_id"`
}

type Store struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

var db *sql.DB
var stores = []Store{}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	Users := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        mail TEXT NOT NULL,
        address TEXT
    );`

	Items := `
    CREATE TABLE IF NOT EXISTS items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price INTEGER NOT NULL,
        description TEXT NOT NULL,
		category TEXT,
        images BLOB,
        store_id INTEGER,
        FOREIGN KEY (store_id) REFERENCES stores(id)
    );`

	Stores := `
    CREATE TABLE IF NOT EXISTS stores (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        address TEXT NOT NULL
    );`

	Orders := `
    CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        store_id INTEGER,
        item_id INTEGER NOT NULL,
        quantity INTEGER NOT NULL,
        time_of_purchase DATE NOT NULL,
        estimated_delivery DATE,
        status TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (store_id) REFERENCES stores(id),
        FOREIGN KEY (item_id) REFERENCES items(id)
    );`

	Clicks := `
    CREATE TABLE IF NOT EXISTS clicks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        item_id INTEGER NOT NULL,
        time_of_click DATE NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (item_id) REFERENCES items(id)
    );`

	statements := []string{Users, Items, Stores, Orders, Clicks}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			log.Fatal(err)
		}
	}

	// if database is not empty, fetch all stores
	rows, err := db.Query("SELECT id, name, address FROM stores")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var store Store
		if err := rows.Scan(&store.ID, &store.Name, &store.Address); err != nil {
			log.Fatal(err)
		}
		stores = append(stores, store)
	}
}

func main() {
	e := echo.New()

	// Enable CORS for all origins
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	e.GET("/stores", getStores)

	// generate roots for different shops
	for _, store := range stores {
		e.GET("/"+store.Name, getItems)
		println(store.Name)
	}

	e.Start(":8080")
}

func getStores(c echo.Context) error {
	return c.JSON(http.StatusOK, stores)
}

func getItems(c echo.Context) error {
	items, err := fetchItems()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, items)
}

func fetchItems() ([]Item, error) {
	rows, err := db.Query("SELECT id, name, price, description, category, images FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		var imageData string
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.Category, &imageData); err != nil {
			return nil, err
		}

		if len(imageData) > 0 {
			// encodedImage := base64.StdEncoding.EncodeToString(imageData)
			item.Image = imageData
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
