package gcmf

import (
	"context"
	"fmt"
	gcmf "github.com/35598253/gcmf/image"
	"path"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Outfile 返回格式
type Outfile struct {
	OldName string `json:"oldName"`
	URL     string `json:"url"`
	Size    int64  `json:"size"`
}
type UploadCfg struct {
	FileExts    string `json:"fileExts"`
	ImgExts     string `json:"imgExts"`
	MaxSize     int64  `json:"maxSize"`
	ImageResize uint   `json:"imageResize"`
}
type uploadConfig struct {
	RealPath  string `json:"realPath"`
	Url       string `json:"url"`
	UploadDir string `json:"uploadDir"`
	*UploadCfg
}

type Uploader interface {
	File(UpFile *ghttp.UploadFile) (*Outfile, error)
	Image(UpImg *ghttp.UploadFile, WaterMark string) (*Outfile, error)
}

// Upload file
func Upload(Ctx context.Context, RealPath, UploadDir, Url string, UpCfg ...*UploadCfg) Uploader {
	cfg, _ := GetConfig(Ctx)

	_cfg := &uploadConfig{
		RealPath:  RealPath,
		Url:       Url,
		UploadDir: UploadDir,
	}
	if len(UpCfg) > 0 && UpCfg[0] != nil {
		dfCfg := UpCfg[0]

		if dfCfg.ImgExts == "" {
			dfCfg.ImgExts = cfg.Upload.ImgExts
		}
		if dfCfg.FileExts == "" {
			dfCfg.FileExts = cfg.Upload.FileExts
		}
		if dfCfg.ImageResize == 0 {
			dfCfg.ImageResize = cfg.Upload.ImageResize
		}
		if dfCfg.MaxSize == 0 {
			dfCfg.MaxSize = cfg.Upload.MaxSize
		}
		_cfg.UploadCfg = dfCfg
	}

	return _cfg
}

// File 上传文件
func (u *uploadConfig) File(UpFile *ghttp.UploadFile) (*Outfile, error) {
	var tPath = fmt.Sprintf("/file/%s/%s/", u.UploadDir, time.Now().Format("2006"))
	var upPath = u.RealPath + tPath
	var outPath = u.Url + tPath

	// 读取默认

	if err := checkFile(UpFile, u.FileExts, u.MaxSize); err != nil {
		return nil, err
	}
	// 保存附件
	filePath, err := UpFile.Save(upPath, true)
	if err != nil {
		return nil, err
	}
	// 添加到附件数据库中
	return &Outfile{OldName: UpFile.Filename, URL: outPath + filePath, Size: UpFile.Size}, nil
}

// Image 图片上传
func (u *uploadConfig) Image(UpImg *ghttp.UploadFile, WaterMark string) (*Outfile, error) {

	if err := checkFile(UpImg, u.ImgExts, u.MaxSize); err != nil {
		return nil, err
	}
	imgOutPath := fmt.Sprintf("/images/%s/%s/", u.UploadDir, time.Now().Format("2006"))
	// 剪裁&水印
	imgTmp, err := gcmf.NewImageUpload(UpImg, u.ImageResize)
	if err != nil {
		return nil, err
	}
	var _tPath string
	if WaterMark != "" {

		_tImg, err := imgTmp.TextWater(WaterMark)
		if err != nil {
			return nil, err
		}
		_tPath, err = _tImg.ToFile(u.RealPath, imgOutPath)
		if err != nil {
			return nil, err
		}
	} else {
		_tPath, err = imgTmp.ToFile(u.RealPath, imgOutPath)
		if err != nil {
			return nil, err
		}
	}

	//
	return &Outfile{OldName: UpImg.Filename, URL: u.Url + _tPath, Size: UpImg.Size}, nil

}

// checkFile 检测上传文件
func checkFile(file *ghttp.UploadFile, Exts string, MaxSize int64) error {

	size := file.Size

	ext := path.Ext(file.FileHeader.Filename)[1:]

	if !inArray(Exts, ext) {
		return ErrLang("Upload_Ext_Is_Err", Exts)

	}
	// checkSize
	if size > MaxSize {
		return ErrLang("Upload_Size_Is_Err", MaxSize/1024/1024)

	}
	return nil
}

// inArray 类型检测
func inArray(arrStr, str string) bool {
	Arr := strings.Split(arrStr, "|")

	return InArray(str, Arr)
}
