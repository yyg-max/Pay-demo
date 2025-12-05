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
	"context"
	"fmt"
	"time"

	"github.com/linux-do/pay/internal/db"
	"github.com/linux-do/pay/internal/logger"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

type OrderType string

const (
	OrderTypeReceive   OrderType = "receive"
	OrderTypePayment   OrderType = "payment"
	OrderTypeTransfer  OrderType = "transfer"
	OrderTypeCommunity OrderType = "community"
)

type OrderStatus string

const (
	OrderStatusSuccess   OrderStatus = "success"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusExpired   OrderStatus = "expired"
	OrderStatusDisputing OrderStatus = "disputing"
	OrderStatusRefund    OrderStatus = "refund"
	OrderStatusRefused   OrderStatus = "refused"
)

type Order struct {
	ID              uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo         string          `json:"order_no" gorm:"-"`
	OrderName       string          `json:"order_name" gorm:"size:64;not null"`
	MerchantOrderNo string          `json:"merchant_order_no" gorm:"size:64;index"`
	ClientID        string          `json:"client_id" gorm:"size:64;index:idx_orders_client_status_created,priority:1;index:idx_orders_client_payee,priority:1;index:idx_orders_client_payer,priority:1"`
	PayerUserID     uint64          `json:"payer_user_id" gorm:"index:idx_orders_payer_status_type_created,priority:1;index:idx_orders_payer_status_type_trade,priority:1;index:idx_orders_client_payer,priority:2"`
	PayeeUserID     uint64          `json:"payee_user_id" gorm:"index:idx_orders_payee_status_type_created,priority:1;index:idx_orders_client_payee,priority:2"`
	PayerUsername   string          `json:"payer_username" gorm:"->"`
	PayeeUsername   string          `json:"payee_username" gorm:"->"`
	Amount          decimal.Decimal `json:"amount" gorm:"type:numeric(20,2);not null;index"`
	Status          OrderStatus     `json:"status" gorm:"type:varchar(20);not null;index:idx_orders_payee_status_type_created,priority:2;index:idx_orders_payer_status_type_created,priority:2;index:idx_orders_client_status_created,priority:2;index:idx_orders_payer_status_type_trade,priority:2"`
	Type            OrderType       `json:"type" gorm:"type:varchar(20);not null;index:idx_orders_payee_status_type_created,priority:3;index:idx_orders_payer_status_type_created,priority:3;index:idx_orders_payer_status_type_trade,priority:3"`
	Remark          string          `json:"remark" gorm:"size:255"`
	PaymentType     string          `json:"payment_type" gorm:"size:20"`
	TradeTime       time.Time       `json:"trade_time" gorm:"index:idx_orders_payer_status_type_trade,priority:4"`
	ExpiresAt       time.Time       `json:"expires_at" gorm:"not null"`
	CreatedAt       time.Time       `json:"created_at" gorm:"autoCreateTime;index:idx_orders_payee_status_type_created,priority:4;index:idx_orders_payer_status_type_created,priority:4;index:idx_orders_client_status_created,priority:3"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// AfterFind 格式化 OrderNo
func (o *Order) AfterFind(*gorm.DB) error {
	o.OrderNo = fmt.Sprintf("%018d", o.ID)
	return nil
}

// ExpirePendingOrders 将已过期且 pending 状态的订单设置为 expired
func ExpirePendingOrders(ctx context.Context) {
	result := db.DB(ctx).Model(&Order{}).
		Where("status = ? AND expires_at <= ?", OrderStatusPending, time.Now()).
		Update("status", OrderStatusExpired)

	if result.Error != nil {
		logger.ErrorF(ctx, "过期 pending 订单失败: %v", result.Error)
	} else {
		logger.InfoF(ctx, "已将 %d 个已过期的 pending 订单设置为 expired", result.RowsAffected)
	}
}
