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
	"strconv"
)

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
		Name:       username,
		Identifier: username,
		Credential: hashedPassword,
		Avatar:     "https://static.wikia.nocookie.net/youtube/images/8/8a/Skittle-Chan.jpg/revision/latest?cb=20211122171729",
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
			log.Printf("generate token error: %v\n", err)
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
			log.Printf("generate token error: %v\n", err)
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
	// uid是自己的id
	var uid uint
	if exist {
		uid = claim.(*middleware.Claims).ID
	}
	userIdStr := c.Query("user_id")
	if userIdStr == "" {
		c.JSON(http.StatusOK, UserResponse{
			Response: structs.Response{StatusCode: 1, StatusMsg: "user_id不能为空"},
		})
	}
	userID, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: structs.Response{StatusCode: 1, StatusMsg: "user_id不合法"},
		})
	}
	var user structs.User
	err = db.GetUser(&user, uint(userID), uid)
	log.Printf("user: %v\n", user)
	if err == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: structs.Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: structs.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	}
}
