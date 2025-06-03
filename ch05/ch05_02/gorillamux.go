package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func main() {
	// gorilla/muxのルーターを作成する
	r := mux.NewRouter()

	// ルートにハンドラーを登録する
	r.HandleFunc("/", helloMuxHandler)

	// ポートにバインドし、gorilla/mux ルーターを使用する。
	log.Fatal(http.ListenAndServe(":8080", r))
}
