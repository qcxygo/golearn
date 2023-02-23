package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

var usersLoginInfo = map[string]User{}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	var num int
	db.QueryRow("select count(*) from tb_user where name=?", username).Scan(&num)

	if num != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	} else {
		var numall int
		numall = 0

		_, err := db.Exec(`INSERT INTO  tb_user ( name,follow_count,follower_count,is_follow,password) VALUES (?,?,?,?,?)`, username, 0, 0, 0, password)
		fmt.Println(err)
		RCSet(token, numall+1, 30*time.Minute)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},

			Token: username + password,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password
	var num int
	db.QueryRow("select count(*) from tb_user where name=? and password=?", username, password).Scan(&num)

	if num > 0 {
		var id int
		db.QueryRow("select id from tb_user where name=? and password=?", username, password).Scan(&id)
		RCSet(token, id, 30*time.Minute)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(id),
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	id := c.Query("id")
	var user User
	row := db.QueryRow("select user_id,name,follow_count,follower_count,is_follow from tb_user where user_id =?", 1)
	row.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &user.IsFollow)
	fmt.Println(id)

	if user.Name != "" {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
