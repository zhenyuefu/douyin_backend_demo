package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/structs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserListResponse struct {
	structs.Response
	UserList []structs.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	users, err := db.GetFollowList(1)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Get user list failed"})
		return
	}
	var userList []structs.User
	for _, user := range users {
		var u []db.UserModel
		// 1 是自己uid
		err := db.DB.Model(&user).Association("Follows").Find(&u, 1)
		if err != nil {
			log.Println("查询follow:" + err.Error())
		}
		userList = append(userList, structs.User{
			Id:            user.ID,
			Name:          user.Identifier,
			FollowCount:   db.DB.Model(&user).Association("Follows").Count(),
			FollowerCount: db.DB.Model(&user).Association("Followers").Count(),
			IsFollow:      len(u) > 0,
		})
	}
	log.Println(userList)
	c.JSON(http.StatusOK, UserListResponse{
		Response: structs.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	//连接数据库
	//var DB *gorm.DB
	//var err error
	//db.Init()
	//
	//m := DB.Migrator()
	//if err = m.AutoMigrate(&User{}); err != nil {
	//	panic(err)
	//}
	//
	//var users = User{}
	//err = DB.Select("Id", "Name", "FollowCount", "FollowerCount", "IsFollow","Avatar").Find(&users, 1).Error
	//if err != nil {
	//	c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Get user list failed"})
	//	return
	//}

	users, err := db.GetFollowList(1)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Get user list failed"})
		return
	}
	var userList []structs.User
	for _, user := range users {
		var u []db.UserModel
		// 1 是自己uid
		err := db.DB.Model(&user).Association("Followers").Find(&u, 1)
		if err != nil {
			log.Println("查询follower:" + err.Error())
		}
		userList = append(userList, structs.User{
			Id:            user.ID,
			Name:          user.Identifier,
			FollowCount:   db.DB.Model(&user).Association("Follows").Count(),
			FollowerCount: db.DB.Model(&user).Association("Followers").Count(),
			IsFollow:      len(u) > 0,
		})
	}
	log.Println(userList)
	c.JSON(http.StatusOK, UserListResponse{
		Response: structs.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}
