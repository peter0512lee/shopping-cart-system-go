package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/shopping-cart/internal/models" // 替換為你的專案名
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartHandler struct {
	cartCollection    *mongo.Collection
	productCollection *mongo.Collection
}

func NewCartHandler(cartCollection, productCollection *mongo.Collection) *CartHandler {
	return &CartHandler{
		cartCollection:    cartCollection,
		productCollection: productCollection,
	}
}

// AddToCart 添加商品到購物車
func (h *CartHandler) AddToCart(c *gin.Context) {
	var input struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 驗證並獲取商品信息
	productID, err := primitive.ObjectIDFromHex(input.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	err = h.productCollection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 獲取或創建購物車
	var cart models.Cart
	err = h.cartCollection.FindOne(context.Background(), bson.M{"user_id": input.UserID}).Decode(&cart)

	if err == mongo.ErrNoDocuments {
		// 創建新購物車
		cart = models.Cart{
			UserID:    input.UserID,
			Items:     make([]models.CartItem, 0),
			Total:     0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	// 更新購物車項目
	foundIndex := -1
	for i, item := range cart.Items {
		if item.ProductID == productID {
			foundIndex = i
			break
		}
	}

	if foundIndex >= 0 {
		// 更新現有項目
		cart.Items[foundIndex].Quantity += input.Quantity
	} else {
		// 添加新項目
		cart.Items = append(cart.Items, models.CartItem{
			ProductID: productID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  input.Quantity,
		})
	}

	// 計算總金額
	var total float64
	for _, item := range cart.Items {
		total += item.Price * float64(item.Quantity)
	}
	cart.Total = total
	cart.UpdatedAt = time.Now()

	// 保存到數據庫
	if cart.ID.IsZero() {
		result, err := h.cartCollection.InsertOne(context.Background(), cart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
		cart.ID = result.InsertedID.(primitive.ObjectID)
	} else {
		_, err = h.cartCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": cart.ID},
			bson.M{"$set": cart},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
			return
		}
	}

	c.JSON(http.StatusOK, cart)
}

// GetCart 獲取購物車內容
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var cart models.Cart
	err := h.cartCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// UpdateCartItem 更新購物車商品數量
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.Param("user_id")
	var input struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity cannot be negative"})
		return
	}

	productID, err := primitive.ObjectIDFromHex(input.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var cart models.Cart
	err = h.cartCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// 更新商品數量
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			if input.Quantity == 0 {
				// 如果數量為0，移除該商品
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i].Quantity = input.Quantity
			}
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found in cart"})
		return
	}

	// 重新計算總金額
	var total float64
	for _, item := range cart.Items {
		total += item.Price * float64(item.Quantity)
	}
	cart.Total = total
	cart.UpdatedAt = time.Now()

	// 更新購物車
	_, err = h.cartCollection.UpdateOne(
		context.Background(),
		bson.M{"user_id": userID},
		bson.M{"$set": cart},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// RemoveFromCart 從購物車移除商品
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.Param("user_id")
	productID := c.Param("product_id")

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var cart models.Cart
	err = h.cartCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// 移除商品
	found := false
	for i, item := range cart.Items {
		if item.ProductID == objID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found in cart"})
		return
	}

	// 重新計算總金額
	var total float64
	for _, item := range cart.Items {
		total += item.Price * float64(item.Quantity)
	}
	cart.Total = total
	cart.UpdatedAt = time.Now()

	// 更新購物車
	_, err = h.cartCollection.UpdateOne(
		context.Background(),
		bson.M{"user_id": userID},
		bson.M{"$set": cart},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// ClearCart 清空購物車
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.Param("user_id")

	result, err := h.cartCollection.UpdateOne(
		context.Background(),
		bson.M{"user_id": userID},
		bson.M{
			"$set": bson.M{
				"items":      []models.CartItem{},
				"total":      0,
				"updated_at": time.Now(),
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
