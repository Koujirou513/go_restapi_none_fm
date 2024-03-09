package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"     // 文字列を整数に変換する
	"github.com/gorilla/mux"  // ルーティング機能を拡張するパッケージ
)

// Book型を定義
type Book struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
}

// Book型のスライスを作成
var books []Book
// idを定義
var id int = 1

// インデックスを探す関数
func findIndex(bookID string) int {  // bookIDを受け取りint型の値を返す
	id, err := strconv.Atoi(bookID)  // bookID文字列をintに変換
	if err != nil {                  // エラーがあれば-1を返す
		return -1
	}
	for i, book := range books {     // iにはインデックス、bookにはbooksスライスの各値が入る
		if book.ID == id {           // 検索しているidがあれば
			return i                 // インデックスを返す
		}
	}
	return -1                        // 無ければ-1を返す
}

// 全てのBookを取得する
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // レスポンスのContent-Typeを設定
	json.NewEncoder(w).Encode(books)  // booksをJSONにエンコードしてレスポンスボディに書き込む
}

// 特定のIDのBookを取得する
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)                    // URLパラメータを取得

	i := findIndex(params["id"])             // bookのインデックスを取得

	// Bookが見つからない場合はエラーを返す
	if i == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No book found with given ID")
		return
	}

	json.NewEncoder(w).Encode(books[i])  // 見つかったBookをJSONにエンコードしてレスポンスボディに書き込む
}

// 新しいBookを作成する
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book                              // Book型のbookを定義
	_ = json.NewDecoder(r.Body).Decode(&book)  // リクエストボディのJSONをBookにデコード
	book.ID = id                               // 新しいIDを割り当てる
	id++                                       // IDをインクリメント
	books = append(books, book)                // booksに新しいBookを追加
	json.NewEncoder(w).Encode(book)            // 作成されたBookをJSONにエンコードしてレスポンスボディに書き込み
}

// 特定のIDのBookを更新する
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	i := findIndex(params["id"])

	if i == -1 {                             // Bookが見つからない場合はエラーを返す
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No book found with given ID")
		return
	}

	var updatedBook Book                              // Book型のupdatedBook変数を定義
	_ = json.NewDecoder(r.Body).Decode(&updatedBook)  // リクエストボディのJSONをBookにデコード
	books[i] = updatedBook                            // 既存の書籍を新しい内容で更新 
	books[i].ID, _ = strconv.Atoi(params["id"])       // IDを保持,strconv.Atoiは2つの値を返す、2つ目はエラーが入る
	json.NewEncoder(w).Encode(books[i])
}   // エラーの値を_に代入し明示的に無視しているが本来は適切にエラーハンドリングをするべき

// 特定のIDのBookを削除する
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
} // appendの第二引数で複数展開したい場合は...を使用する

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