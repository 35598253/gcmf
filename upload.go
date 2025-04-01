package gcmf

import (
	"context"
	"fmt"
	"github.com/35598253/gcmf/image"

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
type FileUpCfg struct {
	Exts    string `json:"exts"`
	MaxSize int64  `json:"maxSize"`
}
type ImageUpCfg struct {
	Exts        string `json:"exts"`
	ImageResize uint   `json:"imageResize"`
	FontPath    string `json:"fontPath"`
	ImageWater  string `json:"imageWater"`
}
type UploadConfig struct {
	RealPath    string      `json:"realPath"`
	Url         string      `json:"url"`
	UploadDir   string      `json:"uploadDir"`
	FileConfig  *FileUpCfg  `json:"fileConfig"`
	ImageConfig *ImageUpCfg `json:"imageConfig"`
}

type Uploader interface {
	File(UpFile *ghttp.UploadFile) (*Outfile, error)
	Image(UpImg *ghttp.UploadFile) (*Outfile, error)
}

// Upload file
func Upload(Ctx context.Context, UpCfg ...*UploadConfig) Uploader {
	cfg, _ := GetConfig(Ctx)

	_cfgUp := cfg.Upload

	if len(UpCfg) > 0 {
		dfCfg := UpCfg[0]

		if dfCfg.ImageConfig != nil {
			if dfCfg.ImageConfig.Exts != "" {
				_cfgUp.ImageConfig.Exts = dfCfg.ImageConfig.Exts
			}
			if dfCfg.ImageConfig.ImageResize > 0 {
				_cfgUp.ImageConfig.ImageResize = dfCfg.ImageConfig.ImageResize
			}
			if dfCfg.ImageConfig.FontPath != "" {
				_cfgUp.ImageConfig.FontPath = dfCfg.ImageConfig.FontPath
			}
			if dfCfg.ImageConfig.ImageWater != "" {
				_cfgUp.ImageConfig.ImageWater = dfCfg.ImageConfig.ImageWater
			}
		}
		if dfCfg.FileConfig != nil {
			if dfCfg.FileConfig.Exts != "" {
				_cfgUp.FileConfig.Exts = dfCfg.FileConfig.Exts
			}
			if dfCfg.FileConfig.MaxSize > 0 {
				_cfgUp.FileConfig.MaxSize = dfCfg.FileConfig.MaxSize
			}
		}
	}

	return _cfgUp
}

// File 上传文件
func (u *UploadConfig) File(UpFile *ghttp.UploadFile) (*Outfile, error) {
	var tPath = fmt.Sprintf("/file/%s/%s/", u.UploadDir, time.Now().Format("2006"))
	var upPath = u.RealPath + tPath
	var outPath = u.Url + tPath

	// 读取默认

	if err := checkFile(UpFile, u.FileConfig.Exts, u.FileConfig.MaxSize); err != nil {
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
func (u *UploadConfig) Image(UpImg *ghttp.UploadFile) (*Outfile, error) {

	imgOutPath := fmt.Sprintf("/images/%s/%s/", u.UploadDir, time.Now().Format("2006"))
	// 剪裁&水印
	imgTmp, err := image.NewImageUpload(UpImg, u.ImageConfig.ImageResize)
	if err != nil {
		return nil, err
	}
	var _tPath string
	if u.ImageConfig.ImageWater != "" {

		_tImg, err := imgTmp.TextWater(u.ImageConfig.ImageWater, u.ImageConfig.FontPath)
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
