package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"     // 文字列を整数に変換する
	"github.com/gorilla/mux"  // ルーティング機能を拡張するパッケージ
)

type Book struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
}

var books []Book
var id int = 1

func findIndex(bookID string) int {
	id, err := strconv.Atoi(bookID)
	if err != nil {
		return -1
	}
	for i, book := range books {
		if book.ID == id {
			return i
		}
	}
	return -1
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // レスポンスのContent-Typeを設定
	json.NewEncoder(w).Encode(books)  // booksをJSONにエンコードしてレスポンスボディに書き込む
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // URLパラメータを取得

	i := findIndex(params["id"])  // bookのインデックスを取得

	// Bookが見つからない場合はエラーを返す
	if i == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No book found with given ID")
		return
	}

	json.NewEncoder(w).Encode(books[i])  // 見つかったBookをJSONにエンコードしてレスポンスボディに書き込む
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)  // リクエストボディのJSONをBookにデコード
	book.ID = id  // 新しいIDを割り当て
	id++  // IDをインクリメント
	books = append(books, book)  // booksに新しいBookを追加
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	i := findIndex(params["id"])

	if i == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No book found with given ID")
		return
	}

	var updatedBook Book
	_ = json.NewDecoder(r.Body).Decode(&updatedBook)  // リクエストボディのJSONをBookにデコード
	books[i] = updatedBook  // 既存の書籍を新しい内容で更新 
	books[i].ID, _ = strconv.Atoi(params["id"])  // IDを保持
	json.NewEncoder(w).Encode(books[i])
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)  // URLパラメータ取得

	i := findIndex(params["id"])  // Bookのインデックスを取得

	// Bookが見つからない場合は400エラーを返す
	if i == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No book found with given ID")
		return
	}

	books = append(books[:i], books[i+1:]...)  // Bookをbooksから削除
}

func main() {
	r := mux.NewRouter()  // ルーターを作成

	// エンドポイントとハンドラ関数を関連付け
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Println("Server starting at 8000")
	http.ListenAndServe(":8000", r)
}