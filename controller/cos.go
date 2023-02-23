package controller

import (
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var Cos *cos.Client

func init() {
	u, _ := url.Parse("")
	bucket := &cos.BaseURL{BucketURL: u}
	Cos = cos.NewClient(bucket, &http.Client{Transport: &cos.AuthorizationTransport{
		SecretID:  "",
		SecretKey: "",
	}})
}
