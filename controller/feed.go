package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	token := c.Query("token")
	t := c.Query("next_time")
	tt, _ := strconv.Atoi(t)

	tm := time.Unix(int64(tt), 0)
	fmt.Println(int64(tt))
	lt := tm.Format("2006-01-02 15:04:05")

	fmt.Println(lt)
	var Videos = []Video{}
	var ltime string
	quality := RCGet(token)
	uid, _ := quality.Int64()
	var v *sql.Rows
	if t != "" {
		vs, err := db.Query("select video_id,user_id,play_url,cover_url,favourite_count,comment_count,title,update_date from tb_video where update_date<DATE_FORMAT(?,\"%Y-%m-%d %H:%i:%s\") ORDER BY update_date DESC limit 0,30", lt)
		fmt.Println(err)
		v = vs
	} else {
		vs, err := db.Query("select video_id,user_id,play_url,cover_url,favourite_count,comment_count,title,update_date from tb_video  ORDER BY update_date DESC limit 0,30")

		fmt.Println(err)
		v = vs
	}
	for v.Next() {
		var uidt int
		var fn int
		var temv Video
		v.Scan(&temv.Id, &uidt, &temv.PlayUrl, &temv.CoverUrl, &temv.FavoriteCount, &temv.CommentCount, &temv.Title, &ltime)
		urows := db.QueryRow("select user_id,name,follow_count,follower_count from tb_user where user_id=?", uidt)
		var tuser User
		urows.Scan(&tuser.Id, &tuser.Name, &tuser.FollowCount, &tuser.FollowerCount)
		db.QueryRow("select count(*) from favourite_tb where user_id=? and video_id=", uid, temv.Id).Scan(fn)
		if fn == 0 {
			temv.IsFavorite = false
		} else {
			temv.IsFavorite = true
		}
		tuser.IsFollow = false
		temv.Author = tuser
		Videos = append(Videos, temv)

	}
	fmt.Println(ltime)
	fmt.Println(ltime[:10] + " " + ltime[11:19])
	tm2, _ := time.ParseInLocation("2006-01-02 15:04:05", ltime[:10]+" "+ltime[11:19], time.Local)
	fmt.Println(tm2)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: Videos,
		NextTime:  tm2.Unix(),
	})
}
