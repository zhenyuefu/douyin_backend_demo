package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/structs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	structs.Response
	UserList []structs.User `json:"user_list"`
}

// RelationAction 关注取关
func RelationAction(c *gin.Context) {
	//获取uid
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	uid := user.(*middleware.Claims).ID
	//获取to_user_id
	toUserId, err := strconv.Atoi(c.Query("to_user_id"))
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "to_user_id is not valid"})
		return
	}
	actionType := c.Query("action_type")
	switch actionType {
	//关注
	case "1":
		{
			if err := db.Follow(uid, uint(toUserId)); err != nil {
				c.JSON(http.StatusOK, structs.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, structs.Response{StatusCode: 0})
			return
		}
		//取关
	case "2":
		err := db.Unfollow(uid, uint(toUserId))
		if err != nil {
			c.JSON(http.StatusOK, structs.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
	}
}

// FollowList 关注列表
func FollowList(c *gin.Context) {
	// 查询的用户的uid
	uidStr := c.Query("user_id")
	uid, err := strconv.Atoi(uidStr)

	// loginUID
	token := c.Query("token")
	var loginUID uint
	if token != "" {
		claims, err := middleware.ParseToken(token)
		if err != nil {
			log.Printf(err.Error())
		}
		loginUID = claims.ID
	}

	users, err := db.GetFollowList(uint(uid))
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Get user list failed"})
		return
	}
	var userList []structs.User
	for _, user := range users {
		var u []db.UserModel
		// 这里是自己uid
		err := db.DB.Model(&user).Association("Followers").Find(&u, loginUID)
		if err != nil {
			log.Println("查询follow:" + err.Error())
		}
		userList = append(userList, structs.User{
			ID:       user.ID,
			Name:     user.Name,
			Avatar:   user.Avatar,
			IsFollow: len(u) > 0,
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

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {
	// 查询的用户的uid
	uidStr := c.Query("user_id")
	uid, err := strconv.Atoi(uidStr)

	// loginUID
	token := c.Query("token")
	var loginUID uint
	if token != "" {
		claims, err := middleware.ParseToken(token)
		if err != nil {
			log.Printf(err.Error())
		}
		loginUID = claims.ID
	}

	users, err := db.GetFollowerList(uint(uid))
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Get user list failed"})
		return
	}
	var userList []structs.User
	for _, user := range users {
		var u []db.UserModel
		// 1 是自己uid
		err := db.DB.Debug().Model(&user).Association("Followers").Find(&u, loginUID)
		if err != nil {
			log.Println("查询follower:" + err.Error())
		}
		userList = append(userList, structs.User{
			ID:       user.ID,
			Name:     user.Name,
			Avatar:   user.Avatar,
			IsFollow: len(u) > 0,
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
