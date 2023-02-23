package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	quality := RCGet(token)
	_, err := quality.Result()
	if err != nil {
		panic("get id failed,err: " + err.Error())
		return
	}
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	file, err := data.Open()

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	defer file.Close()
	FilePath := "./upload/"

	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + data.Filename
	fw, err := os.Create(FilePath + fileName)
	if err != nil {
		panic(err)
	}
	defer fw.Close()
	_, err = io.Copy(fw, file)
	_, err = Cos.Object.PutFromFile(context.Background(), "/videos/"+fileName, FilePath+fileName, nil)
	if err != nil {
		panic(err)
	}

	videoPath := FilePath + fileName
	picpath := FilePath + fileName
	buf := bytes.NewBuffer(nil)
	err = ffmpeg_go.Input(videoPath).
		Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	img, err := imaging.Decode(buf)

	err = imaging.Save(img, picpath+".jpeg")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
	}

	names := strings.Split(picpath, "/")
	fileImage := names[len(names)-1] + ".jpeg"

	_, err = Cos.Object.PutFromFile(context.Background(), "/pictures/"+fileImage, picpath+".jpeg", nil)
	if err != nil {
		log.Fatal("存储失败：", err)
	}
	id, err := quality.Int64()
	var numall int
	numall = 0
	db.QueryRow("select count(*) from tb_video").Scan(&numall)
	var uname string
	db.QueryRow("select name from tb_user where user_id=?", id).Scan(&uname)
	_, err = db.Exec(`INSERT INTO tb_video(video_id,user_name,user_id,play_url,cover_url,title) VALUES (?,?,?,?,?,?)`, strconv.Itoa(numall+1), uname, strconv.Itoa(int(id)), "https://niganma-1301440536.cos.ap-guangzhou.myqcloud.com/videos/"+fileName, "https://niganma-1301440536.cos.ap-guangzhou.myqcloud.com/pictures/"+fileImage, title)
	if err != nil {
		panic("get id failed,err: " + err.Error())
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  fileName + " uploaded successfully",
	})

}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	uid := c.Query("user_id")
	var Videos = []Video{}
	vrows, _ := db.Query("select video_id,user_id,play_url,cover_url,favourite_count,comment_count,title from tb_video where user_id=?", uid)
	for vrows.Next() {
		var temv Video
		var uidt int
		vrows.Scan(&temv.Id, &uidt, &temv.PlayUrl, &temv.CoverUrl, &temv.FavoriteCount, &temv.CommentCount, &temv.Title)

		urows := db.QueryRow("select user_id,name,follow_count,follower_count from tb_user where user_id=?", uidt)
		var tuser User
		urows.Scan(&tuser.Id, &tuser.Name, &tuser.FollowCount, &tuser.FollowerCount)

		temv.Author = tuser
		Videos = append(Videos, temv)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: Videos,
	})
}
