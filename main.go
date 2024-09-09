package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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
	StoreID     int    `json:"store_id"` // TODO: Add store_id to Click struct
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Address  string `json:"address"`
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
        image BLOB NOT NULL,
        owner_id INTEGER NOT NULL,
        FOREIGN KEY (owner_id) REFERENCES users(id)
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
		FOREIGN KEY (store_id) REFERENCES stores(id)
    );`

	// Searches for Stores
	Searches := `
	CREATE TABLE IF NOT EXISTS searches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		search_query TEXT NOT NULL,
		time_of_search DATE NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	SearchesItems := `
	CREATE TABLE IF NOT EXISTS searches_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		store_id INTEGER NOT NULL,
		query TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (store_id) REFERENCES stores(id)
	);`

	statements := []string{Users, Items, Stores, Orders, Clicks, Images, Searches, SearchesItems}
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

	// Routes
	e.POST("/register", registerUser)
	e.POST("/login", loginUser)

	e.GET("/stores", getStores)
	e.POST("/click", recordClick)

	// Route for fetching item images
	e.GET("/images/:item_id", getItemImages)

	// Route for fetching statistics for a specific store
	e.GET("/statistics/:store_name", getStatistics)

	// Route for fetching items for a specific store
	e.GET("/store/:store_name", getItems)

	// Route for recording search queries
	e.POST("/search", recordSearch)

	e.POST("/searchItems", recordSearchItem)

	e.Start(":8080")
}

// jwtSecret represents the secret key used for signing the JWT
var jwtSecret = []byte("your-secret-key")

// jwtCustomClaims struct now includes custom fields along with the standard claims.
type jwtCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Helper function to generate JWT
func generateJWT(userID int) (string, error) {
	// Create the claims with custom fields (UserID) and the standard claims
	claims := &jwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Updated to use NewNumericDate in v5
		},
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(jwtSecret)
}

// User registration handler
func registerUser(c echo.Context) error {
	// write the request data
	fmt.Println(c.Request().Body)
	u := new(User)
	if err := c.Bind(u); err != nil {
		log.Println("Error binding request data:", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Log received user data for debugging
	log.Printf("User registration attempt: Username=%s, Email=%s\n", u.Username, u.Email)

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.JSON(http.StatusInternalServerError, "Failed to hash password")
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO users (username, password, mail, address) VALUES (?, ?, ?, ?)",
		u.Username, hashedPassword, u.Email, u.Address)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		return c.JSON(http.StatusInternalServerError, "Failed to register user")
	}

	log.Println("User registered successfully")
	return c.JSON(http.StatusCreated, "User registered successfully")
}

// User login handler
func loginUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Fetch user from the database by username or email
	var id int
	var hashedPassword string
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ? OR mail = ?", u.Username, u.Email).
		Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, "Invalid username or email")
		}
		return c.JSON(http.StatusInternalServerError, "Failed to query user")
	}

	// Compare the stored hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid password")
	}

	// Generate JWT
	token, err := generateJWT(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to generate token")
	}

	// Respond with JWT token
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   token,
		"user_id": id,
	})
}

// Retrieve click data for a specific store from the database
func getClickData(storeName string) ([]Click, error) {
	query := `
		SELECT clicks.id, clicks.user_id, clicks.item_id, clicks.time_of_click 
		FROM clicks 
		INNER JOIN items ON clicks.item_id = items.id 
		INNER JOIN stores ON items.store_id = stores.id 
		WHERE stores.name = ?`

	rows, err := db.Query(query, storeName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []Click
	for rows.Next() {
		var click Click
		if err := rows.Scan(&click.ID, &click.UserID, &click.ItemID, &click.TimeOfClick); err != nil {
			continue
		}
		clicks = append(clicks, click)
	}
	return clicks, nil
}

func getStatistics(c echo.Context) error {
	// Extract store name from the URL path
	storeName := c.Param("store_name")

	clicks, err := getClickData(storeName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Return the list of clicks as a JSON response
	return c.JSON(http.StatusOK, clicks)
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

	// Find the store ID based on the store name
	var storeID int
	for _, store := range stores {
		if store.Name == storeName {
			storeID = store.ID
			break
		}
	}

	userID := c.QueryParam("user_id") // Retrieve user_id from query parameters
	itemID := c.QueryParam("item_id") // Retrieve item_id from query parameters
	if userID == "" || itemID == "" {
		return c.JSON(http.StatusBadRequest, "User ID and Item ID are required")
	}

	_, err := db.Exec("INSERT INTO clicks (user_id, item_id, time_of_click, store_id) VALUES (?, ?, ?, ?)", userID, itemID, time.Now().Format("2006-01-02 15:04:05"), storeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Click recorded successfully")
}

func recordSearch(c echo.Context) error {
	// Define a struct to hold the request data
	var requestData struct {
		UserID    string `json:"user_id"`
		StoreName string `json:"store_name"`
	}

	// Bind the JSON payload to the requestData struct
	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// Extract values from requestData
	userID := requestData.UserID
	storeName := requestData.StoreName

	fmt.Println("User ID:", userID)
	fmt.Println("Store Name:", storeName)

	// Insert the search record into the database
	_, err := db.Exec("INSERT INTO searches (user_id, search_query, time_of_search) VALUES (?, ?, ?)", userID, storeName, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Search recorded successfully")
}

type SearchItem struct {
	UserID    string `json:"user_id"`
	StoreName string `json:"store_name"`
	Query     string `json:"query"`
}

func recordSearchItem(c echo.Context) error {
	// Define a struct to hold the request data
	requestData := SearchItem{}

	// Bind the JSON payload to the requestData struct
	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	// Extract values from requestData
	userID := requestData.UserID
	storeName := requestData.StoreName
	query := requestData.Query

	// Convert store name to store ID
	var storeID int
	storeFound := false
	for _, store := range stores {
		if store.Name == storeName {
			storeID = store.ID
			storeFound = true
			break
		}
	}

	fmt.Println("User ID:", userID)
	fmt.Println("Store ID:", storeID)
	fmt.Println("Query:", query)

	if userID == "" || !storeFound || query == "" {
		return c.JSON(http.StatusBadRequest, "User ID, Store ID, and Query are required")
	}

	// Insert the search record into the database
	_, err := db.Exec("INSERT INTO searches_items (user_id, store_id, query) VALUES (?, ?, ?)", userID, storeID, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Search recorded successfully")
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
