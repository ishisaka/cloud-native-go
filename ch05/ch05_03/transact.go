package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sync"
)

// EventType は、トランザクションログのイベントの種類を表します。
type EventType byte

const (
	_                     = iota // iota == 0; ignore the zero value
	EventDelete EventType = iota // iota == 1
	EventPut                     // iota == 2; implicitly repeat
)

// Event は、トランザクションログのイベントを表します。
type Event struct {
	Sequence  uint64    // A unique record ID
	EventType EventType // The action taken
	Key       string    // The key affected by this transaction
	Value     string    // The value of a PUT the transaction
}

// TransactionLogger は、トランザクションログを記録するためのインターフェースです。
type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	Close() error
	Wait()

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

// FileTransactionLogger は、ファイルベースのトランザクションロガーを表します。
type FileTransactionLogger struct {
	events       chan<- Event // Write-only channel for sending events
	errors       <-chan error // Read-only channel for receiving errors
	lastSequence uint64       // The last used event sequence number
	file         *os.File     // The location of the transaction log
	wg           *sync.WaitGroup
}

// WritePut は、トランザクションログにPUTイベントを書き込みます。
func (l *FileTransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

// WriteDelete は、トランザクションログに削除イベントを書き込みます。
func (l *FileTransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventDelete, Key: key}
}

// Err は、エラーチャネルを返します。
func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

// Run は、トランザクションログのイベントを処理するgoroutineを開始します。
func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() {
		for e := range events { // Retrieve the next Event

			l.lastSequence++ // Increment sequence number

			_, err := fmt.Fprintf( // Write the event to the log
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
				return
			}

			l.wg.Done()
		}
	}()
}

// ReadEvents は、トランザクションログからイベントを読み取ります。
func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			fmt.Sscanf(
				line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value)

			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			uv, err := url.QueryUnescape(e.Value)
			if err != nil {
				outError <- fmt.Errorf("value decoding failure: %w", err)
				return
			}

			e.Value = uv
			l.lastSequence = e.Sequence

			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

// Wait は、トランザクションログの処理が完了するまで待機します。
func (l *FileTransactionLogger) Wait() {
	l.wg.Wait()
}

// Close は、トランザクションログを閉じます。
func (l *FileTransactionLogger) Close() error {
	l.Wait()

	if l.events != nil {
		close(l.events) // Terminates Run loop and goroutine
	}

	return l.file.Close()
}

// NewFileTransactionLogger は、新しいファイルベースのトランザクションロガーを作成します。
// filenameはトランザクションログファイルのパスを指定します。
// ファイルのオープンに失敗した場合、エラーを返します。
func NewFileTransactionLogger(filename string) (*FileTransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: file, wg: &sync.WaitGroup{}}, nil
}
