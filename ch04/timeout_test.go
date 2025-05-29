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
	"testing"
	"time"
)

// TestTimeoutNo は、コンテキストのタイムアウトがヒットしない場合の動作を確認するテストです。
// コンテキストのタイムアウトは 2 秒に設定され、`Timeout` 関数が正常に動作するか検証します。
// 実行中に予期しないエラーが発生した場合はテストが失敗します。
func TestTimeoutNo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	timeout := Timeout(Slow)
	_, err := timeout(ctx, "some input")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// TestTimeoutYes は、コンテキストのタイムアウトがヒットする場合の動作を確認するテストです。
// コンテキストのタイムアウトは 1 秒に設定され、`Timeout` 関数が正常に動作するか検証します。
// 実行中に予期しないエラーが発生した場合はテストが失敗します。
func TestTimeoutYes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/2)
	defer cancel()

	timeout := Timeout(Slow)
	_, err := timeout(ctx, "some input")
	if err == nil {
		t.Fatal("Didn't get expected timeout error")
	}
	if err != context.DeadlineExceeded {
		t.Fatalf("Unexpected error: %v", err)
	}
	fmt.Println("Got expected timeout error")
}

// Slow は、1秒間の待機を含む関数です。
func Slow(s string) (string, error) {
	time.Sleep(time.Second)
	fmt.Println("Slow")
	return "Got input: " + s, nil
}
