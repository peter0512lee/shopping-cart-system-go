package main

import (
	"log"

	"github.com/yourusername/shopping-cart/api/handlers" // 替換為你的專案名
	"github.com/yourusername/shopping-cart/pkg/database" // 替換為你的專案名

	"github.com/gin-gonic/gin"
)

func main() {
	// 連接資料庫
	client := database.ConnectDB()

	// 初始化集合
	productCollection := database.GetCollection(client, "products")
	cartCollection := database.GetCollection(client, "carts")

	// 初始化處理器
	productHandler := handlers.NewProductHandler(productCollection)
	cartHandler := handlers.NewCartHandler(cartCollection, productCollection)

	// 設置路由
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 基本的健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由
	api := r.Group("/api/v1")
	{
		// 商品路由
		api.POST("/products", productHandler.CreateProduct)
		api.GET("/products", productHandler.GetProducts)
		api.DELETE("/products", productHandler.ClearProducts)
		api.POST("/products/bulk", productHandler.CreateBulkProducts)

		// 購物車路由
		api.POST("/cart", cartHandler.AddToCart)
		api.GET("/cart/:user_id", cartHandler.GetCart)
		api.PUT("/cart/:user_id", cartHandler.UpdateCartItem)                // 更新商品數量
		api.DELETE("/cart/:user_id/:product_id", cartHandler.RemoveFromCart) // 移除商品
		api.DELETE("/cart/:user_id", cartHandler.ClearCart)                  // 清空購物車

	}

	// 啟動服務器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
