package dao

import (
	"HelpStudent/internal/app/users/model"
	"context"
	"gorm.io/gorm"
)

type users struct {
	*gorm.DB
}

func (u *users) Init(db *gorm.DB) (err error) {
	u.DB = db
	return db.AutoMigrate(&model.Users{})
}

func (u *users) CreateWithBind(ctx context.Context, user *model.Users, bind *model.UserBind) error {
	tx := u.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existedUser := model.Users{}
	result := tx.Model(&model.Users{}).Where("staff_id = ?", user.StaffId).First(&existedUser)
	if result.RowsAffected == 1 {
		if res := tx.Model(&model.UserBind{}).Where("user_id = ?", existedUser.Id).
			Update("union_id", bind.UnionId); res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		tx.Commit()
		return nil
	}

	if result = tx.Create(user); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	bind.UserId = user.Id
	if result = tx.Create(bind); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}
