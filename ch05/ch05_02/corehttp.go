// もっとも基本的で最小のhttpサーバー
package main

import (
	"log"
	"net/http"
)

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Go!"))
}

func main() {
	// httpのパスと関数を関連付ける
	http.HandleFunc("/", helloGoHandler)

	// ポートにバインドし、デフォルトでDefaultServeMuxハンドラーを使用する。
	log.Fatal(http.ListenAndServe(":8080", nil))
}
