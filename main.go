package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int
	Username string
	Password string
	Mail     string
	Address  string
}

type Item struct {
	ID          int
	Name        string
	Price       int
	Description string
	Image       string
	StoreID     int
}

type Store struct {
	ID      int
	Name    string
	Address string
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	// Ensure static directories exist
	err = os.MkdirAll("./static/images", os.ModePerm)
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
}

func main() {
	e := echo.New()
	e.Static("/static", "./static")

	// Setting up the template renderer
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}

	e.GET("/", serveItems)

	e.Start(":8080")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func serveItems(c echo.Context) error {
	items, err := fetchItems()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "items.html", map[string]interface{}{
		"Items": items,
	})
}

func fetchItems() ([]Item, error) {
	rows, err := db.Query("SELECT id, name, price, description, images FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		var imageData []byte
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &imageData); err != nil {
			return nil, err
		}

		// Log the length of the image data
		// log.Printf("Image data length: %d", len(imageData))

		if len(imageData) > 0 {
			// Save image as a file
			imagePath := "./static/images/item_" + strconv.Itoa(item.ID) + ".jpg"
			err = saveImageToFile(imagePath, imageData)
			if err != nil {
				return nil, err
			}
			item.Image = "/static/images/item_" + strconv.Itoa(item.ID) + ".jpg"
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func saveImageToFile(filePath string, imageData []byte) error {
	err := os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return err
	}
	return nil
}
