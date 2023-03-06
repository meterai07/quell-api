package initializers

import (
	"quell-api/entity"
)

func SyncDatabase() {
	DB.AutoMigrate(
		entity.User{},
		entity.Category{},
		entity.Post{},
		entity.Attachment{},
		entity.Saving{},
		entity.SavingCategory{},
	)
}
