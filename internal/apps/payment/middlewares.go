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

package payment

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linux-do/pay/internal/db"
	"github.com/linux-do/pay/internal/model"
	"github.com/linux-do/pay/internal/util"
	"github.com/shopspring/decimal"
)

// CreateOrderRequest 商户创建订单统一请求
type CreateOrderRequest struct {
	OrderName       string          `json:"order_name" binding:"required,max=64"`
	MerchantOrderNo string          `json:"merchant_order_no"`
	Amount          decimal.Decimal `json:"amount" binding:"required"`
	Remark          string          `json:"remark" binding:"max=200"`
	PaymentType     string          `json:"payment_type"`
}

// EPayRequest 易支付请求
type EPayRequest struct {
	ClientID        string          `form:"pid" binding:"required"`
	OrderName       string          `form:"name" binding:"required,max=64"`
	MerchantOrderNo string          `form:"out_trade_no" binding:"required"`
	Amount          decimal.Decimal `form:"money" binding:"required"`
	NotifyURL       string          `form:"notify_url"`
	ReturnURL       string          `form:"return_url"`
	Device          string          `form:"device"`
	Sign            string          `form:"sign" binding:"required"`
	PayType         string          `form:"type" binding:"required"`
	SignType        string          `form:"sign_type"`
}

// ToCreateOrderRequest 转换为通用创建订单请求
func (r *EPayRequest) ToCreateOrderRequest() *CreateOrderRequest {
	return &CreateOrderRequest{
		OrderName:       r.OrderName,
		MerchantOrderNo: r.MerchantOrderNo,
		Amount:          r.Amount,
		PaymentType:     r.PayType,
	}
}

// RequireMerchantAuth 验证商户 ClientID/ClientSecret（Basic Auth）
func RequireMerchantAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization: Basic base64(ClientID:ClientSecret)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Err("缺少认证信息"))
			return
		}

		// 解析 Basic Auth
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Err("认证格式错误"))
			return
		}

		// 解码 base64
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Err("认证信息解码失败"))
			return
		}

		// 解析 ClientID:ClientSecret
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Err("认证信息格式错误"))
			return
		}

		clientID := credentials[0]
		clientSecret := credentials[1]

		var apiKey model.MerchantAPIKey
		if err := db.DB(c.Request.Context()).
			Where("client_secret = ? AND client_id = ?", clientSecret, clientID).
			First(&apiKey).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Err("认证失败"))
			return
		}

		util.SetToContext(c, APIKeyObjKey, &apiKey)

		c.Next()
	}
}

// RequireSignatureAuth 验证签名
func RequireSignatureAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		PayType := c.PostForm("type")

		var apiKey model.MerchantAPIKey

		switch PayType {
		case util.PayTypeEPay:
			if createOrderReq, err := VerifySignature(c, &apiKey); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.Err(err.Error()))
				return
			} else {
				util.SetToContext(c, CreateOrderRequestKey, createOrderReq)
			}
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, util.Err("不支持的请求类型"))
			return
		}

		util.SetToContext(c, APIKeyObjKey, &apiKey)

		c.Next()
	}
}
