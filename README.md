# go_restapi_none_fm

go run main.go　でサーバー立ち上げ

確認方法
GET /books - 全ての本の情報を取得
curl http://localhost:8000/books

GET /books/{id} - 特定のIDを持つ本の情報を取得
curl http://localhost:8000/books/1

POST /books - 新しい本を追加
curl -X POST http://localhost:8000/books -H "Content-Type: application/json" -d '{"title":"New Book","author":"Author Name"}'

PUT /books/{id} - 特定のIDを持つ本の情報を更新
curl -X PUT http://localhost:8000/books/1 -H "Content-Type: application/json" -d '{"title":"Updated Title","author":"Updated Author"}'

DELETE /books/{id} - 特定のIDを持つ本を削除
curl -X DELETE http://localhost:8000/books/1
