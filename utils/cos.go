package utils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

/**
 * @Author jiang
 * @Description 初始化COS客户端
 * @Date 9:00 2023/1/29
 **/
var COSClient *cos.Client

func InitCos() {
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶 region 可以在 COS 控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(viper.GetString("cos.addr"))
	b := &cos.BaseURL{BucketURL: u}
	COSClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			// SecretID: os.Getenv("SECRETID"), // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			SecretID: viper.GetString("cos.SecretID"),
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			// SecretKey: os.Getenv("SECRETKEY"), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			SecretKey: viper.GetString("cos.SecretKey"),
		},
	})

	fmt.Println("Cos client inited ...... ")
}
