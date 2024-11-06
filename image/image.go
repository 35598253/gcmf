package gcmf

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/gogf/gf/v2/os/gfile"
	"golang.org/x/image/bmp"

	"github.com/nfnt/resize"

	"github.com/gogf/gf/v2/net/ghttp"
)

type imgFile struct {
	Img image.Image
	Ext string
}

// OutImage 导出接口
type OutImage interface {
	TextWater(Text string, FontPath ...string) (*imgFile, error)
	ImgWater(logoPath string, size uint) (*imgFile, error)
	ToFile(RootPath, outPath string) (string, error)
	ToBase64() (string, error)
}

func NewImageFile(imgPath string) (OutImage, error) {
	img, err := os.OpenFile(imgPath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open Img `%s`: %w", imgPath, err)
	}

	defer func() {
		_ = img.Close()
	}()
	_img, _ext, err := image.Decode(img)
	if err != nil {
		return nil, err
	}
	return &imgFile{
		Img: _img,
		Ext: _ext,
	}, nil
}
func NewImageUpload(UpImg *ghttp.UploadFile, imgResize uint) (OutImage, error) {

	ff, err := UpImg.Open()
	if err != nil {
		return nil, err
	}

	defer func(ff multipart.File) {
		_ = ff.Close()
	}(ff)
	img, exts, err := image.Decode(ff)
	if err != nil {
		return nil, err
	}
	//读取图片宽度
	x := img.Bounds().Max.X

	if x > int(imgResize) {
		img = resize.Resize(imgResize, 0, img, resize.NearestNeighbor)
	}
	return &imgFile{
		Img: img,
		Ext: exts,
	}, nil

}
func (s *imgFile) ToFile(RootPath, outPath string) (string, error) {

	// 相对路径
	_tPath := RootPath + outPath
	name := strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6))
	imageName := name + "." + s.Ext
	savePath := _tPath + imageName
	outFile := outPath + imageName
	buf, err := gfile.Create(savePath)

	if err != nil {
		return "", err
	}
	defer func(buf multipart.File) {
		_ = buf.Close()
	}(buf)

	err = imgEncode(buf, s.Img, s.Ext)
	if err != nil {
		return "", err
	}
	return outFile, nil
}
func (s *imgFile) ToBase64() (string, error) {

	buf := new(bytes.Buffer)
	err := imgEncode(buf, s.Img, s.Ext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil

}

func imgEncode(out io.Writer, img image.Image, Ext string) error {

	var err error
	switch Ext {
	case "jpeg":
		err = jpeg.Encode(out, img, nil)
	case "png":
		err = png.Encode(out, img)
	case "gif":
		err = gif.Encode(out, img, nil)
	case "bmp":
		err = bmp.Encode(out, img)
	default:
		err = errors.New("image type error")
	}
	return err
}
