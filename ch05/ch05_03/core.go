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

package main

import (
	"errors"
	"sync"
)

// store はキーと値のペアを保持するマップです。
// スレッドセーフです。
var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var ErrorNoSuchKey = errors.New("no such key")

// Delete は指定されたキーをstoreから削除します。
// キーが存在しない場合でもエラーは返さずnilを返します。
func Delete(key string) error {
	// 書き込みロックを獲得します。
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}

// Get は指定されたキーに対応する値を返します。
// キーが存在しない場合、ErrorNoSuchKeyエラーを返します。
func Get(key string) (string, error) {
	// 読み取りロックを獲得します。
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

// Put は指定されたキーと値をstoreに保存します。
func Put(key string, value string) error {
	// 書き込みロックを獲得します。
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}
