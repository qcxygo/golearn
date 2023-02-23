package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	vid := c.Query("video_id")
	quality := RCGet(token)
	uid, _ := quality.Int64()

	action_type := c.Query("action_type")
	var numall int

	db.QueryRow("select count(*) from tb_favourite where  user_id=? and video_id=?", uid, vid).Scan(&numall)
	fmt.Println(numall)
	var fn int

	db.QueryRow("select favourite_count from tb_video where video_id=?", vid).Scan(&fn)
	fmt.Println(fn)
	if action_type == "1" {
		if numall > 0 {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点过赞了"})
		} else {
			_, err := db.Exec("INSERT INTO tb_favourite(user_id,video_id) VALUES (?,?)", uid, vid)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
				fmt.Println(err)
				return
			}
			_, err = db.Exec("update  tb_video set favourite_count =? where video_id=?", fn+1, vid)

			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "点赞失败"})
				fmt.Println(err)
				return
			} else {
				c.JSON(http.StatusOK, Response{StatusCode: 0})
				return
			}
		}
	} else if action_type == "2" {
		if numall < 1 {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "没点过赞"})
			return
		} else {
			_, err := db.Exec("delete from  tb_favourite where user_id=? and video_id=?", uid, vid)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "取消失败"})
				return
			}
			_, err = db.Exec("update  tb_video set favourite_count =? where video_id=?", fn-1, vid)
			if err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "取消失败"})
				return
			} else {
				c.JSON(http.StatusOK, Response{StatusCode: 0})
				return
			}
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "状态异常"})
	}
}

func FavoriteList(c *gin.Context) {
	var Videos = []Video{}
	//token := c.Query("token")
	uid := c.Query("user_id")
	vids := []int{}
	rows, err := db.Query("select video_id from tb_favourite where user_id=?", uid)
	fmt.Println(err)
	for rows.Next() {
		var vd int
		rows.Scan(&vd)

		vids = append(vids, vd)
	}

	for i := 0; i < len(vids); i++ {
		var temv Video
		var uidt int
		fmt.Println(vids[i])
		vrows := db.QueryRow("select video_id,user_id,play_url,cover_url,favourite_count,comment_count,title from tb_video where video_id=?", vids[i])
		fmt.Println(vrows)
		vrows.Scan(&temv.Id, &uidt, &temv.PlayUrl, &temv.CoverUrl, &temv.FavoriteCount, &temv.CommentCount, &temv.Title)

		urows := db.QueryRow("select user_id,name,follow_count,follower_count from tb_user where user_id=?", uidt)
		var tuser User
		urows.Scan(&tuser.Id, &tuser.Name, &tuser.FollowCount, &tuser.FollowerCount)
		tuser.IsFollow = false
		temv.Author = tuser
		Videos = append(Videos, temv)

	}
	fmt.Println(Videos)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: Videos,
	})
}
