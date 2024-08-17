package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Discount    int    `json:"discount"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	StoreID     int    `json:"store_id"`
}

type Store struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Category string `json:"category"`
	Image    string `json:"image"`
}

type Click struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	ItemID      int    `json:"item_id"`
	TimeOfClick string `json:"time_of_click"`
}

var db *sql.DB
var stores = []Store{}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	Images := `
	CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		image BLOB NOT NULL,
		store_id INTEGER NOT NULL,
		item_id INTEGER NOT NULL,
		FOREIGN KEY (store_id) REFERENCES stores(id),
		FOREIGN KEY (item_id) REFERENCES items(id)
	);`

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
		discount INTEGER,
		category TEXT,
        images BLOB,
        store_id INTEGER,
        FOREIGN KEY (store_id) REFERENCES stores(id)
    );`

	Stores := `
    CREATE TABLE IF NOT EXISTS stores (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        address TEXT NOT NULL,
		category TEXT NOT NULL,
		image BLOB NOT NULL
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

	statements := []string{Users, Items, Stores, Orders, Clicks, Images}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			log.Fatal(err)
		}
	}

	// if database is not empty, fetch all stores
	rows, err := db.Query("SELECT id, name, address, category, image FROM stores")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var store Store
		if err := rows.Scan(&store.ID, &store.Name, &store.Address, &store.Category, &store.Image); err != nil {
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
	e.POST("/click", recordClick)

	// Route for fetching item images
	e.GET("/images/:item_id", getItemImages)

	// Route for fetching items for a specific store
	e.GET("/store/:store_name", getItems)

	e.Start(":8080")
}

func getStores(c echo.Context) error {
	return c.JSON(http.StatusOK, stores)
}

func getItems(c echo.Context) error {
	// Extract store name from the URL path
	storeName := c.Param("store_name")

	// Find the store ID based on the store name
	var storeID int
	for _, store := range stores {
		if store.Name == storeName {
			storeID = store.ID
			break
		}
	}

	// If no matching store is found, return an error
	if storeID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Store not found"})
	}

	// Fetch items for the specific store
	items, err := fetchItems(storeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

func fetchItems(storeId int) ([]Item, error) {
	rows, err := db.Query("SELECT id, name, price, description, discount, category, images FROM items WHERE store_id = ?", storeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		var imageData string
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.Discount, &item.Category, &imageData); err != nil {
			return nil, err
		}

		if len(imageData) > 0 {
			item.Image = imageData
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func recordClick(c echo.Context) error {
	storeName := c.QueryParam("store_name") // Retrieve store_name from query parameters
	if storeName == "" {
		return c.JSON(http.StatusBadRequest, "Store name is required")
	}

	userID := c.QueryParam("user_id") // Retrieve user_id from query parameters
	itemID := c.QueryParam("item_id") // Retrieve item_id from query parameters
	if userID == "" || itemID == "" {
		return c.JSON(http.StatusBadRequest, "User ID and Item ID are required")
	}

	// log.Println("Store Name:", storeName)
	// log.Println("User ID:", userID)
	// log.Println("Item ID:", itemID)

	_, err := db.Exec("INSERT INTO clicks (user_id, item_id, time_of_click) VALUES (?, ?, ?)", userID, itemID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Click recorded successfully")
}

func getItemImages(c echo.Context) error {
	itemID := c.Param("item_id")
	if itemID == "" {
		return c.JSON(http.StatusBadRequest, "Item ID is required")
	}

	var images []string

	// Fetch multiple image data
	rows, err := db.Query("SELECT image FROM images WHERE item_id = ?", itemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var imageData string
		if err := rows.Scan(&imageData); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		images = append(images, imageData)
	}

	// Check if any images were found
	if len(images) == 0 {
		return c.JSON(http.StatusNotFound, "No images found")
	}

	// Return the images data as a JSON array
	return c.JSON(http.StatusOK, map[string][]string{
		"images": images,
	})
}
