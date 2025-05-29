/*
 * Copyright 2024 Matthew A. Titmus
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ch04

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Circuit は、コンテキストを受け取り、文字列とエラーを返す関数型です。
// 非同期処理やリトライ、スロットリングなどの機構で使用されます。
type Circuit func(context.Context) (string, error)

// Breaker は、指定された回数の失敗後、指定された時間後に再試行する機能を持つラッパーを返します。
// threshold は、失敗回数の閾値を指定します。
// コンテキストの終了シグナルを監視し、キャンセル時に即時終了します。
// Effector が成功した場合、すぐに結果を返し、失敗した場合はリトライを続けます。
func Breaker(circuit Circuit, threshold int) Circuit {
	var failures int
	var last = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock() // Establish a "read lock"

		d := failures - threshold

		if d >= 0 {
			shouldRetryAt := last.Add((2 << d) * time.Second)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock() // Release read lock

		response, err := circuit(ctx) // Issue the request proper

		m.Lock() // Lock around shared resources
		defer m.Unlock()

		last = time.Now() // Record time of attempt

		if err != nil { // Circuit returned an error,
			failures++           // so we count the failure
			return response, err // and return
		}

		failures = 0 // Reset failures counter

		return response, nil
	}
}
