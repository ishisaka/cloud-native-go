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
)

// TimeoutFunction は、文字列を受け取り、文字列とエラーを返す関数型です。
type TimeoutFunction func(string) (string, error)

// WithContext は、コンテキストと文字列を受け取り、文字列とエラーを返す関数型です。
type WithContext func(context.Context, string) (string, error)

// Timeout は、コンテキストを受け取り、文字列とエラーを返す関数型です。
// 非同期処理やリトライ、スロットリングなどの機構で使用されます。
// 最後の呼び出しか処理しないようにします。
func Timeout(f TimeoutFunction) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		ch := make(chan struct {
			result string
			err    error
		}, 1)

		go func() {
			defer close(ch)

			res, err := f(arg)
			select {
			case ch <- struct {
				result string
				err    error
			}{res, err}:
			case <-ctx.Done():
				// コンテキストがキャンセルされた場合は早期リターン
				return
			}
		}()

		select {
		case res := <-ch:
			return res.result, res.err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
