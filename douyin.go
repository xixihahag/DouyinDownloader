package main

import (
	"douyin/internal/model"
	"douyin/pkg/util/httputil"
	"douyin/pkg/util/logutil"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const (
	regSharedUrl = `http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`
	realUrl      = `https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=`
)

var Logger *zap.SugaredLogger

type Douyin struct {
	reg      *regexp.Regexp
	videoReg *regexp.Regexp
	userReg  *regexp.Regexp
	content  *fyne.Container
	scroll   *container.Scroll
}

func NewDouyin(content *fyne.Container, scroll *container.Scroll) *Douyin {
	z, _ := zap.NewDevelopment()
	Logger = z.Sugar()
	return &Douyin{
		reg:      regexp.MustCompile(regSharedUrl),
		videoReg: regexp.MustCompile("https://www.iesdouyin.com/share/video/([0-9]+)"),
		userReg:  regexp.MustCompile("https://www.iesdouyin.com/share/user/(\\w+)"),
		content:  content,
		scroll:   scroll,
	}
}

func (d *Douyin) Get(sharedContent string) error {
	d.addLogAndScroll("Start getting data, please wait ...")
	url := d.reg.FindString(sharedContent)
	if url == "" {
		d.addLogAndScroll("分享链接不正确, 请输入例如：https://v.douyin.com/6oGRLEE/")
		return fmt.Errorf("分享链接不正确, content: %s", sharedContent)
	}

	ids, err := d.getIDs(url)
	if err != nil {
		return err
	}

	count := len(ids)
	d.addLogAndScroll(fmt.Sprintf("start download, count %d", count))
	for index, id := range ids {
		d.addLogAndScroll(fmt.Sprintf("deal %d/%d ...", index+1, count))
		info, err := d.getVideoInfo(realUrl + id)
		if err != nil {
			continue
		}
		v, err := d.parseInfo(info)
		if err != nil {
			continue
		}
		if err := d.Download(&v); err != nil {
			continue
		}
		d.addLogAndScroll(fmt.Sprintf("Download Success! %s", v.Desc))
	}
	d.addLogAndScroll("all data download done")
	return nil
}

// 根据 url 的返回结果判断
func (d *Douyin) getIDs(url string) ([]string, error) {
	logutil.Infof("getIDs url: %s", url)
	// 根据返回body 判断是单个视频/video 还是 抖音首页/user
	body, err := httputil.Get(url)
	if err != nil {
		return nil, err
	}
	// logutil.Infof("body: %s", body)

	exp, err := regexp.Compile("\\d+")
	if err != nil {
		return nil, err
	}
	result := exp.FindString(string(body))
	if result == "" {
		return nil, fmt.Errorf("解析参数失败 ->" + string(body))
	}

	var ids []string
	if id := d.videoReg.FindStringSubmatch(body); id != nil {
		// 找到的是一个 video
		logutil.Info("single video")
		logutil.Infof("id: %s", id[1])
		ids = append(ids, string(id[1]))
	} else if uid := d.userReg.FindStringSubmatch(body); uid != nil {
		// 找到的是首页
		logutil.Info("first page")
		logutil.Infof("id: %s", uid[1])
		i, err := getIDsFromHomePage(uid[1])
		if err != nil {
			return nil, err
		}
		ids = append(ids, i...)
	}
	return ids, nil
}

func getIDsFromHomePage(secUid string) ([]string, error) {
	init := true
	cursor := 0
	count := 30
	ids := []string{}

	for {
		if cursor == 0 && !init {
			break
		}

		apiURL := fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/aweme/post/?sec_uid=%s&count=%d&max_cursor=%d&aid=1128&_signature=PDHVOQAAXMfFyj02QEpGaDwx1S&dytk=",
			secUid, count, cursor)
		body, err := httputil.Get(apiURL)
		if err != nil {
			return nil, err
		}

		var userPost model.UserPost
		if err := json.Unmarshal([]byte(body), &userPost); err != nil {
			logutil.Errorf("unmarshal fail, body: %s, err: %v", body, err)
			return nil, err
		}
		if userPost.StatusCode != 0 {
			msg := fmt.Sprintf("return code != 0, rsp: %v", body)
			logutil.Error(msg)
			return nil, errors.New(msg)
		}

		init = false
		cursor = int(userPost.MaxCursor)
		for _, item := range userPost.AwemeList {
			ids = append(ids, item.AwemeID)
		}
	}
	return ids, nil
}

func (d *Douyin) getVideoInfo(url string) (string, error) {
	return httputil.Get(url)
}

func (d *Douyin) parseInfo(info string) (Video, error) {
	item := gjson.Get(info, "item_list.0")
	res := item.Get("video.play_addr.url_list.0")
	video := Video{}
	if !res.Exists() {
		return video, fmt.Errorf("未找到视频地址, info: %s", info)
	}
	// 获取播放地址
	video.PlayAddr = strings.ReplaceAll(res.Str, "playwm", "play")

	// 获取创建时间
	ct := item.Get("create_time").Int()
	video.CreateTime = time.Unix(ct, 0).Format("2006-01-02 15.04.05")

	// 获取播放时长
	res = item.Get("duration")
	// 视频类的有播放时长，图文类无播放时长
	if res.Exists() && res.Raw != "0" {
		video.VideoType = VideoPlay
	} else {
		video.VideoType = Imageplay
		res = item.Get("images")
		if res.Exists() && res.IsArray() {
			for _, image := range res.Array() {
				imageRes := image.Get("url_list.0")
				if imageRes.Exists() {
					video.Images = append(video.Images, ImageItem{
						ImageUrl: imageRes.Str,
						ImageID:  image.Get("uri").Str,
					})
				}
			}
		}
	}

	// 获取播放地址
	res = item.Get("video.play_addr.uri")
	if res.Exists() {
		video.PlayID = res.Str
	}

	// 获取视频唯一 ID
	res = item.Get("aweme_id")
	if res.Exists() {
		video.VideoID = res.Str
	}

	//获取封面
	res = item.Get("video.cover.url_list.0")
	if res.Exists() {
		video.Cover = res.Str
	}

	// 获取原始封面
	res = item.Get("video.origin_cover.url_list.0")
	if res.Exists() {
		video.OriCover = res.Str
	}
	res = item.Get("video.origin_cover.url_list")
	if res.Exists() {
		res.ForEach(func(key, value gjson.Result) bool {
			video.OriCoverList = append(video.OriCoverList, value.Str)
			return true
		})
	}

	//获取作者信息
	res = item.Get("author.uid")
	if res.Exists() {
		video.Author.Id = res.Str
	}
	res = item.Get("author.short_id")
	if res.Exists() {
		video.Author.ShortId = res.Str
	}
	res = item.Get("author.nickname")
	if res.Exists() {
		video.Author.Nickname = res.Str
	}
	res = item.Get("author.signature")
	if res.Exists() {
		video.Author.Signature = res.Str
	}
	//回获取作者大头像
	res = item.Get("author.avatar_larger.url_list.0")
	if res.Exists() {
		video.Author.AvatarLarger = res.Str
	}

	//获取视频描述
	res = item.Get("desc")
	if res.Exists() {
		video.Desc = res.Str
	}
	return video, nil
}

func (d *Douyin) Download(v *Video) error {
	return v.Download()
}

func (d *Douyin) addLogAndScroll(text string) {
	d.content.Add(
		widget.NewLabel(text),
	)
	d.scroll.ScrollToBottom()
}
