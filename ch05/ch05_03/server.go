package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var logger TransactionLogger

// initializeTransactionLog は、トランザクションログを初期化します。
func initializeTransactionLog() error {
	var err error

	logger, err = NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors: // Retrieve any errors
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Got a DELETE event!
				err = Delete(e.Key)
			case EventPut: // Got a PUT event!
				err = Put(e.Key, e.Value)
			}
		}
	}

	logger.Run()

	return err
}

// loggingMiddleware は、リクエストをログに記録するミドルウェアです。
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// helloMuxHandler は、/ に対する GET リクエストを処理する。
func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Mux!\n"))
}
func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

// keyValuePutHandler は、/v1/{key} に対する PUT リクエストを処理する。
func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	value, err := io.ReadAll(r.Body) // The request body has our value
	defer r.Body.Close()

	if err != nil { // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value)) // Store the value as a string
	if err != nil {               // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	// トランザクションログに書き込みログを追加
	logger.WritePut(key, string(value))

	w.WriteHeader(http.StatusCreated) // All good! Return StatusCreated
}

// keyValueGetHandler は、/v1/{key} に対する GET リクエストを処理する。
func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	value, err := Get(key) // Get value for key
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value)) // Write the value to the response
}

// keyValueDeleteHandler は、/v1/{key} に対する DELETE リクエストを処理する。
func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// トランザクションログに削除ログを追加
	logger.WriteDelete(key)

	log.Printf("DELETE key=%s\n", key)
}

func main() {
	// トランザクションログを初期化する
	err := initializeTransactionLog()
	if err != nil {
		panic(err)
	}

	// gorilla/muxのルーターを作成する
	r := mux.NewRouter()

	// ミドルウェアを登録する
	r.Use(loggingMiddleware)

	// ルートにハンドラーを登録する
	r.HandleFunc("/", notAllowedHandler)

	// ルートにput用のハンドラーを登録する
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	// ルートにget用のハンドラーを登録する
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	// ルートにdelete用のハンドラーを登録する
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")

	r.HandleFunc("/v1", notAllowedHandler)
	r.HandleFunc("/v1/{key}", notAllowedHandler)

	// ポートにバインドし、gorilla/mux ルーターを使用する。
	log.Fatal(http.ListenAndServe(":8080", r))
}
