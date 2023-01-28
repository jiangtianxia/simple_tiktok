package utils

import (
	"fmt"
	"os/exec"
	"simple_tiktok/logger"
)

/**
 * @Author jiang
 * @Description 使用ffmpeg截取视频封面
 * @Date 10:00 2023/1/26
 **/
/**
 * 运行cmd命令
 **/
func CallCommandRun(cmdName string, args []string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	fmt.Println("CallCommand Run 参数=> ", args)
	fmt.Println("CallCommand Run 执行命令=> ", cmd)
	bytes, err := cmd.Output()
	if err != nil {
		fmt.Println("CallCommand Run 出错了.....", err.Error())
		logger.SugarLogger.Error("CallCommand Run 出错了.....", err.Error())
		return "", err
	}
	resp := string(bytes)
	fmt.Println(resp)
	fmt.Println("CallCommand Run 调用完成.....")
	return resp, nil
}

/**
 * 根据视频URL调用ffmpeg截取视频封面
 * 参数：
 * 		-i video_url    		需要截取的视频url
 *		-ss 00:00:01    		截取的视频时间点
 *		-vframes 1 image_url 	将截取的封面保存到该url
 *		-f image2				使用image2格式进行解析
 *		-vcodec png				截取封面图为png格式
 **/
func GetIpcScreenShot(ffmpegPath string, url string, screenShotPath string) (string, error) {
	var params []string
	// ffmpeg -i video_url -ss 00:00:01 -vframes 1 image_url -f image2 -vcodec png
	params = append(params, "-i")
	params = append(params, url)
	params = append(params, "-ss")
	params = append(params, "00:00:01")
	params = append(params, "-vframes")
	params = append(params, "1")
	params = append(params, screenShotPath)
	params = append(params, "-f")
	params = append(params, "image2")
	params = append(params, "-vcodec")
	params = append(params, "png")

	resp, err := CallCommandRun(ffmpegPath, params)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("获取截图出错，url为--->", url)
		logger.SugarLogger.Error("获取截图出错，url为--->", url, "Error为--->", err.Error())
		return "", err
	}
	return resp, nil
}

// func main() {
// 	GetIpcScreenShot("ffmpeg", "./test/test.mp4", "./test.png")
// }
