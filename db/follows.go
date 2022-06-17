package db

type FollowModel struct {
	uid uint `gorm:"primaryKey;column:uid;index"`
	fid uint `gorm:"primaryKey;column:fid;index"`
}

func (f *FollowModel) TableName() string {
	return "follows"
}
