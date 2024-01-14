package controller

import (
	"LamodaTestTask/app/randomizer"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

// work
func TestCreateWarehouse(t *testing.T) {
	// Коннектимся к базе
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=lamoda_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создам новый склад
	w := &Warehouse{
		Name:        randomizer.RandomString(6),
		IsAvailable: true,
	}

	// Вызываем функцию создания нового склада
	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем успешно ли был создан склад
	if w.ID == 0 {
		t.Errorf("Expected warehouse ID to be non-zero, got %d", w.ID)
	}
}

// work
func TestCreateProduct(t *testing.T) {
	// Подключаемся к базе
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=lamoda_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем новый склад (снова т.к мы не знаем айдишник и название склада заренее)
	w := &Warehouse{
		Name:        randomizer.RandomString(6),
		IsAvailable: true,
	}

	// Вызываем функцию для создания склада
	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем продукт
	p := &Product{
		Name:        randomizer.RandomString(6),
		Size:        randomizer.RandomString(6),
		Code:        randomizer.RandomString(6),
		Quantity:    randomizer.RandomInt(6),
		WarehouseID: w.ID,
	}

	// Вызываем функцию создания продукта
	err = CreateProduct(db, p)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем успешно ли создан продукт
	if p.ID == 0 {
		t.Errorf("Expected product ID to be non-zero, got %d", p.ID)
	}
}

// work
func TestReserveProductsEmptyProductCodes(t *testing.T) {
	// Подключаемся к базе
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=lamoda_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = ReserveProducts(db, []string{})
	if err == nil {
		t.Error("Expected an error with empty product codes, but got nil")
	}
}

// work
func TestReserveProductsInvalidProductCode(t *testing.T) {
	// Подключаемся к базе
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=lamoda_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем новый склад
	w := &Warehouse{
		Name:        randomizer.RandomString(6),
		IsAvailable: true,
	}
	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем новый продукт и добавляем его на склад
	p := &Product{
		Name:        randomizer.RandomString(6),
		Size:        randomizer.RandomString(6),
		Code:        randomizer.RandomString(6),
		Quantity:    1,
		WarehouseID: w.ID,
	}
	err = CreateProduct(db, p)
	if err != nil {
		t.Fatal(err)
	}

	// Пытаемся зарезервировать продукт с неверным кодом
	err = ReserveProducts(db, []string{"invalid-code"})
	if err == nil {
		t.Error("Expected an error with invalid product code, but got nil")
	}
}

// work
func TestReserveProductsProductOutOfStock(t *testing.T) {
	// Подключаемся к базе
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=lamoda_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем новый склад
	w := &Warehouse{
		Name:        randomizer.RandomString(6),
		IsAvailable: true,
	}
	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем новый продукт и добавляем его на склад
	p := &Product{
		Name:        randomizer.RandomString(6),
		Size:        randomizer.RandomString(6),
		Code:        randomizer.RandomString(6),
		Quantity:    0, // устанавливаем количество 0, чтобы продукт был недоступен для бронирования
		WarehouseID: w.ID,
	}
	err = CreateProduct(db, p)
	if err != nil {
		t.Fatal(err)
	}

	// Пытаемся зарезервировать продукт, который отсутствует на складе
	err = ReserveProducts(db, []string{p.Code})
	if err == nil {
		t.Errorf("Expected error, but got nil")
	} else if err.Error() != "product is out of stock" {
		t.Errorf("Expected error message 'product is out of stock', but got '%s'", err.Error())
	}
}
