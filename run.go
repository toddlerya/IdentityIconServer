package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/issue9/identicon"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	back  = color.RGBA{R: 255, A: 100}
	fore  = color.RGBA{G: 255, B: 255, A: 100}
	size  = 128
)

func genMd5(str string) string {
	data := []byte(str)
	md5Str := fmt.Sprintf("%x", md5.Sum(data))
	return md5Str
}

func generateIdIcon(str string) string {
	img, err := identicon.Make(size, back, fore, []byte(str))
	if err != nil {
		log.Printf("无法根据输入信息生成头像: %s ERROR: %s\n", str, err)
	}
	imageName := fmt.Sprintf("%s.png", genMd5(str))
	if err != nil {
		log.Printf("无法生成UUID: %s ERROR: %s", str, err)
	}
	fi, err := os.Create(fmt.Sprintf("./images/%s", imageName))
	if err != nil {
		log.Printf("无法创建图片: %s ERROR: %s\n", str, err)
	}
	err = png.Encode(fi, img)
	if err != nil {
		log.Printf("无法编码头像图片: %s ERROR: %s\n", str, err)
	}
	err = fi.Close()
	if err != nil {
		log.Printf("无法关闭图片: %s ERROR: %s\n", str, err)
	}
	return imageName
}

func ServeHTTP(port string) {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	// TODO 后续投入生产要考虑日志分割，日志大小等问题
	f, _ := os.Create("serve.log")

	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()
	router.Use(cors.Default()) // 允许任何服务ajax跨域调用

	router.Static("/image", "./images")
	router.StaticFS("/images", http.Dir("./images"))
	v1 := router.Group("api/v1")
	{
		v1.GET("/getIdImg", getIdImgUrl)
	}

	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil{
		log.Fatalf("在%s端口启动服务失败！\n", port)
	}
}

func getIdImgUrl(c *gin.Context) {
	// todo: 需要对Query参数进行bind，先粗暴的判断下长度
	str := c.Query("str")
	if len(str) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status": "error",
				"code":   "400"},
			"data": gin.H{
				"msg": "请输入有效的参数"},
		})
	} else {
		imageName := generateIdIcon(str)
		idImageUrl := fmt.Sprintf("/images/%s", imageName)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status": "ok",
				"code":   "200"},
			"data": gin.H{
				"IdImageUrl": idImageUrl},
		})
	}

}

func main() {
	ServeHTTP("9981")
}
