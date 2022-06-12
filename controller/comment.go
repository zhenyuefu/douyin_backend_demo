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

type CommentListResponse struct {
	structs.Response
	CommentList []structs.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	structs.Response
	Comment structs.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	//获取uid
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	uid := user.(*middleware.Claims).ID

	actionType := c.Query("action_type")
	vidStr := c.Query("video_id")
	vid, err := strconv.Atoi(vidStr)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Video id is not valid"})
	}
	switch actionType {
	// 创建评论
	case "1":
		{
			commentText := c.Query("comment_text")
			comment := db.CommentModel{
				UID:     uid,
				VID:     uint(vid),
				Content: commentText,
			}
			err := db.CreateComment(&comment)
			if err != nil {
				c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: err.Error()})
				return
			}
			var user structs.User
			err = db.GetUser(&user, uid, uid)
			if err != nil {
				log.Println(err)
			}
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: structs.Response{
					StatusCode: 0,
					StatusMsg:  "success",
				},
				Comment: structs.Comment{
					ID:         comment.ID,
					User:       user,
					Content:    comment.Content,
					CreateDate: comment.CreatedAt.Format("01-02"),
				},
			})
		}
	// 删除评论
	case "2":
		{
			commentIdStr := c.Query("comment_id")
			commentId, err := strconv.Atoi(commentIdStr)
			if err != nil {
				c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Comment id is not valid"})
				return
			}
			err = db.DeleteComment(uint(commentId))
			if err != nil {
				c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Delete comment failed"})
				return
			}
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: structs.Response{
					StatusCode: 0,
					StatusMsg:  "success",
				},
			})
		}
	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	var uid uint
	if token != "" {
		claims, err := middleware.ParseToken(token)
		if err != nil {
			log.Printf(err.Error())
		}
		uid = claims.ID
	}

	vidStr := c.Query("video_id")
	vid, err := strconv.Atoi(vidStr)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Video id is not valid"})
		return
	}
	var commentList []db.CommentModel
	err = db.GetComment(&commentList, uint(vid))
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	var comments []structs.Comment
	for _, comment := range commentList {
		var user structs.User
		err := db.GetUser(&user, comment.UID, uid)
		if err != nil {
			log.Println(err)
		}
		comments = append(comments, structs.Comment{
			ID:         comment.ID,
			User:       user,
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response: structs.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		CommentList: comments,
	})
}
