package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init init DB
func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(constants.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

	m := DB.Migrator()

	if err = m.AutoMigrate(&UserModel{}, &VideoModel{}, &LikeModel{}, &CommentModel{}); err != nil {
		panic(err)
	}

}
