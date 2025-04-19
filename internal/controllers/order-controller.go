package controllers

import (
	"context"
	"fmt"
	"graves/internal/repository"
	"graves/pkg/dto"
	"graves/pkg/payos"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	payosapi "github.com/payOSHQ/payos-lib-golang"
)

func CreatePaymentLink(c *gin.Context) {
	body := payosapi.CheckoutRequestType{}
	c.ShouldBind(&body)

	for _, item := range body.Items {
		body.Amount += item.Price * item.Quantity
	}

	p := payos.GetInstance()
	repo, err := repository.GetInstance()
	if err != nil {
		c.Error(err)
		return
	}

	userId, err := GetUserId(c)
	if err != nil {
		c.Error(err)
		return
	}

	data, err := p.CreatePaymentLink(body)
	if err != nil {
		c.Error(err)
		return
	}

	if err := repo.CreateOrder(c, userId, *data); err != nil {
		c.Error(err)
		return
	}

	go func(ctx context.Context) {
		time.Sleep(time.Second*time.Duration(p.GetTimeOut()) + 5)
		detailData, err := p.GetPaymentLinkInfo(fmt.Sprintf("%d", data.OrderCode))
		if err != nil {
			return
		}
		if err := repo.UpdateOrderStatus(ctx, int32(detailData.OrderCode), detailData.Status); err != nil {
			return
		}
	}(context.Background())

	c.JSON(http.StatusOK, data)
}

func GetPaymentLinkInfo(c *gin.Context) {
	orderIdStr := c.Query("orderId")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.Error(err)
		return
	}

	p := payos.GetInstance()

	repo, err := repository.GetInstance()
	if err != nil {
		c.Error(err)
		return
	}

	userId, err := GetUserId(c)
	if err != nil {
		c.Error(err)
		return
	}

	if _, err = repo.GetOrderById(c, userId, int32(orderId)); err != nil {
		c.Error(err)
		return
	}

	data, err := p.GetPaymentLinkInfo(orderIdStr)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func CancelPaymentLink(c *gin.Context) {
	orderIdStr := c.Query("orderId")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.Error(err)
		return
	}

	p := payos.GetInstance()

	repo, err := repository.GetInstance()
	if err != nil {
		c.Error(err)
		return
	}

	userId, err := GetUserId(c)
	if err != nil {
		c.Error(err)
		return
	}

	if _, err = repo.GetOrderById(c, userId, int32(orderId)); err != nil {
		c.Error(err)
		return
	}

	data, err := p.CancelPaymentLink(orderIdStr)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetUserId(c *gin.Context) (uint64, error) {
	userId := c.GetHeader("x-user-id")
	if userId == "" {
		return 0, fmt.Errorf("x-user-id header is required")
	}
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("x-user-id header must be a valid uint64")
	}
	return id, nil
}

func ListOrders(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		c.Error(err)
		return
	}

	repo, err := repository.GetInstance()
	if err != nil {
		c.Error(err)
		return
	}

	req := dto.ListOrders{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	data, err := repo.ListOrders(c, userId, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}
