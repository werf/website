package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var dbConn *sql.DB

func main() {
	// Get database filename.
	dbFilePath := os.Getenv("SQLITE_FILE")
	if dbFilePath == "" {
		log.Println("Environment variable SQLITE_FILE is not set! I will use in-memory database.")
		dbFilePath = ":memory:"
	}

	// Open database file.
	var err error
	dbConn, err = sql.Open("sqlite3", dbFilePath+"?_busy_timeout=99999")
	if err != nil {
		log.Fatalf("Could not connect to SQLite database: %v", err)
	}
	log.Printf("Connected to SQLite database.")

	// Close db on SIGTERM.
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		dbConn.Close()
		os.Exit(1)
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/labels", func(r chi.Router) {
		r.Get("/", listLabels)
		r.Post("/", addLabel)
		r.Route("/{labelId}", func(r chi.Router) {
			r.Get("/", getLabel)
			r.Post("/", modifyLabel)
			r.Delete("/", deleteLabel)
		})
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:3000", r))
}

type Label struct {
	Id    int64  `json:"id"`
	Label string `json:"label"`
}

type AppError struct {
	Result  string `json:"result"`
	Comment string `json:"comment"`
}

// handleAppError return error as JSON
func handleAppError(w http.ResponseWriter, err error) {
	log.Printf("ERROR %v", err)
	appError := AppError{
		Result:  "error",
		Comment: err.Error(),
	}
	errBytes, _ := json.Marshal(appError)
	w.Write(errBytes)
}

// handleError returns code 500
func handleError500(w http.ResponseWriter, err error) {
	log.Printf("ERROR %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

const listSql = `
SELECT * FROM labels
`

func listLabels(w http.ResponseWriter, r *http.Request) {
	rows, err := dbConn.Query(listSql)
	if err != nil {
		handleAppError(w, err)
		return
	}
	defer rows.Close()
	labels := make([]Label, 0)
	for rows.Next() {
		var label = Label{}
		err = rows.Scan(
			&label.Id,
			&label.Label,
		)
		if err != nil {
			handleAppError(w, err)
			return
		}
		labels = append(labels, label)
	}
	outBytes, _ := json.Marshal(labels)
	w.Write(outBytes)
}

const insertSql = `
INSERT INTO labels (label) VALUES (?)
`

func addLabel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handleAppError(w, err)
		return
	}

	labelText := r.Form.Get("label")

	result, err := dbConn.Exec(insertSql, labelText)

	if err != nil {
		handleAppError(w, err)
		return
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		handleAppError(w, err)
		return
	}
	labelJson := fmt.Sprintf(`{"id":%d, "label":"%s"}`,
		lastId,
		labelText,
	)
	w.Write([]byte(labelJson))
}

const getSql = `
SELECT * 
FROM labels 
WHERE id = ?
`

func getLabel(w http.ResponseWriter, r *http.Request) {
	labelId := chi.URLParam(r, "labelId")

	rows, err := dbConn.Query(getSql, labelId)
	if err != nil {
		handleAppError(w, err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		handleAppError(w, fmt.Errorf("entry not found"))
		return
	}

	var label = Label{}
	err = rows.Scan(
		&label.Id,
		&label.Label,
	)
	if err != nil {
		handleAppError(w, err)
		return
	}

	outBytes, _ := json.Marshal(label)
	w.Write(outBytes)
}

const updateSql = `
UPDATE labels
SET label=?
WHERE id=?
`

func modifyLabel(w http.ResponseWriter, r *http.Request) {
	labelId := chi.URLParam(r, "labelId")

	err := r.ParseForm()
	if err != nil {
		handleAppError(w, err)
		return
	}

	labelText := r.Form.Get("label")

	_, err = dbConn.Exec(updateSql, labelText, labelId)
	if err != nil {
		handleAppError(w, err)
		return
	}

	labelJson := fmt.Sprintf(`{"id":%s, "label":"%s"}`,
		labelId,
		labelText,
	)
	w.Write([]byte(labelJson))
}

const deleteSql = `DELETE FROM labels WHERE id=?`

func deleteLabel(w http.ResponseWriter, r *http.Request) {
	labelId := chi.URLParam(r, "labelId")

	_, err := dbConn.Exec(deleteSql, labelId)
	if err != nil {
		handleAppError(w, err)
		return
	}

	w.Write([]byte(`{"result":true}`))
}
