package main

import (
    "database/sql"
    "html/template"
    "io"
    "log"
    "net/http"

    "github.com/labstack/echo/v4"
    _ "github.com/mattn/go-sqlite3"
)

// Template struct to handle rendering
type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("sqlite3", "./users.db")
    if err != nil {
        log.Fatal(err)
    }

    createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

    if _, err := db.Exec(createTable); err != nil {
        log.Fatal(err)
    }
}

// SaveUser saves a new user to the database
func SaveUser(username, password string) error {
    stmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
    if err != nil {
        return err
    }
    _, err = stmt.Exec(username, password)
    return err
}

// CheckUser checks if a user exists and the password matches
func CheckUser(username, password string) bool {
    var dbPassword string
    err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
    if err != nil {
        return false
    }
    return dbPassword == password
}

func main() {
    e := echo.New()

    // Setting up the template renderer
    t := &Template{
        templates: template.Must(template.ParseGlob("*.html")),
    }
    e.Renderer = t

    e.GET("/", serveIndex)
    e.POST("/login", loginHandler)
    e.POST("/register", registerHandler)

    e.Start(":8080")
}

// Serve the HTML page
func serveIndex(c echo.Context) error {
    return c.Render(http.StatusOK, "index.html", nil)
}

// Handle login request
func loginHandler(c echo.Context) error {
    username := c.FormValue("username")
    password := c.FormValue("password")

    var response map[string]string

    if CheckUser(username, password) {
        response = map[string]string{"message": "Login successful!"}
    } else {
        response = map[string]string{"message": "Invalid username or password."}
    }

    return c.JSON(http.StatusOK, response)
}

// Handle registration request
func registerHandler(c echo.Context) error {
    username := c.FormValue("username")
    password := c.FormValue("password")

    var response map[string]string

    if err := SaveUser(username, password); err != nil {
        response = map[string]string{"message": "Error registering user."}
    } else {
        response = map[string]string{"message": "User registered successfully!"}
    }

    return c.JSON(http.StatusOK, response)
}
