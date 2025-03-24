package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

// Task は、タスク情報を表す構造体
type Task struct {
	ID          int        `json:"id"`           // タスクのID
	Body        string     `json:"body"`         // タスク本文
	TimeLimit   *time.Time `json:"time_limit"`   // タスクの期限
	CompletedAt *time.Time `json:"completed_at"` // タスク完了時間
	CreatedAt   time.Time  `json:"created_at"`   //タスクの作成時間
}

// TaskToCreate は、作成するタスクの構造体
type TaskToCreate struct {
	Body      string     `json:"body"`       // 作成するタスクの本文
	TimeLimit *time.Time `json:"time_limit"` // 作成するタスクの期限
}

// OnlyID タスクのIDの情報を表す構造体
type OnlyID struct {
	ID int `json:"id"` // 指定するタスクのID
}

// TaskToUpdate は、タスクの更新の情報を表す構造体
type TaskToUpdate struct {
	ID        int        `json:"id"`         // 更新するタスクのID
	Body      string     `json:"body"`       // 更新後のタスクの本文
	TimeLimit *time.Time `json:"time_limit"` // 変更後タスクの期限
}

// tasks は、削除されていないtaskのデータをjsonに変換して返してる
func tasks(w http.ResponseWriter, r *http.Request) {
	// dbの接続
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	// エラーが発生した場合、ログにエラーメッセージを表示する
	if err != nil {
		log.Fatal(err)
	}

	// クエリを元にdbからデータを取得する
	rows, err := db.Query("SELECT id, body, time_limit, completed_at, created_at FROM TASK WHERE deleted_at IS NULL")
	if err != nil {
		log.Fatal(err)
	}

	// タスク構造体の配列の定義
	var tasks []Task

	// レコード単位で繰り返す
	for rows.Next() {
		// レコードに入っているタスクの定義
		var task Task
		// taskにレコードの各情報を入れる
		err := rows.Scan(&task.ID, &task.Body, &task.TimeLimit, &task.CompletedAt, &task.CreatedAt)
		// エラーが発生した場合、ログにエラーメッセージを表示する
		if err != nil {
			log.Fatal(err)
		}
		// tasksにtaskを追加
		tasks = append(tasks, task)
	}

	// レスポンスヘッダーにContent-Type: application/jsonを追加する
	w.Header().Set("Content-Type", "application/json")
	// tasksをjsonに変換する
	b, _ := json.Marshal(tasks)
	// jsonに変換されたtasksを返してる
	w.Write(b)
}

// create は、タスクを新規作成し、ログにIDと作成日時を表示する間数
func create(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外はエラーを返す
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}
	// TaskToCreate は新規タスクを作成するためのデータ構造
	var task TaskToCreate

	// リクエストのボディからJSONデータを読み取り、TaskToCreateオブジェクトにデコード
	// エラーが発生した場合は、400 Bad Requestステータスコードとエラーメッセージを返す
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// id は新規IDを作成するためのデータ構造
	var id int
	// createdAt は新規作成日を作成するためのデータ構造
	var createdAt time.Time

	// dbの接続
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

	// taskテーブルにBODY,TIME_LIMITの登録しIDとCREATED_ATを返す
	query := `INSERT INTO TASK (BODY, TIME_LIMIT)
	VALUES ($1, $2)
	RETURNING ID, CREATED_AT`

	// クエリの実行
	err = db.QueryRow(query, task.Body, task.TimeLimit).Scan(&id, &createdAt)
	if err != nil {
		log.Fatalf("クエリエラー: %v", err)
	}
	//ログにIDの表示
	log.Println("ID:", id)
	//ログに作成日時を表示
	log.Println("作成日時:", createdAt)
}

// complete　指定したIDのタスクを完了する間数
func complete(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外はエラーを返す
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}
	// IDの指定
	var oi OnlyID

	// リクエストのボディからJSONデータを読み取り、OnlyIDオブジェクトにデコード
	// エラーが発生した場合は、400 Bad Requestステータスコードとエラーメッセージを返す
	if err := json.NewDecoder(r.Body).Decode(&oi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 完了させるタスクのID
	id := oi.ID

	// dbの接続
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	defer db.Close()
	// エラーが発生した場合、ログにエラーメッセージを表示する
	if err != nil {
		log.Fatal(err)
	}

	// 指定されたIDのcompleted_atの更新
	db.QueryRow("UPDATE task SET completed_at = $1 WHERE id = $2", time.Now(), id)

}

// update は、指定されたIDのタスクのbodyとtime_limitを更新する間数
func update(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外はエラーを返す
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}

	// Task は、指定されたIDのbody, time_limitを更新するためのデータ構造
	var task TaskToUpdate

	// リクエストのボディからJSONデータを読み取り、TaskToUpdateオブジェクトにデコード
	// エラーが発生した場合は、400 Bad Requestステータスコードとエラーメッセージを返す
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// dbの接続
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	defer db.Close()
	// エラーが発生した場合、ログにエラーメッセージを表示する
	if err != nil {
		log.Fatal(err)
	}
	// 更新させるタスクのID
	id := task.ID

	// 指定されたIDのbody, time_limitの更新
	db.QueryRow("UPDATE task SET body = $1, time_limit = $2 WHERE id = $3", task.Body, task.TimeLimit, id)
}

// _delete は、指定されたIDのタスクのdeleted_atを更新して論理削除をする間数
func _delete(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外はエラーを返す
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}
	//　IDの指定
	var oi OnlyID

	// リクエストのボディからJSONデータを読み取り、OnlyIDオブジェクトにデコード
	// エラーが発生した場合は、400 Bad Requestステータスコードとエラーメッセージを返す
	if err := json.NewDecoder(r.Body).Decode(&oi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// dbの接続
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	defer db.Close()
	// エラーが発生した場合、ログにエラーメッセージを表示する
	if err != nil {
		log.Fatal(err)
	}

	// 削除するタスクのID
	id := oi.ID

	// 指定されたIDのdeleted_atを更新して論理削除を行う
	db.QueryRow("UPDATE task SET deleted_at = $1 WHERE id = $2", time.Now(), id)

}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Mux を作成
	mux := http.NewServeMux()

	// PATHの指定
	mux.HandleFunc("/api/tasks", tasks)
	mux.HandleFunc("/api/tasks/create", create)
	mux.HandleFunc("/api/tasks/update", update)
	mux.HandleFunc("/api/tasks/complete", complete)
	mux.HandleFunc("/api/tasks/delete", _delete)

	// CORS ミドルウェアを適用
	handler := CORSMiddleware(mux)

	// PORT番号の指定
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", handler)
}
