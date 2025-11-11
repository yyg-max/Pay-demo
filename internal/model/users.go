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
	"gorm.io/gorm"
	"time"
)

type OAuthUserInfo struct {
	Id         uint64     `json:"id"`
	Username   string     `json:"username"`
	Name       string     `json:"name"`
	Active     bool       `json:"active"`
	AvatarUrl  string     `json:"avatar_url"`
	TrustLevel TrustLevel `json:"trust_level"`
}

type User struct {
	ID               uint64     `json:"id" gorm:"primaryKey"`
	Username         string     `json:"username" gorm:"size:64;unique;index"`
	Nickname         string     `json:"nickname" gorm:"size:100"`
	AvatarUrl        string     `json:"avatar_url" gorm:"size:100"`
	TrustLevel       TrustLevel `json:"trust_level" gorm:"index"`
	PayScore         int64      `json:"pay_score" gorm:"default:0;index"`
	PayKey           string     `json:"pay_key" gorm:"size:10;index"`
	SignKey          string     `json:"sign_key" gorm:"size:64;unique;index;not null"`
	TotalBalance     int64      `json:"total_balance" gorm:"default:0;index"`
	AvailableBalance int64      `json:"available_balance" gorm:"default:0;index"`
	IsActive         bool       `json:"is_active" gorm:"default:true"`
	IsAdmin          bool       `json:"is_admin" gorm:"default:false"`
	LastLoginAt      time.Time  `json:"last_login_at" gorm:"index"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime;index"`
}

func (u *User) Exact(tx *gorm.DB, id uint64) error {
	if err := tx.Where("id = ?", id).First(u).Error; err != nil {
		return err
	}
	return nil
}
