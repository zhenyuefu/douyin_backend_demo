package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"gorm.io/gorm"
)

type FollowModel struct {
	gorm.Model
	UID    uint `gorm:"index"`
	FID    uint `gorm:"index"`
	Cancel bool `gorm:"default:false"`
}

func (f *FollowModel) TableName() string {
	return constants.FollowTableName
}
