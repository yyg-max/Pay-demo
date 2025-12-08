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

package model

import (
	"time"

	"gorm.io/gorm"
)

type MerchantAPIKey struct {
	ID             uint64         `json:"id" gorm:"primaryKey"`
	UserID         uint64         `json:"user_id" gorm:"not null;index:idx_merchant_api_keys_user_created,priority:1"`
	ClientID       string         `json:"client_id" gorm:"size:64;uniqueIndex;index:idx_client_credentials,priority:2;not null"`
	ClientSecret   string         `json:"client_secret" gorm:"size:64;index:idx_client_credentials,priority:1;not null"`
	AppName        string         `json:"app_name" gorm:"size:20;not null"`
	AppHomepageURL string         `json:"app_homepage_url" gorm:"size:100;not null"`
	AppDescription string         `json:"app_description" gorm:"size:100"`
	RedirectURI    string         `json:"redirect_uri" gorm:"size:100"`
	NotifyURL      string         `json:"notify_url" gorm:"size:100;not null"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime;index:idx_merchant_api_keys_user_created,priority:2"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
