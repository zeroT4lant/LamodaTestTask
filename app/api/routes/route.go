package routes

import (
	"LamodaTestTask/app/api/controller"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "LamodaTestTask/app/docs"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

//в маршрутах "route" - обрабатываем http запросы

// ErrorResponse структура ошибки ответа
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewRouter инициализация роутов
func NewRouter(db *sql.DB) *gin.Engine {
	// Создаём роутер gin
	r := gin.Default()

	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler()))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Обработчик для создания нового склада
	r.POST("/create-warehouse", func(c *gin.Context) {
		// Считываем данные склада из тела запроса
		wh := controller.Warehouse{}
		err := c.BindJSON(&wh)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse data"})
			return
		}

		// Создаем новый склад в базе данных
		err = controller.CreateWarehouse(db, &wh)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Отправляем ответ с ID нового склада
		c.JSON(http.StatusCreated, gin.H{"id": wh.ID})
	})

	// Обработчик для создания нового продукта на заданном складе
	r.POST("/create-product", func(c *gin.Context) {
		var p controller.Product
		err := c.BindJSON(&p)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
			return
		}

		err = controller.CreateProduct(db, &p)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": p.ID})
	})

	// Удаление продукта
	r.DELETE("/delete-product/:id", func(c *gin.Context) {
		//на всякий случай, id может быть у нас строкой
		//из запроса парсим - id
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		if err := controller.DeleteProduct(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	// Резервирование товаров
	r.POST("/reserve-products", func(c *gin.Context) {
		//привязываем массив уникальных кодов товара
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid request body",
			})
			return
		}

		//сохраняем по айди
		err := controller.ReserveProducts(db, productCodes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	})

	// Отмена резервирования продуктов / освобождение резерва товаров
	r.POST("/release-products", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid request body",
			})
			return
		}

		err := controller.ReleaseProducts(db, productCodes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	})

	// Получения оставшегося количества продуктов на складе
	r.GET("/remaining-products/:warehouseID", func(c *gin.Context) {
		warehouseID := c.Param("warehouseID")
		var id int
		if _, err := fmt.Sscan(warehouseID, &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid warehouse ID",
			})
			return
		}

		products, err := controller.GetRemainingProducts(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, products)
	})

	return r
}
