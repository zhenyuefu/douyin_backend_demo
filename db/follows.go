package db

type FollowModel struct {
	uid uint `gorm:"primaryKey;column:uid"`
	fid uint `gorm:"primaryKey;column:fid"`
}

func (f *FollowModel) TableName() string {
	return "follows"
}
