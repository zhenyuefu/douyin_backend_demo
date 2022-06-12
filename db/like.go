package db

type LikeModel struct {
	UID uint `gorm:"primaryKey;column:uid"`
	VID uint `gorm:"primaryKey;column:vid"`
}

func (l *LikeModel) TableName() string {
	return "like"
}

// Like 创建点赞
func Like(like *LikeModel) error {
	return DB.Create(like).Error
}

func UnLike(uid uint, vid uint) error {
	return DB.Where("uid = ? and vid = ?", uid, vid).Delete(&LikeModel{}).Error
}
