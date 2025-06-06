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
	"fmt"
	"testing"
	"time"
)

var count int

// EmulateTransientError は一時的な失敗をシミュレーションする関数です。
// 最初の3回はエラーを返し、それ以降は成功を返します。
func EmulateTransientError(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "intentional fail", errors.New("error")
	} else {
		return "success", nil
	}
}

// TestRetry はエフェクター関数のリトライ動作をテストするための関数です。
// リトライ回数と遅延間隔を設定し、コンテキストを通じて実際の挙動を確認します。
func TestRetry(t *testing.T) {
	ctx := context.Background()
	r := Retry(EmulateTransientError, 5, 2*time.Second)
	res, err := r(ctx)

	fmt.Println(res, err)
}
