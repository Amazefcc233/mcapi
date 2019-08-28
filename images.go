package main

import (
	"fmt"
	"image"
	_ "image/png"
	"strconv"
	"time"
	"flag" // 添
    "io/ioutil" 
    "log" 

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	// "golang.org/x/image/font/inconsolata"
	"github.com/golang/freetype/truetype"
)

var ( // 添
	fontfile = flag.String("fontfile", "msyh.ttc", "filename of the ttf font")
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")  
    size     = flag.Float64("size", 16, "font size in points")	
)

const (
	imageWidth  = 400 //325
	imageHeight = 64
)

const (
	imageBlockWidth = 64
	fromImage       = 4
	offsetText      = float64(imageBlockWidth + fromImage)
)

func respondServerImage(c *gin.Context) {
	c.Request.ParseForm()

	ip := c.Request.Form.Get("ip")
	port := c.Request.Form.Get("port")
	title := c.Request.Form.Get("title")
	theme := c.Request.Form.Get("theme")

	var serverAddr string
	var serverDisp string

	if port == "" {
		serverAddr = ip + ":25565"
		serverDisp = ip
	} else {
		serverAddr = ip + ":" + port
		serverDisp = serverAddr
	}

	if title != "" {
		serverDisp = title
	}
	
	// 59-73添
	fontBytes, err := ioutil.ReadFile(*fontfile) 
    if err != nil { 
        log.Println(err) 
        return 
    } 
	font, err := truetype.Parse(fontBytes)
	if err != nil {
        log.Println(err) 
        return 
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: *size,
		DPI:*dpi,
	}) 

	status := getStatusFromCacheOrUpdate(serverAddr, c, true)

	if status == nil {
		dc := gg.NewContext(imageWidth, imageHeight)

		// dc.SetFontFace(inconsolata.Regular8x16)
		dc.SetFontFace(face)
		if theme == "dark" {
			dc.SetRGB(1, 1, 1)
		} else {
			dc.SetRGB(0, 0, 0)
		}

		dc.DrawStringAnchored("Too many bad requests.", imageWidth/2, imageHeight/2, 0.5, 0.5)

		dc.EncodePNG(c.Writer)
	}

	var imgToDraw image.Image

	if status.Favicon == "" {
		img, err := gg.LoadPNG("files/grass_sm.png")
		if err != nil {
			c.Error(err)
			return
		}

		imgToDraw = img
	} else {
		img, err := status.Image()
		if err != nil {
			c.Error(err)
			return
		}

		imgToDraw = img
	}

	bounds := imgToDraw.Bounds()
	height, width := bounds.Dy(), bounds.Dx()

	dc := gg.NewContext(imageWidth, imageHeight)

	dc.DrawImage(imgToDraw, (imageBlockWidth-width)/2, (imageHeight-height)/2)

	// dc.SetFontFace(inconsolata.Regular8x16)
	dc.SetFontFace(face)
	if theme == "dark" {
		dc.SetRGB(1, 1, 1)
	} else {
		dc.SetRGB(0, 0, 0)
	}
	_, tH := dc.MeasureString(serverDisp)
	dc.DrawString(serverDisp, offsetText, tH)

	lastHeight := 1 + tH

	var online string

	if status.Online {
		online = "在线!"
	} else {
		online = "离线"
	}

	tW, tH := dc.MeasureString(online)
	dc.DrawString(online, offsetText, lastHeight+tH+2)

	lastHeight += tH + 2

	if status.Online {
		msg := fmt.Sprintf("    %d/%d 玩家", status.Players.Now, status.Players.Max)
		_, tH = dc.MeasureString(msg)
		dc.DrawString(msg, float64(width+fromImage*2)+tW, lastHeight)
	}

	i, _ := strconv.ParseInt(status.LastUpdated, 10, 64)
	last := time.Unix(i, 0)
	minutesAgo := int(time.Now().Sub(last).Minutes())

	plural := ""
	if minutesAgo != 1 {
		plural = "s"
	}

	msg := fmt.Sprintf("数据于 %d min%s 前更新 • 代码源自 mcapi.us", minutesAgo, plural)

	dc.DrawString(msg, offsetText, imageHeight-4)

	dc.EncodePNG(c.Writer)
}
