package initializers

import (
	"quell-api/entity"
)

func SyncDatabase() {
	DB.AutoMigrate(
		entity.User{},
		entity.Category{},
		entity.Post{},
		entity.UserPremium{},
		entity.Attachment{},
		entity.UserTransaction{},
		entity.SavingCategory{},
		entity.Saving{},
	)
}
