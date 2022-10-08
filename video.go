package main

import (
	"douyin/pkg/util/httputil"
	"douyin/pkg/util/logutil"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type VideoType int

const (
	//VideoPlayType 视频类
	VideoPlay VideoType = 0
	//ImagePlayType 图文类
	Imageplay VideoType = 1
)

type Video struct {
	PlayAddr     string
	VideoType    VideoType
	Images       []ImageItem
	PlayID       string
	VideoID      string
	Cover        string   // 封面
	OriCover     string   // 原始封面
	OriCoverList []string // 原始封面列表
	MusicAddr    string   // 音乐地址
	Author       struct {
		Id           string `json:"id"`
		ShortId      string `json:"short_id"`
		Nickname     string `json:"nickname"`      // 昵称
		AvatarLarger string `json:"avatar_larger"` // 大头像
		Signature    string `json:"signature"`     // 签名
	} `json:"author"`
	Desc       string // 视频描述
	CreateTime string // 创建时间
}

// 图片信息
type ImageItem struct {
	ImageUrl string
	ImageID  string
}

func (v *Video) Download() error {
	defer func() {
		if err := recover(); err != nil {
			logutil.Error("panic")
		}
	}()
	defer func() {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}()
	dir, err := v.getFilePath()
	if err != nil {
		return err
	}
	if v.VideoType == Imageplay {
		// 拉图片
		// 新建一个文件夹放图片
		dir = filepath.Join(dir, v.CreateTime+v.Desc)
		_, err := os.Stat(dir)
		if !os.IsNotExist(err) {
			return nil
		}
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0655); err != nil {
				return err
			}
		}

		for _, image := range v.Images {
			ext := ".jepg"
			uri, err := url.Parse(image.ImageUrl)
			if err != nil {
				return err
			}
			ext = filepath.Ext(uri.Path)
			imageID := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(image.ImageID, "//", ""), "\\\\", "/"), "/", "-")
			imageName := filepath.Join(dir, imageID+ext)

			body, err := httputil.Get(image.ImageUrl)
			if err != nil {
				continue
			}

			if err := ioutil.WriteFile(imageName, []byte(body), 0655); err != nil {
				Logger.Errorf("save iamge fail, err: %v", err)
				continue
			}
		}
		return nil
	}
	// 存视频
	body, err := httputil.Get(v.PlayAddr)
	if err != nil {
		return err
	}

	var videoName string
	if ext := filepath.Ext(v.PlayID); ext != "" {
		videoName = v.CreateTime + v.Desc + ext
	} else {
		videoName = v.CreateTime + v.Desc + ".mp4"
	}
	videoName = filepath.Join(dir, videoName)

	// 检查这个文件是否存在（是否已经被下载过）
	if _, err := os.Stat(videoName); !os.IsNotExist(err) {
		return nil
	}

	// Logger.Infof("videoName: %v", videoName)
	f, err := os.Create(videoName)
	if err != nil {
		return fmt.Errorf("create path fail, path: %v, err: %v", videoName, err)
	}
	defer f.Close()
	if _, err := io.Copy(f, strings.NewReader(body)); err != nil {
		return fmt.Errorf("write file fail, err: %v", err)
	}

	return nil
}

func (v *Video) getFilePath() (string, error) {
	savePath, _ := os.Getwd()
	filepath.Join(savePath, "Download")
	filePath, err := filepath.Abs(savePath)
	if err != nil {
		log.Printf("获取文件路径失败,err: %v", err)
		return "", err
	}
	dir := filepath.Join(filePath, v.Author.Nickname)

	// 创建文件夹
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0655); err != nil {
			return "", nil
		}
	}
	return dir, nil
}
