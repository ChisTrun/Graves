package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"graves/internal/repository"
	"graves/pkg/logger/pkg/logging"
	"graves/pkg/payos"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	payosapi "github.com/payOSHQ/payos-lib-golang"
)

func VerifyPaymentWebhookData(c *gin.Context) {
	var webhookDataReq payosapi.WebhookType

	body, _ := io.ReadAll(c.Request.Body)

	p := payos.GetInstance()

	err := json.Unmarshal(body, &webhookDataReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webhookData, err := p.VerifyPaymentWebhookData(webhookDataReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(ctx context.Context) {
		detailData, err := p.GetPaymentLinkInfo(fmt.Sprintf("%d", webhookData.OrderCode))
		if err != nil {
			logging.Logger(ctx).Error(fmt.Sprintf("Failed to get payment link info: %v", err))
			return
		}
		repo, err := repository.GetInstance()
		if err != nil {
			logging.Logger(ctx).Error(fmt.Sprintf("Failed to get repository instance: %v", err))
			return
		}
		if err := repo.UpdateOrderStatus(ctx, int32(detailData.OrderCode), detailData.Status); err != nil {
			logging.Logger(ctx).Error(fmt.Sprintf("Failed to update order status: %v", err))
			return
		}
	}(context.Background())
	c.JSON(http.StatusOK, webhookData)
}
