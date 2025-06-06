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
	"fmt"
	"sync"
	"time"
)

// Effector はコンテキストを受け取り、文字列とエラーを返す関数型です。
// 非同期処理やリトライ、スロットリングなどの機構で使用されます。
type Effector func(context.Context) (string, error)

// Throttle は Effector を指定された最大数とリフリー数と間隔で制限する機能を持つラッパーを返します。
// max は最大実行回数を指定し、refill はリフリー回数を指定します。
// d はリフリー間隔を指定します。
// コンテキストの終了シグナルを監視し、キャンセル時に即時終了します。
// Effector が成功した場合、すぐに結果を返し、失敗した場合はリフリーを続けます。
func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	var tokens = max
	var once sync.Once
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return

					case <-ticker.C:
						m.Lock()
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
						m.Unlock()
					}
				}
			}()
		})

		m.Lock()
		defer m.Unlock()

		if tokens <= 0 {
			return "", fmt.Errorf("too many calls")
		}

		tokens--

		return e(ctx)
	}
}
