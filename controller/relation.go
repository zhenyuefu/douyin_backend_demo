package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/structs"
	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, UserListResponse{
		Response: structs.Response{
			StatusCode: 0,
		},
		UserList: []structs.User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {

	var users = db.UserModel{}
	err := db.DB.Select("Id", "Name", "FollowCount", "FollowerCount", "IsFollow",
		"Avatar").Find(&users, 1).Error
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Get user list failed"})
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: structs.Response{
			StatusCode: 0,
		},
		UserList: []structs.User{DemoUser},
	})
}
