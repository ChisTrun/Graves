package routers

import (
	"graves/internal/controllers"
	"graves/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.ErrorHandlerMiddleware())

	protected := router.Group("/graves")
	{
		protected.POST("/create-payment-link", controllers.CreatePaymentLink)
		protected.GET("/payment-link-info", controllers.GetPaymentLinkInfo)
		protected.POST("/cancel-payment-link", controllers.CancelPaymentLink)
		protected.POST("/list-orders", controllers.ListOrders)
		protected.POST("/verify-payment-webhook", controllers.VerifyPaymentWebhookData)
	}

	return router
}
