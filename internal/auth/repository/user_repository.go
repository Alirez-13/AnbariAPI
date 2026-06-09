// # SINGLE REASON: Preserve legacy user repository constructor.
package repository

import (
	"AnbariAPI/internal/auth/infrastructure"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return infrastructure.NewGormUserRepository(db)
}
