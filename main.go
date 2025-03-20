package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	ID          int        `json:"id"`
	Body        string     `json:"body"`
	TimeLimit   *time.Time `json:"time_limit"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
type TaskToCreate struct {
	Body      string     `json:"body"`
	TimeLimit *time.Time `json:"time_limit"`
}
type OnlyID struct {
	ID int `json:"id"`
}
type TaskToUpdate struct {
	ID        int        `json:"id"`
	Body      string     `json:"body"`
	TimeLimit *time.Time `json:"time_limit"`
}

func tasks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, body, time_limit, completed_at, created_at FROM TASK WHERE deleted_at IS NULL")

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Body, &task.TimeLimit, &task.CompletedAt, &task.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(tasks)
	w.Write(b)
}

func create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}
	var task TaskToCreate
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var id int
	var createdAt time.Time
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	query := `INSERT INTO TASK (BODY, TIME_LIMIT)
	VALUES ($1, $2)
	RETURNING ID, CREATED_AT`
	err = db.QueryRow(query, task.Body, task.TimeLimit).Scan(&id, &createdAt)
	if err != nil {
		log.Fatalf("クエリエラー: %v", err)
	}
	log.Println("ID:", id)
	log.Println("作成日時:", createdAt)
}

func complete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}
	var oi OnlyID
	if err := json.NewDecoder(r.Body).Decode(&oi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 完了させるタスクのID
	id := oi.ID

	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("UPDATE task SET completed_at = $1 WHERE id = $2", time.Now(), id)

}

func update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}
	var task TaskToUpdate
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	id := task.ID

	db.QueryRow("UPDATE task SET body = $1, time_limit = $2 WHERE id = $3", task.Body, task.TimeLimit, id)
}

func _delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}
	var oi OnlyID
	if err := json.NewDecoder(r.Body).Decode(&oi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 削除するタスクのID
	id := oi.ID

	db.QueryRow("UPDATE task SET deleted_at = $1 WHERE id = $2", time.Now(), id)

}

func main() {
	http.HandleFunc("/api/tasks", tasks)
	http.HandleFunc("/api/tasks/create", create)
	http.HandleFunc("/api/tasks/update", update)
	http.HandleFunc("/api/tasks/complete", complete)
	http.HandleFunc("/api/tasks/delete", _delete)
	http.ListenAndServe(":8080", nil)
}
