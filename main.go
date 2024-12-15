package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
	_ "m1/docs"
	"m1/tasks"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Address struct {
	gorm.Model
	ID         uint
	CustomerID uint
	Country    string
	City       string
	Street     string
	House      int
}

type Customer struct {
	gorm.Model
	ID         uint
	Name       string
	Surname    string
	Patronymic string
	Telephone  string
	CardNumber string
	Mail       string
	Password   string
	Addresses  []Address
}

type Product struct {
	Id                string
	ManufacturerId    int
	ProductCategoryId int
	Name              string
	Price             float64
	Quantity          int
	ArticleNumber     int
	Description       string
	ImageUrl          string
}

type Product_Order struct {
	Id        string
	ProductId int
	OrderId   int
}

type Category struct {
	gorm.Model
	ID   uint
	Name string
}

type Manufacturer struct {
	gorm.Model
	ID      uint
	Name    string
	Country string
}

type Order struct {
	gorm.Model
	ID         uint
	Date       time.Time
	Time       time.Time
	CustomerID uint
	Deliveries []Delivery
	Products   []Product
}

type Review struct {
	gorm.Model
	ID         uint
	CustomerID uint
	ProductID  uint
	Date       time.Time
	Text       string
	Rating     int
}

type Favorite struct {
	gorm.Model
	ID         uint
	ProductID  uint
	CustomerID uint
}

type Provider struct {
	gorm.Model
	ManufacturerID uint
	ProductID      uint
}

type Basket struct {
	gorm.Model
	ID         uint
	CustomerID uint
	Products   []Product
}

type Delivery struct {
	gorm.Model
	ID           uint
	OrderID      uint
	AddressID    uint
	Status       string
	SendDate     time.Time
	ExpectedDate time.Time
}

/*
var products = []Product{
	{Id: "1", ManufacturerId: 1, ProductCategoryId: 1, Name: "Samsung Galaxy S24", Price: 78000, Quantity: 10, ArticleNumber: 8100, Description: "Samsung Galaxy S24 — компактный смартфон из флагманской линейки, получивший мощный фирменный процессор, большой объем оперативной памяти, узнаваемый дизайн и продуманную эргономику. ", ImageUrl: "https://ir.ozone.ru/s3/multimedia-m/c1000/6900636406.jpg"},
	{Id: "2", ManufacturerId: 1, ProductCategoryId: 1, Name: "Samsung Galaxy S23 FE", Price: 36000, Quantity: 50, ArticleNumber: 7100, Description: "Смартфон Samsung Galaxy S23 FE - это новейшее устройство от известного южнокорейского производителя Samsung. Он предлагает пользователям высокую производительность и отличные функциональные возможности.", ImageUrl: "https://img.championat.com/i/d/o/168794395937683540.jpg"},
	{Id: "3", ManufacturerId: 2, ProductCategoryId: 2, Name: "Xiaomi Pad 6", Price: 26500, Quantity: 12, ArticleNumber: 5300, Description: "Xiaomi Pad 6: Стиль и производительность в одном устройстве Стильный дизайн, который впечатляет: Тонкий металлический корпус (6.51 мм) с элегантной фактурой", ImageUrl: "https://cdn1.ozone.ru/s3/multimedia-c/6718970892.jpg"},
	{Id: "4", ManufacturerId: 2, ProductCategoryId: 2, Name: "Xiaomi Pad 6 Pro", Price: 49000, Quantity: 5, ArticleNumber: 1000, Description: "Планшет Xiaomi Pad 6S Pro 8/256 GB - это устройство, которое сочетает в себе высокую производительность и стильный дизайн.", ImageUrl: "https://www.gizmochina.com/wp-content/uploads/2023/06/Xiaomi-Pad-6-Featured-scaled.jpg"},
	{Id: "5", ManufacturerId: 3, ProductCategoryId: 3, Name: "Asus TUF Gaming A15", Price: 85000, Quantity: 10, ArticleNumber: 1100, Description: "Игровой ноутбук ASUS TUF Gaming A15 FA507NU-LP141 обеспечивает впечатляющую производительность для игр и других задач.", ImageUrl: "https://avatars.mds.yandex.net/i?id=a03c8b88bfe7858a4bfea3daa18d8fcf_l-6998621-images-thumbs&n=13"},
	{Id: "6", ManufacturerId: 3, ProductCategoryId: 3, Name: "Asus ROG Zephyrus G15", Price: 15000, Quantity: 5, ArticleNumber: 19000, Description: "Мощный и портативный, ноутбук ROG Zephyrus G15 представляет собой игровую платформу на базе операционной системы Windows 10, выполненную в ультратонком корпусе весом всего 1,9 кг.", ImageUrl: "https://static.onlinetrade.ru/img/items/m/noutbuk_asus_rog_zephyrus_m15_gu502lv_az105t_15.6_fhd_ag_ips_240hz_i7_10750h_16gb_1024gb_ssd_nodvd_rtx_2060_6gb_w10_black_1459559_3.jpg"},
	{Id: "7", ManufacturerId: 3, ProductCategoryId: 3, Name: "Asus Vivobook 17", Price: 41500, Quantity: 12, ArticleNumber: 900, Description: "ASUS Vivobook 17 X1704ZA-AU341 90NB10F2-M00DD0 — производительный ноутбук в ударопрочном корпусе с активной системой охлаждения.", ImageUrl: "https://avatars.mds.yandex.net/i?id=b40b973a43129287ac2cfdb6ff688283_l-6458590-images-thumbs&n=13"},
	{Id: "8", ManufacturerId: 4, ProductCategoryId: 4, Name: "Sony WH-1000XM4", Price: 25000, Quantity: 25, ArticleNumber: 3400, Description: "Беспроводные наушники Sony WH-1000XM4 — это флагманские наушники с активным шумоподавлением и высоким качеством звука.", ImageUrl: "https://ir.ozone.ru/s3/multimedia-2/c1000/6766801130.jpg"},
	{Id: "9", ManufacturerId: 4, ProductCategoryId: 4, Name: "Sony WH-CH720N", Price: 8000, Quantity: 50, ArticleNumber: 9000, Description: "Благодаря технологии шумоподавления, легкой конструкции и длительному времени работы от аккумулятора вы сможете наслаждаться музыкой дольше и без отвлекающих окружающих звуков.", ImageUrl: "https://cdn1.ozone.ru/s3/multimedia-x/c600/6670694193.jpg"},
	{Id: "10", ManufacturerId: 5, ProductCategoryId: 3, Name: "Apple MacBook Air 13 Retina", Price: 65000, Quantity: 10, ArticleNumber: 7900, Description: "Самый тонкий и лёгкий ноутбук Apple MacBook Air 13 model: A2337 теперь стал суперсильным благодаря чипу Apple M1.", ImageUrl: "https://msk.aura-rent.ru/wp-content/uploads/2020/02/noutbook-air-2.2.jpg"},
	{Id: "11", ManufacturerId: 5, ProductCategoryId: 3, Name: "Apple MacBook Pro 14 2023", Price: 135000, Quantity: 7, ArticleNumber: 14100, Description: "Ноутбук Apple Macbook Pro 14 M3 8/512 Silver (MR7J3) – мощный и стильный помощник для работы и развлечений.", ImageUrl: "https://mtscdn.ru/upload/iblock/ee8/mbp14_silver_gallery7_202301.jpg"},
	{Id: "12", ManufacturerId: 5, ProductCategoryId: 5, Name: "Apple iMac 24", Price: 190000, Quantity: 10, ArticleNumber: 410, Description: "Моноблок Apple iMac 24 2023 года - это мощный и стильный компьютер, который станет незаменимым помощником в вашей повседневной работе.", ImageUrl: "https://avatars.mds.yandex.net/get-mpic/3732535/2a0000018a6d36516ed5f119259f018168b1/orig"},
	{Id: "13", ManufacturerId: 6, ProductCategoryId: 6, Name: "Nintendo Switch OLED", Price: 25500, Quantity: 20, ArticleNumber: 16700, Description: "Консоль Nintendo Switch OLED с красочным 7-дюймовым экраном. При практически одинаковых размерах с Nintendo Switch консоль Nintendo Switch OLED отличается более крупным 7-дюймовым OLED-экраном с глубокими цветами и высоким контрастом.", ImageUrl: "https://cdn1.ozone.ru/s3/multimedia-s/6080757748.jpg"},
	{Id: "14", ManufacturerId: 7, ProductCategoryId: 6, Name: "Sony PlayStation 5 Slim", Price: 47000, Quantity: 16, ArticleNumber: 4500, Description: "Игровая приставка Sony PlayStation 5 Slim: улучшенный дизайн и расширенные возможности хранения данных Обновленная версия популярной консоли PlayStation 5", ImageUrl: "https://appmistore.ru/upload/iblock/c4e/25bnswfucwuectbrw1s9vlgvr4bkswud.webp"},
	{Id: "15", ManufacturerId: 8, ProductCategoryId: 6, Name: "Microsoft Xbox Series S", Price: 34000, Quantity: 25, ArticleNumber: 19300, Description: "Игровая консоль Microsoft Xbox Series S рассчитана на использование игр, загружаемых из цифровой библиотеки. ", ImageUrl: "https://media.wired.co.uk/photos/606d9dbb20fc96acca6d3a5a/1:1/w_2000,h_2000,c_limit/3.jpg"},
	{Id: "16", ManufacturerId: 9, ProductCategoryId: 7, Name: "JBL Charge 5", Price: 13000, Quantity: 40, ArticleNumber: 56000, Description: "Charge 5 - портативная колонка от JBL, предназначенная для использования в любых условиях. Она имеет защиту от воды и пыли по стандарту IP67, что позволяет использовать ее на пляже или в походе.", ImageUrl: "https://kazandigital.ru/uploaded/images/abouts/157610-speakers-review-jbl-charge-5-review-image1-4bvjkgsxy5.jpg"},
	{Id: "17", ManufacturerId: 9, ProductCategoryId: 7, Name: "JBL Flip 6", Price: 9000, Quantity: 40, ArticleNumber: 109000, Description: "Страна-производитель Китай Общие параметры Тип портативная колонка Модель JBL Flip 6 Код производителя [JBLFLIP6BLK] Основной цвет черный Акустические характеристики", ImageUrl: "https://avatars.mds.yandex.net/get-mpic/2017233/2a0000018de355c9d630bec9a88e2b702591/optimize"},
	{Id: "18", ManufacturerId: 10, ProductCategoryId: 8, Name: "NVIDIA GeForce RTX 3060 Dual LHR", Price: 29000, Quantity: 20, ArticleNumber: 67100, Description: "Видеокарта Palit GeForce RTX 3060 Dual 12G поможет тебе получить стабильный FPS выше 60 кадров в секунду при максимальных настройках графики и разрешении Full HD.", ImageUrl: "https://avatars.mds.yandex.net/get-mpic/10352132/2a0000018ec730cbd98af51c3c2abe51c47d/optimize"},
	{Id: "19", ManufacturerId: 11, ProductCategoryId: 8, Name: "AMD Radeon RX 6600 PULSE", Price: 36000, Quantity: 41, ArticleNumber: 41100, Description: "Техпроцесс: 7 нм; Тип видеокарты: игровая; Графический процессор: Radeon RX 6600;", ImageUrl: "https://avatars.mds.yandex.net/get-mpic/4120567/2a000001922a9d36a89c564a9f2f29901969/optimize"},
	{Id: "20", ManufacturerId: 1, ProductCategoryId: 9, Name: "Samsung Galaxy Watch 6", Price: 19000, Quantity: 40, ArticleNumber: 9900, Description: "Смарт-часы Samsung Galaxy Watch6 44 мм Silver (SM-R940) обладают дисплеем диагональю 1,47 дюйма и разрешением 480x480 пикселей — увеличить экран и добавить пространства для свайпов позволила тонкая рамка.", ImageUrl: "https://avatars.mds.yandex.net/get-mpic/10352132/2a0000018c5f0c1ddd3b729147a231ecc7b7/optimize"},
}
*/

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=Kuc1804SX dbname=pr10 port=5432 sslmode=disable search_path=public"
	}
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Миграция схемы
	db.AutoMigrate(&Address{}, &Customer{}, &Product{}, &Order{}, &Review{}, &Favorite{}, &Category{}, &Manufacturer{}, &Provider{}, &Basket{}, &Delivery{})
	if err != nil {
		log.Fatalf("Не удалось выполнить миграцию базы данных: %v", err)
	}
}

var router = gin.Default()

var basket = []BasketItem{}

type BasketItem struct {
	ProductId string `json:"productId"`
	Quantity  int
}

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var users = []Credentials{
	{Username: "user", Password: "password"},
	{Username: "user1", Password: "password1"},
	{Username: "user2", Password: "password2"},
	{Username: "user3", Password: "password3"},
}

func handleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					handleError(c, http.StatusUnauthorized, "token expired")
					c.Abort()
					return
				}
			}
			handleError(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		if !token.Valid {
			handleError(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

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
