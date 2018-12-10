package commonControllers

import (
	//"fmt"
	"net/http"
	"errors"
	"mime/multipart"
	"path"
	"os"
	"io"
	"image"
	"image/jpeg"
	"image/png"
	"image/gif"
	//"strconv"
	"strings"
	"encoding/hex"
	"time"
	"crypto/md5"

	"common/system"
	
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

const FIELD = "upfile"

// 支持的图片后缀名
var supportImageExtNames = []string{".jpg", ".jpeg", ".png", ".ico", ".svg", ".bmp", ".gif"}

func UploadImage(c *gin.Context) {
	imageConfig := system.GetImageConfig()
	var (
		maxUploadSize = imageConfig.Maxsize
		distDir       string
		distFile      string
		err           error
		file          *multipart.FileHeader
		src			  multipart.File
		dist          *os.File
	)

	if file, err = c.FormFile(FIELD); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":1, "msg":"上传文件不能为空"})
		return
	}

	// 判断类型
	extname := strings.ToLower(path.Ext(file.Filename))
	if IsImage(extname) == false {
		c.JSON(http.StatusBadRequest, gin.H{"code":2, "msg":"上传类型错误"})
		return
	}
	// 判断大小
	if file.Size > int64(maxUploadSize) {
		c.JSON(http.StatusBadRequest, gin.H{"code":3, "msg":"文件尺寸超过最大限制"})
		return
	}
	if src, err = file.Open(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":4, "msg":"获取文件内容失败"})
		return
	}
	defer src.Close()
	
	// 内容生成md5key
	hash := md5.New()
	io.Copy(hash, src)
	md5string := hex.EncodeToString(hash.Sum([]byte("")))
	fileName := md5string + extname

	// 根据日期创建文件
	datePath := datePathName()
	distDir = path.Join(imageConfig.Path, datePath)
	distFile = path.Join(distDir, fileName)
	err = os.MkdirAll(distDir, 0775)
    if err != nil {
    	c.JSON(http.StatusBadRequest, gin.H{"code":5, "msg":"创建存储文件失败"})
    	return
    }
	if dist, err = os.Create(distFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":6, "msg":"创建文件失败"})
		return
	}
	defer dist.Close()

	// 获取图片内容
	if src, err = file.Open(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":7, "msg":"获取文件内容失败"})
		return
	}
	io.Copy(dist, src)

	// 生成缩略图
	var thumbPath string
	if imageConfig.Isthumb == 1 {
		// 这里不处理生成失败，不影响主流程
		if thumbPath, err = thumbnailify(distDir, md5string, extname); err != nil {

		}
	}
	/*
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"hash": md5string,
		"filename": fileName,
		"origin": file.Filename,
		"size": file.Size,
		"path": distFile,
		"thumbpath": thumbPath,
	})
	*/
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"url": "/" + distFile,
		"originalName": fileName,
		"size": file.Size,
		"type": extname,
		"state": "SUCCESS",
		"name": fileName,
		"thumbpath": thumbPath,
	})
}

// 生成缩略图
func thumbnailify(imageDir string, filename string, ext string) (outputPath string, err error) {
	var (
		file    *os.File
		img     image.Image
	)

	thumbFilename := filename + "_s" + ext
	outputPath = path.Join(imageDir, thumbFilename)

	// 读取文件
	imagePath := path.Join(imageDir, filename + ext)
	if file, err = os.Open(imagePath); err != nil {
		return
	}
	defer file.Close()

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		break
	case ".png":
		img, err = png.Decode(file)
		break
	case ".gif":
		img, err = gif.Decode(file)
		break
	default:
		err = errors.New("不支持的类型:"+ext)
		return
	}
	if img == nil {
		err = errors.New("生成缩略图失败")
		return
	}

	imageConfig := system.GetImageConfig()
	m := resize.Thumbnail(uint(imageConfig.Thumbwidth), uint(imageConfig.Thumbheight), img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return
	}
	defer out.Close()

	switch ext {
	case ".jpg", ".jpeg":
		jpeg.Encode(out, m, nil)
		break;
	case ".png":
		png.Encode(out, m)
		break;
	case ".gif":
		gif.Encode(out, m, nil)
		break;
	default:
		err = errors.New("不支持的类型:"+ext)
		return
	}
	return
}

// 判断类型是否符合
func IsImage(extName string) bool {
	for i := 0; i < len(supportImageExtNames); i++ {
		if supportImageExtNames[i] == extName {
			return true
		}
	}
	return false
}

// 日期路径
func datePathName() string {
	var now = time.Now()
	var pathName = now.Format("2006") + "/" + now.Format("0102")
	return pathName
}





