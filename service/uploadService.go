package service

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

/**
 * @Author jiang
 * @Description 上传视频到腾讯云COS
 * @Date 16:00 2023/1/28
 **/
func UploadCOS(c *gin.Context, srcFile multipart.File, head *multipart.FileHeader, title string, userid int, username string) (int, string) {
	// 1、判断文件是否为视频
	// 读取后缀名以及相关信息
	suffix := path.Ext(head.Filename)
	// 判断后缀是否为视频后缀
	if suffix != ".avi" && suffix != ".wmv" && suffix != ".mpg" && suffix != ".mpeg" && suffix != ".flv" && suffix != ".mov" && suffix != ".rm" && suffix != ".ram" && suffix != ".swf" && suffix != ".mp4" {
		return -1, "请上传视频文件"
	}

	// 2、将视频存到腾讯云COS
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶 region 可以在 COS 控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(viper.GetString("cos.addr"))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
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

	// 视频保存名
	// 生成identity
	identity, err := utils.GetID()
	if err != nil {
		logger.SugarLogger.Error("GetID Error:" + err.Error())
		fmt.Println("GetID Error:" + err.Error())
		return -1, "投稿失败"
	}

	filename := strconv.Itoa(int(time.Now().Unix()) + int(identity))
	key := filename + suffix

	res, err := client.Object.Put(c, key, srcFile, nil)
	if err != nil {
		logger.SugarLogger.Error("Put Video Error：" + err.Error())
		fmt.Println("Put Video Error：" + err.Error())
		return -1, "投稿失败"
	}

	fmt.Println(res.Request.URL)
	return 0, "投稿成功"
}
