/*
 * MIT License
 *
 * Copyright (c) 2025 linux.do
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/google/uuid"
)

// GenerateUniqueIDSimple 生成 64 位唯一标识符
func GenerateUniqueIDSimple() string {
	randomBytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		// 如果随机数生成失败，使用 UUID 作为后备
		uuidBytes := []byte(uuid.NewString())
		hash := sha256.Sum256(uuidBytes)
		copy(randomBytes, hash[:])
	}
	return hex.EncodeToString(randomBytes)
}
