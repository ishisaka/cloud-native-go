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
	"sync"
	"testing"
	"time"
)

// TestDebounceFirstDataRace tests for data races.
func TestDebounceFirstDataRace(t *testing.T) {
	ctx := context.Background()

	circuit := failAfter(1)
	debounce := DebounceFirst(circuit, time.Second)

	wg := sync.WaitGroup{}

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := debounce(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	time.Sleep(time.Second * 2)
	// 1秒以上経過しているのでcircuitが開いているので以下の処理はエラーになる
	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := debounce(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	wg.Wait()
}
