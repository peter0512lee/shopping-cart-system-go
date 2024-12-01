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

type ProductHandler struct {
	collection *mongo.Collection
}

func NewProductHandler(collection *mongo.Collection) *ProductHandler {
	return &ProductHandler{
		collection: collection,
	}
}

// CreateProduct 創建新商品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.collection.InsertOne(context.Background(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, product)
}

// GetProducts 獲取所有商品
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var products []models.Product

	cursor, err := h.collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) CreateBulkProducts(c *gin.Context) {
	var products []models.Product
	if err := c.ShouldBindJSON(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	var insertProducts []interface{}
	for _, product := range products {
		product.CreatedAt = now
		product.UpdatedAt = now
		insertProducts = append(insertProducts, product)
	}

	result, err := h.collection.InsertMany(context.Background(), insertProducts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating products"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Products created successfully",
		"count":   len(result.InsertedIDs),
	})
}

func (h *ProductHandler) ClearProducts(c *gin.Context) {
	// 刪除集合中的所有文檔
	err := h.collection.Drop(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error clearing products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All products have been cleared"})
}
