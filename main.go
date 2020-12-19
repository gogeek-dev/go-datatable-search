package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Book Struct contains useful information
type Book struct {
	Title   string `json:"Title"`
	Country string `json:"Country"`
	Date    string `json:"Date"`
	Author  string `json:"Author"`
	ID      string `json:"identifier"`
}

func dbConn() (db *sql.DB) {

	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// Books is a list of Book
type Books struct {
	Info []Book `json:"books"`
}

//search function returns the result of the query
func search(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var book Book
	t, _ := template.ParseFiles("html/index.html")
	// result := search(`Select * From Book;`)

	rows, err := db.Query("select *from Book")
	if err != nil {
		log.Print(err.Error())
	}

	res := []Book{}
	for rows.Next() {
		// get RawBytes from data
		var id string
		var title, author, country, date string
		err = rows.Scan(&id, &title, &author, &country, &date)
		if err != nil {
			log.Println(err.Error())
		}

		book.ID = id
		book.Author = author
		book.Date = date
		book.Country = country
		book.Title = title

		res = append(res, book)
	}
	t.Execute(w, res)
	defer db.Close()

}

func main() {

	err1 := godotenv.Load(".env")

	if err1 != nil {
		log.Fatalf("Error loading .env file")
	}

	fmt.Println("Server started on: http://localhost:8080")
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", search)
	fileServer := http.FileServer(http.Dir("html/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Println(" -  Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(" -  ListenAndServe: ", err)
	}
}
