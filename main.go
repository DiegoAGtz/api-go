package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	ID          int    `json:"id"`
	ID_author   int    `json:"id_author"`
	ID_category int    `json:"id_category"`
	Title       string `json:"title"`
	Pages       int    `json:"pages"`
}

var db *sql.DB

func main() {
	// Conectarse a la base de datos MySQL
	var err error
	db, err = sql.Open("mysql", "root:3435@tcp(localhost:3306)/library")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Inicializar el enrutador "mux"
	r := mux.NewRouter()

	// Rutas de la API
	r.HandleFunc("/author", GetAuthors).Methods("GET")
	r.HandleFunc("/author/{id}", GetAuthor).Methods("GET")
	r.HandleFunc("/author", CreateAuthor).Methods("POST")
	r.HandleFunc("/author/{id}", UpdateAuthor).Methods("PUT")
	r.HandleFunc("/author/{id}", DeleteAuthor).Methods("DELETE")

	r.HandleFunc("/category", GetCategories).Methods("GET")
	r.HandleFunc("/category/{id}", GetCategory).Methods("GET")
	r.HandleFunc("/category", CreateCategory).Methods("POST")
	r.HandleFunc("/category/{id}", UpdateCategory).Methods("PUT")
	r.HandleFunc("/category/{id}", DeleteCategory).Methods("DELETE")

	r.HandleFunc("/", GetBooks).Methods("GET")
	r.HandleFunc("/{id}", GetBook).Methods("GET")
	r.HandleFunc("/", CreateBook).Methods("POST")
	r.HandleFunc("/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/{id}", DeleteBook).Methods("DELETE")

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor escuchando en el puerto 8030...")
	log.Fatal(http.ListenAndServe(":8030", r))
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM book")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.ID_author, &book.ID_category, &book.Title, &book.Pages)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	err := db.QueryRow("SELECT * FROM book WHERE id = ?", params["id"]).Scan(&book.ID, &book.ID_author, &book.ID_category, &book.Title, &book.Pages)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO book (id_author, id_category, title, pages) VALUES (?, ?, ?, ?)", book.ID_author, book.ID_category, book.Title, book.Pages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	book.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book.ID, err = strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE book SET id_author = ?, id_category = ?, title = ?, pages = ?  WHERE id = ?", book.ID_author, book.ID_category, book.Title, book.Pages, params["id"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rows != 1 {
		mapD := map[string]string{"error": "The book does not exist"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mapD)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	_, err := db.Exec("DELETE FROM book WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM author")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	authors := []Author{}
	for rows.Next() {
		var author Author
		err := rows.Scan(&author.ID, &author.Name)
		if err != nil {
			log.Fatal(err)
		}
		authors = append(authors, author)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}

func GetAuthor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var author Author
	err := db.QueryRow("SELECT * FROM author WHERE id = ?", params["id"]).Scan(&author.ID, &author.Name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

func CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author Author
	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO author (name) VALUES (?)", author.Name)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	author.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var author Author
	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	author.ID, err = strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE author SET name = ? WHERE id = ?", author.Name, params["id"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rows != 1 {
		mapD := map[string]string{"error": "The author does not exist"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mapD)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	_, err := db.Exec("DELETE FROM author WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM category")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var category Category
	err := db.QueryRow("SELECT * FROM category WHERE id = ?", params["id"]).Scan(&category.ID, &category.Name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO category (name) VALUES (?)", category.Name)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	category.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	category.ID, err = strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE category SET name = ? WHERE id = ?", category.Name, params["id"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rows != 1 {
		mapD := map[string]string{"error": "The category does not exist"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mapD)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	_, err := db.Exec("DELETE FROM category WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
