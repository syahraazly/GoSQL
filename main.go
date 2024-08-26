package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	user     = "root"
	password = "password"
	dbname   = "mini_challenge"
)

var db *sql.DB

type Product struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Variant struct {
	ID          int
	VariantName string
	Quantity    int
	ProductID   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductWithVariant struct {
	ID        int
	Name      string
	Variants  []Variant
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	mysqlInfo := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	db, err = sql.Open("mysql", mysqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database")

	// CreateProduct("Iphone15")
	// UpdateProduct(6, "Macbook Pro M1")
	// GetProductByID(6)
	// CreateVariant("Pink", 10, 6)
	UpdateVariantByID(7, "Pink", 15, 7)
	// GetProductWithVariant(6)
	// DeleteVariantByID(9)
}

func CreateProduct(name string) {
	query := `INSERT INTO products (name, created_at, updated_at) VALUES (?, ?, ?)`

	_, err := db.Exec(query, name, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}

	fmt.Println("New product created successfully")
}

func UpdateProduct(id int, newName string) {
	query := `UPDATE products SET name = ?, updated_at = ? WHERE id = ?`

	_, err := db.Exec(query, newName, time.Now(), id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Product updated successfully")
}

func GetProductByID(productID int) {
    var product Product
    var createdAtStr, updatedAtStr string

    sqlStatement := `SELECT id, name, created_at, updated_at FROM products WHERE id = ?`
    row := db.QueryRow(sqlStatement, productID)

    err := row.Scan(&product.ID, &product.Name, &createdAtStr, &updatedAtStr)
    if err != nil {
        panic(err)
    }

    // parsing strings to time.Time
    product.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
    product.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

    fmt.Printf("Product: %+v\n", product)
}


func CreateVariant(variantName string, quantity, productID int) {
	query := `INSERT INTO variants (variant_name, quantity, product_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(query, variantName, quantity, productID, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}

	fmt.Println("New variant created successfully")
}

func UpdateVariantByID(id int, variantName string, quantity, productID int) {
	sqlStatement := `UPDATE variants SET variant_name = ?, quantity = ?, product_id = ?, updated_at = ? WHERE id = ?`

	_, err := db.Exec(sqlStatement, variantName, quantity, productID, time.Now(), id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Variant updated successfully")
}

func DeleteVariantByID(id int) {
	sqlStatement := `DELETE FROM variants WHERE id = ?`

	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Variant deleted successfully")
}

func GetProductWithVariant(productID int) {
	var product ProductWithVariant

	query := `
	SELECT p.id, p.name, p.created_at, p.updated_at,
	       v.id, v.variant_name, v.quantity, v.product_id, v.created_at, v.updated_at
	FROM products p
	LEFT JOIN variants v ON p.id = v.product_id
	WHERE p.id = ?`

	rows, err := db.Query(query, productID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var variant Variant
		var createdAt, updatedAt, variantCreatedAt, variantUpdatedAt string

		err := rows.Scan(&product.ID, &product.Name, &createdAt, &updatedAt,
			&variant.ID, &variant.VariantName, &variant.Quantity, &variant.ProductID, &variantCreatedAt, &variantUpdatedAt)
		if err != nil {
			panic(err)
		}

		// parsing strings to time.Time
		product.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		product.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		variant.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", variantCreatedAt)
		variant.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", variantUpdatedAt)

		// append variant to the product's list of variants
		product.Variants = append(product.Variants, variant)
	}

	fmt.Printf("Product with Variants: %+v\n", product)
}
