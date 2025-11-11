package merchant

import (
	"github.com/gin-gonic/gin"
	"github.com/linux-do/pay/internal/model"
)

func GetAPIKeyFromContext(c *gin.Context) (*model.MerchantAPIKey, bool) {
	apiKey, exists := c.Get(APIKeyObjKey)
	if !exists {
		return nil, false
	}
	key, ok := apiKey.(*model.MerchantAPIKey)
	return key, ok
}

func SetAPIKeyToContext(c *gin.Context, apiKey *model.MerchantAPIKey) {
	c.Set(APIKeyObjKey, apiKey)
}
