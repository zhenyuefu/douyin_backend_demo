package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"gorm.io/gorm"
)

type CommentModel struct {
	gorm.Model
	UID     uint
	VID     uint `gorm:"index"`
	Content string
}

func (c *CommentModel) TableName() string {
	return constants.CommentTableName
}

// CreateComment 创建评论
func CreateComment(comment *CommentModel) error {
	return DB.Create(comment).Error
}

// DeleteComment 删除评论
func DeleteComment(id uint, uid uint) error {
	return DB.Where("uid=?", uid).Delete(&CommentModel{}, id).Error
}

// GetComment 获取评论
func GetComment(commentList *[]CommentModel, vid uint) error {
	res := DB.Model(&CommentModel{}).Where("v_id = ?", vid).Find(&commentList)
	return res.Error
}
