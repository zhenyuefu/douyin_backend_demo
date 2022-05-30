package controller

import (
	"errors"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/structs"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var usersLoginInfo = map[string]structs.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	structs.Response
	UserId uint   `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	structs.Response
	User structs.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := db.UserModel{
		Name:       "testName",
		Identifier: username,
		Credential: hashedPassword,
	}

	err := db.CreateUser(&user)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: structs.Response{StatusCode: 1, StatusMsg: "用户已经注册"},
			})
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: structs.Response{StatusCode: 1, StatusMsg: "服务器错误"},
		})
	} else {
		token, err := middleware.GenerateToken(user.ID, user.Identifier)
		if err != nil {
			log.Fatalf("generate token error: %v", err)
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: structs.Response{StatusCode: 0},
			UserId:   user.ID,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	user, err := db.VerifyCredential(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: structs.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		token, err := middleware.GenerateToken(user.ID, user.Identifier)
		if err != nil {
			log.Fatalf("generate token error: %v", err)
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: structs.Response{StatusCode: 0},
			UserId:   user.ID,
			Token:    token,
		})
	}
}

func UserInfo(c *gin.Context) {
	claim, exist := c.Get("user")
	var uid uint
	var user db.UserModel
	if exist {
		uid = claim.(*middleware.Claims).ID
	}
	err := db.GetUser(&user, uid)
	if err == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: structs.Response{StatusCode: 0},
			User: structs.User{
				Id:            user.ID,
				Name:          user.Identifier,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      false,
				Avatar:        "http://momentcon-1255653016.file.myqcloud.com/2110464774/20001/84192D32EF6436116E845A18A22FB074.png",
			},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: structs.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	}
}
