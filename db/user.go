package db

import (
	"errors"
	"github.com/RaymondCode/simple-demo/constants"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Identifier string `gorm:"uniqueIndex;size:32"`
	Credential []byte
	Name       string        `json:"name"`
	Follows    []FollowModel `gorm:"foreignKey:UID"`
	Followers  []FollowModel `gorm:"foreignKey:FID"`
	Videos     []VideoModel  `gorm:"foreignKey:AuthorID"`
	Likes      []LikeModel   `gorm:"foreignKey:UID"`
}

func (u *UserModel) TableName() string {
	return constants.UserTableName
}

// CreateUser 创建用户
func CreateUser(user *UserModel) error {
	return DB.Create(user).Error
}

// GetUser 获取用户
func GetUser(user *UserModel, uid uint) error {
	return DB.First(user, uid).Error
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
