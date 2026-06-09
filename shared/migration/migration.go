package migration

import (
	models3 "AnbariAPI/internal/auth/domain"
	"AnbariAPI/internal/catalog/domain"
	domain2 "AnbariAPI/internal/inventory/domain"
	"log"

	"gorm.io/gorm"
)

// Migrate runs the database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Auto-migrate schema
	err := db.AutoMigrate(
		&domain.Category{},
		&domain.Product{},
		&domain2.Transaction{},
		&domain2.TransactionLine{},
		&domain2.InventoryBatch{},
		&domain2.BatchAllocation{},
		&models3.User{},    // ADDED User
		&models3.Session{}, // ADDED Session
	)

	if err != nil {
		return err
	}

	// Create indexes manually if needed (GORM tags handle most basic indexes)
	indexes := []struct {
		Table   string
		Name    string
		Columns string
	}{
		{"products", "idx_product_category", "category_id"},
		{"products", "idx_product_name", "name"},
		{"transactions", "idx_transaction_type", "type"},
		{"transactions", "idx_transaction_date", "date"},
		{"transaction_lines", "idx_line_transaction", "transaction_id"},
		{"transaction_lines", "idx_line_product", "product_id"},
	}

	for _, idx := range indexes {
		query := "CREATE INDEX IF NOT EXISTS " + idx.Name + " ON " + idx.Table + "(" + idx.Columns + ")"
		if err := db.Exec(query).Error; err != nil {
			log.Printf("Failed to create index %s: %v", idx.Name, err)
			return err
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
