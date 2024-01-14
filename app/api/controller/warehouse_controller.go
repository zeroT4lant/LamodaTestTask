package controller

//	@title			Warehouse API Documentation
//	@description	This is a sample API for a warehouse application
//	@version		1
//	@host			localhost:8080

import (
	"database/sql"
	"errors"
)

// Product структура продукта
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Size        string `json:"size"`
	Code        string `json:"code"`
	Quantity    int    `json:"quantity"`
	WarehouseID int    `json:"warehouse_id"`
}

// Warehouse структура склада
type Warehouse struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	IsAvailable bool   `json:"is_available" db:"is_available"`
}

//	@Summary		Create a new warehouse.
//	@Description	Create a new warehouse in the database.
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			warehouse	body		Warehouse		true	"Warehouse information"
//	@Success		200			{string}	string			"Warehouse created"
//	@Failure		400			{object}	ErrorResponse	"Invalid request format"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/create-warehouse [post]
//
// CreateWarehouse создает новый склад и записывает в базу
func CreateWarehouse(db *sql.DB, w *Warehouse) error {
	// Подготовка запроса для вставки нового склада
	stmt, err := db.Prepare("INSERT INTO warehouse(name, is_available) VALUES($1, $2) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Вставка нового склада и получение его идентификатора
	//здесь Scan, с помощью которого можно считать все полученные данные в переменные
	err = stmt.QueryRow(w.Name, w.IsAvailable).Scan(&w.ID)
	if err != nil {
		return err
	}

	return nil
}

//	@Summary		Create a new product.
//	@Description	Create a new product on a specified warehouse.
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			product	body		Product			true	"Product information"
//	@Success		200		{string}	string			"Product created"
//	@Failure		400		{object}	ErrorResponse	"Invalid request format"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/create-product [post]
//
// CreateProduct создает новый продукт на заданном складе
func CreateProduct(db *sql.DB, p *Product) error {
	// Подготовка запроса для вставки нового продукта
	stmt, err := db.Prepare("INSERT INTO products(name, size, code, quantity, warehouse_id) VALUES($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Вставка нового продукта и получение его идентификатора
	err = stmt.QueryRow(p.Name, p.Size, p.Code, p.Quantity, p.WarehouseID).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

//	@Summary		Delete a product
//	@Description	Delete a product by its ID.
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Product ID"
//	@Success		200	{string}	string			"Product deleted successfully"
//	@Failure		400	{object}	ErrorResponse	"Invalid request format"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/delete-product/:id [delete]
//
// DeleteProduct удаляет продукт по ID
func DeleteProduct(db *sql.DB, id int) error {
	// Удаляем продукт из базы данных
	_, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

//	@Summary		Reserves products
//	@Description	Reserves products and updates their quantities
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			productCodes	query		[]string	true	"Product codes"
//	@Success		204				{string}	string		""
//	@Failure		400				{object}	ErrorResponse
//	@Failure		500				{object}	ErrorResponse
//	@Router			/reserve-products [post]
//
// ReserveProducts резервирует продукты
func ReserveProducts(db *sql.DB, productCodes []string) error {
	if len(productCodes) == 0 {
		return errors.New("empty product codes")
	}

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Зарезервируем каждый продукт в цикле
	for _, code := range productCodes {
		row := tx.QueryRow("SELECT id, name, size, code, quantity FROM products WHERE code = $1 FOR UPDATE", code)

		//записываем продукт
		var p Product
		err := row.Scan(&p.ID, &p.Name, &p.Size, &p.Code, &p.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Проверяем, доступен ли продукт для бронирования
		if p.Quantity < 1 {
			tx.Rollback()
			return errors.New("product is out of stock")
		}

		// Обновляем количество продукта
		_, err = tx.Exec("UPDATE products SET quantity = quantity - 1 WHERE id = $1", p.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//	@Summary		Releases products
//	@Description	Releases reserved products and updates their quantities
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			productCodes	query		[]string	true	"Product codes"
//	@Success		204				{string}	string		""
//	@Failure		400				{object}	ErrorResponse
//	@Failure		500				{object}	ErrorResponse
//	@Router			/release-products [post]
//
// ReleaseProducts реализует товары
func ReleaseProducts(db *sql.DB, productCodes []string) error {
	// Проверяем массив на пустоту массива кодов
	if len(productCodes) == 0 {
		return errors.New("empty product codes")
	}

	// Начинаем новую транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Создаем новую транзакцию
	updateStmt, err := tx.Prepare("UPDATE products SET quantity = quantity + 1 WHERE code = $1")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer updateStmt.Close()

	// Проходимся по каждому продукту
	for _, code := range productCodes {
		// Проверяем существует ли продукт
		var p Product
		err := db.QueryRow("SELECT id, name, size, code, quantity FROM products WHERE code = $1", code).Scan(&p.ID, &p.Name, &p.Size, &p.Code, &p.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Обновляем количество продукта
		_, err = updateStmt.Exec(p.Code)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// @Description Get remaining products for a given warehouse.
// @Tags products
// @Accept json
// @Produce json
// @Param warehouseID path int true "Warehouse ID"
// @Success 200 {array} Product "Remaining products"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /remaining-products/{warehouseID} [get]
// GetRemainingProducts возвращает оставшееся количество продуктов на складе
func GetRemainingProducts(db *sql.DB, warehouseID int) ([]Product, error) {
	// Проходимся по строкам, возвращенным запросом, и добавляем каждую строку к слайсу продуктов.
	rows, err := db.Query("SELECT id, name, size, code, quantity FROM products WHERE warehouse_id = $1", warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Создаем пустой слайс для хранения результатов
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Size, &p.Code, &p.Quantity); err != nil {
			return nil, err
		}
		p.WarehouseID = warehouseID
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
