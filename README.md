# Industrial-Programming-Technologies-Practice-14 (ЭФМО-02-24 Кучер Артем Сергеевич)
## Swagger + Docker (Интернет-магазина электронники)
### Swagger UI

![image](https://github.com/user-attachments/assets/91fa5eb3-8942-4e47-ae52-0db2dfa79dcd)

![image](https://github.com/user-attachments/assets/b6cd13c7-0959-4746-8fda-f847df1c86ec)

#### func login
```
// @Summary Login
// @Description User login
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body Credentials true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	var validUser *Credentials
	for _, user := range users {
		if user.Username == creds.Username && user.Password == creds.Password {
			validUser = &user
			break
		}
	}

	if validUser == nil {
		handleError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	token, err := generateToken(creds.Username)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Could not create token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
```
#### func refreshToken
```
// @Summary Refresh Token
// @Description Refresh JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /refresh [post]
func refreshToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				newToken, err := generateToken(claims.Username)
				if err != nil {
					handleError(c, http.StatusInternalServerError, "could not refresh token")
					return
				}
				c.JSON(http.StatusOK, gin.H{"token": newToken})
				return
			}
		}
		handleError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if !token.Valid {
		handleError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	handleError(c, http.StatusBadRequest, "token is still valid")
}
```
#### func main
```
// @title My API
// @version 1.0
// @description This is a sample API for an electronic store.
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
func main() {
	initDB()
	router.POST("/login", login)
	router.POST("/refresh", refreshToken)
	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		// Получение всех товаров
		router.GET("/products", getProducts)
		// Получение товара по ID
		router.GET("/products/:id", getProductByID)
		// Создание нового товара
		router.POST("/products", createProduct)
		// Обновление существующего товара
		router.PUT("/products/:id", updateProduct)
		// Удаление товара
		router.DELETE("/products/:id", deleteProduct)

		// Получение всех товаров в корзине
		protected.GET("/basket", getBasket)

		// Добавление товара в корзину
		protected.POST("/basket", addToBasket)

		// Удаление товара из корзины
		protected.DELETE("/basket/:productId", deleteFromBasket)

		protected.GET("/productswithtimeout", getProductsWithTimeout)
	}
	router.POST("/tasks", func(c *gin.Context) {
		taskID := tasks.CreateTask()
		go tasks.RunTask(taskID)
		c.JSON(201, gin.H{"task_id": taskID})
	})

	// Получение статуса задачи
	router.GET("/tasks/:id", func(c *gin.Context) {
		taskID := c.Param("id")
		task := tasks.GetTask(taskID)
		if task == nil {
			c.JSON(404, gin.H{"error": "Task not found"})
			return
		}

		c.JSON(200, task)
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
```
#### func getBasket
```
// @Summary Get Basket
// @Description Get a list of products in the user's basket
// @Tags basket
// @Accept json
// @Produce json
// @Success 200 {object} []BasketItem
// @Router /basket [get]
func getBasket(c *gin.Context) {
	c.JSON(http.StatusOK, basket)
}
```
#### func addToBasket
```
// @Summary Add To Basket
// @Description Add a new product to the user's basket
// @Tags basket
// @Accept json
// @Produce json
// @Param item body BasketItem true "Basket item"
// @Success 201 {object} BasketItem
// @Failure 400 {object} map[string]string
// @Router /basket [post]
func addToBasket(c *gin.Context) {
	var newItem BasketItem

	if err := c.BindJSON(&newItem); err != nil {
		handleError(c, http.StatusBadRequest, "invalid request")
		return
	}

	for i, item := range basket {
		if item.ProductId == newItem.ProductId {
			basket[i].Quantity += newItem.Quantity
			c.JSON(http.StatusOK, basket[i])
			return
		}
	}

	basket = append(basket, newItem)
	c.JSON(http.StatusCreated, newItem)
}
```
#### func deleteFromBasket
```
// @Summary Delete From Basket
// @Description Remove a product from the user's basket
// @Tags basket
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /basket/{productId} [delete]
func deleteFromBasket(c *gin.Context) {
	productId := c.Param("productId")

	for i, item := range basket {
		if item.ProductId == productId {
			basket = append(basket[:i], basket[i+1:]...)
			handleError(c, http.StatusOK, "product removed from basket")
			return
		}
	}
	handleError(c, http.StatusNotFound, "product not found in basket")

}
```
#### func getProducts
```
// @Summary Get Products
// @Description Get a list of products with pagination, sorting, and filtering options
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(5)
// @Param name query string false "Product name"
// @Param category query string false "Product category"
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func getProducts(c *gin.Context) {
	var products []Product
	var total int64
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "5")
	name := c.Query("name")
	category := c.Query("category")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	offset := (pageInt - 1) * limitInt
	query := db.Limit(limitInt).Offset(offset)

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if category != "" {
		query = query.Where("category ILIKE ?", "%"+category+"%")
	}

	query.Find(&products).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
	sort := c.Query("sort")
	if sort != "" {
		query = query.Order(sort)
	}
}
```
#### func getProductByID
```
// @Summary Get Product by ID
// @Description Get a single product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func getProductByID(c *gin.Context) {
	id := c.Param("id")
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		handleError(c, http.StatusNotFound, "Product not found")
		return
	}
	c.JSON(http.StatusOK, product)
}
```
#### func createProduct
```
// @Summary Create Product
// @Description Add a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "New product"
// @Success 201 {object} Product
// @Failure 400 {object} map[string]string
// @Router /products [post]
func createProduct(c *gin.Context) {
	var newProduct Product
	if err := c.BindJSON(&newProduct); err != nil {
		handleError(c, http.StatusNotFound, "Product not found")
		return
	}
	db.Create(&newProduct)
	c.JSON(http.StatusCreated, newProduct)
}
```
#### func updateProduct
```
// @Summary Update Product
// @Description Update an existing product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body Product true "Updated product"
// @Success 200 {object} Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct Product
	if err := c.BindJSON(&updatedProduct); err != nil {
		handleError(c, http.StatusNotFound, "Product not found")
		return
	}
	if err := db.Model(&Product{}).Where("id = ?", id).Updates(updatedProduct).Error; err != nil {
		handleError(c, http.StatusNotFound, "Product not found")
		return
	}
	c.JSON(http.StatusOK, updatedProduct)
}
```
#### func deleteProduct
```
// @Summary Delete Product
// @Description Remove a product from the catalog
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Product{}, id).Error; err != nil {
		handleError(c, http.StatusNotFound, "Product not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
```
#### func getProductsWithTimeout
```
// @Summary Get Products With Timeout
// @Description Get a list of products with a request timeout
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} []Product
// @Router /productswithtimeout [get]
func getProductsWithTimeout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var products []Product
	if err := db.WithContext(ctx).Find(&products).Error; err != nil {
		handleError(c, http.StatusRequestTimeout, "Request timed out")
		return
	}

	c.JSON(http.StatusOK, products)
}
```

### Docker
#### Dockerfile
```
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:latest  

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
```

#### docker-compose.yml
```
version: '3'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
    environment:
      - DATABASE_URL=postgres://postgres:пароль@db:5432/pr10?sslmode=disable

  db:
    image: postgres:13
    environment:
      - POSTGRES_PASSWORD=пароль
      - POSTGRES_DB=pr10
    volumes:
      - ./data:/var/lib/postgresql/data
```


![image](https://github.com/user-attachments/assets/4c0a09b0-9564-44bd-8a0f-2bdf3212935c)
