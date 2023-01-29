package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"simple_tiktok/logger"
	"simple_tiktok/models"
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
func UploadCOS(c *gin.Context, srcFile multipart.File, head *multipart.FileHeader, title string, userid uint64) (int, string) {
	// 1、判断文件是否为视频
	// 读取后缀名以及相关信息
	suffix := path.Ext(head.Filename)
	// 判断后缀是否为视频后缀
	if suffix != ".avi" && suffix != ".wmv" && suffix != ".mpg" && suffix != ".mpeg" && suffix != ".flv" && suffix != ".mov" && suffix != ".rm" && suffix != ".ram" && suffix != ".swf" && suffix != ".mp4" {
		return -1, "请上传视频文件"
	}

	// 2、将视频存到腾讯云COS
	// 生成identity
	identity, err := utils.GetID()
	if err != nil {
		logger.SugarLogger.Error("GetID Error:" + err.Error())
		fmt.Println("GetID Error:" + err.Error())
		return -1, "投稿失败"
	}

	// 视频保存名
	filename := strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(int(identity))
	key := filename + suffix
	fmt.Println(identity)
	fmt.Println(filename)

	_, err = utils.COSClient.Object.Put(c, key, srcFile, nil)
	if err != nil {
		logger.SugarLogger.Error("Put Video Error：" + err.Error())
		fmt.Println("Put Video Error：" + err.Error())
		return -1, "投稿失败"
	}

	// 3、使用腾讯云COS接口截取视频封面，保存到本地服务器下
	// 1）在本地文件夹下创建文件
	dir := viper.GetString("uplaodBase") + strconv.Itoa(int(userid))
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		logger.SugarLogger.Error("MkdirAll Error：" + err.Error())
		fmt.Println("MkdirAll Error：" + err.Error())
		return -1, "投稿失败"
	}

	path := dir + "/" + "cover" + filename + ".png"
	fd, err := os.Create(path)
	if err != nil {
		logger.SugarLogger.Error("File Create Error：" + err.Error())
		fmt.Println("File Create Error：" + err.Error())
		return -1, "投稿失败"
	}

	// 2）读取COS的封面信息，保存到本地
	opt := &cos.GetSnapshotOptions{
		Time:   1,
		Format: "png",
	}

	resp, err := utils.COSClient.CI.GetSnapshot(c, key, opt)
	if err != nil {
		logger.SugarLogger.Error("GetSnapshot Error：" + err.Error())
		fmt.Println("GetSnapshot Error：" + err.Error())
		return -1, "投稿失败"
	}
	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		logger.SugarLogger.Error("io.Copy Error：" + err.Error())
		fmt.Println("io.Copy Error：" + err.Error())
		return -1, "投稿失败"
	}
	fd.Close()

	// 4、保存数据到数据库
	coveurl := path[1:]
	videoInfo := models.VideoBasic{
		Identity:     identity,
		UserIdentity: userid,
		PlayUrl:      key,
		CoverUrl:     coveurl,
		Title:        title,
		PublishTime:  time.Now(),
	}

	fmt.Println(videoInfo)
	return 0, "投稿成功"
}
