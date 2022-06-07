package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"gorm.io/gorm"
)

type LikeModel struct {
	gorm.Model
	UID    uint `gorm:"index"`
	VID    uint `gorm:"index"`
	Cancel bool `gorm:"default:false"`
}

func (l *LikeModel) TableName() string {
	return constants.LikeTableName
}

// CreateFavorite
func CreateLike(like *LikeModel) error {
	return DB.Create(like).Error
}

// GetFavoriteList
func GetLikeList(uid uint) ([]LikeModel, error) {
	var likes []LikeModel
	result := DB.Where("user_id = ?", uid).Find(&likes)
	return likes, result.Error
}
