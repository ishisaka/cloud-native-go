/* log/slogのサンプル */
package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	// Infoレベル(0)のログを出力
	slog.Info("Hello, world!")
	// プリセットの無いログレベル2のログを出力
	slog.Log(context.TODO(), 2, "Hello, world!")
	// ログメッセージに属性を含める
	slog.Info("Hello", "number", 3)
	// slog.Attrを使用する
	slog.Info("hello", slog.Int("number", 3))
	// Logger.Withメソッドを使って新しいLoggerを構築し、全てのレコードにその属性を含める
	logger := slog.Default()
	logger2 := logger.With("url", "https://example.com")
	logger2.Info("Hello, world!") // 2025/06/03 15:16:13 INFO Hello, world! url=https://example.com
	// JSONハンドラを使ってJSON形式で出力する
	loggerwJsonHandler := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	loggerwJsonHandler.Info("Hello, world!") // {"time":"2025-06-03T15:18:52.53884+09:00","level":"INFO","msg":"Hello, world!"}
	loggerwJsonHandler.Info("Hello", slog.Int("number", 3))
}
