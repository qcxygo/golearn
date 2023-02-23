package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	cid := c.Query("comment_id")

	quality := RCGet(token)

	_, err := quality.Result()
	if err != nil {
		panic("get id failed,err: " + err.Error())
		return
	}

	if actionType == "1" {
		text := c.Query("comment_text")
		uid, err := quality.Int64()
		vid := c.Query("video_id")
		var numall int

		db.QueryRow("select count(*) from comment").Scan(&numall)
		_, err = db.Exec(`INSERT INTO comment(comment_id,user_id,video_id,content,deleted) VALUES (?,?,?,?,?)`, numall+1, uid, vid, text, 0)
		urows := db.QueryRow("select user_id,name,follow_count,follower_count from tb_user where user_id=?", uid)
		var tuser User
		urows.Scan(&tuser.Id, &tuser.Name, &tuser.FollowCount, &tuser.FollowerCount)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
			return
		}
		timeStr := time.Now().Format("2006-01-02 15:04:05")

		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: Comment{
				Id:         int64(numall + 1),
				User:       tuser,
				Content:    text,
				CreateDate: timeStr,
			}})
		return
	} else {
		_, err = db.Exec("update  comment set deleted =? where comment_id=?", 1, cid)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "delete wrong"})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	var comments = []Comment{}
	vid := c.Query("video_id")
	crows, _ := db.Query("select user_id,comment_id,content,ctime from comment where video_id=? and deleted=0", vid)
	for crows.Next() {
		var temc Comment
		var tuid int
		var t time.Time
		crows.Scan(&tuid, &temc.Id, &temc.Content, &t)
		urows := db.QueryRow("select user_id,name,follow_count,follower_count from tb_user where user_id=?", tuid)
		var tuser User
		urows.Scan(&tuser.Id, &tuser.Name, &tuser.FollowCount, &tuser.FollowerCount)

		temc.CreateDate = t.Format("2006-01-02 15:04:05")
		temc.User = tuser
		comments = append(comments, temc)

	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
