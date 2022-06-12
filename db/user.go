package db

import (
	"errors"
	"github.com/RaymondCode/simple-demo/constants"
	"github.com/RaymondCode/simple-demo/structs"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Identifier string `gorm:"uniqueIndex;size:32"`
	Credential []byte
	Name       string       `json:"name"`
	Avatar     string       `json:"avatar"`
	Follows    []UserModel  `gorm:"many2many:follows;joinForeignKey:uid;joinReferences:fid"`
	Followers  []UserModel  `gorm:"many2many:follows;joinForeignKey:fid;joinReferences:uid"`
	Videos     []VideoModel `gorm:"foreignKey:AuthorID"`
	Likes      []VideoModel `gorm:"many2many:like;joinForeignKey:uid;joinReferences:vid"`
}

func (u *UserModel) TableName() string {
	return constants.UserTableName
}

// CreateUser 创建用户
func CreateUser(user *UserModel) error {
	return DB.Create(user).Error
}

// GetUser 获取用户
func GetUser(user *structs.User, uid uint, loginUID uint) error {
	var userModel UserModel
	result := DB.First(&userModel, uid).Scan(user)
	user.FollowerCount = userModel.FollowerCount()
	user.FollowCount = userModel.FollowCount()
	user.IsFollow = userModel.IsFollow(loginUID)
	user.TotalFavorite = userModel.TotalFavorite()
	user.FavoriteCount = userModel.FavoriteCount()
	user.WorkCount = userModel.WorkCount()
	return result.Error
}

func GetUserModel(user *UserModel, uid uint) error {
	return DB.First(&user, uid).Error
}

// VerifyCredential 验证密码
func VerifyCredential(identifier string, credential string) (*UserModel, error) {
	var user UserModel
	result := DB.Where("identifier = ?", identifier).First(&user)
	if result.Error != nil {
		return &user, result.Error
	}
	err := bcrypt.CompareHashAndPassword(user.Credential, []byte(credential))
	if err != nil {
		return &user, errors.New("用户名或密码错误")
	}

	return &user, nil
}

// GetFollowList 获取关注列表
func GetFollowList(uid uint) ([]UserModel, error) {
	var user UserModel
	var follows []UserModel
	DB.First(&user, uid)
	err := DB.Model(&user).Association("Follows").Find(&follows)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return follows, err
	}
	return follows, nil
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(uid uint) ([]UserModel, error) {
	var user UserModel
	var follows []UserModel
	DB.First(&user, uid)
	err := DB.Model(&user).Association("Followers").Find(&follows)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return follows, err
	}
	return follows, nil
}

// Follow 关注 uid 对应的用户关注了fid 对应的用户
func Follow(uid uint, fid uint) error {
	var user UserModel
	var userFollow UserModel
	err := DB.First(&userFollow, fid).Error
	if err != nil {
		return err
	}
	return DB.Debug().First(&user, uid).Association("Follows").Append(&userFollow)
}

// Unfollow 取消关注
func Unfollow(uid uint, fid uint) error {
	var followModel FollowModel
	return DB.Debug().Model(&followModel).Where("uid = ? and fid = ?", uid, fid).Delete(&followModel).Error
}

func (u *UserModel) FollowCount() int64 {
	return DB.Debug().Model(&u).Association("Follows").Count()
}

func (u *UserModel) FollowerCount() int64 {
	return DB.Debug().Model(&u).Association("Followers").Count()
}

func (u UserModel) FavoriteCount() int64 {
	return DB.Debug().Model(&u).Association("Likes").Count()
}

func (u UserModel) TotalFavorite() int64 {
	var count int64
	result := DB.Debug().Model(&LikeModel{}).Select("count(*)").Joins("LEFT JOIN video v on v.id = `like`.vid").Where("author_id = ?", u.ID).Group("author_id").Scan(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func (u *UserModel) IsFollow(fid uint) bool {
	var followModel FollowModel
	result := DB.Debug().Model(&followModel).Where("uid = ? and fid = ?", u.ID, fid).Limit(1).Find(&followModel)
	return result.Error == nil
}

func (u *UserModel) WorkCount() int64 {
	return DB.Debug().Model(&u).Association("Videos").Count()
}
