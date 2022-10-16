package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/downloader"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

var OutputDir = "output"

type DeDaoDownloader interface {
	Download() error
}

type CourseDownload struct {
	DownloadType int // 1:mp3, 2:PDF文档, 3:markdown文档
	ID           int
	AID          int
}

type OdobDownload struct {
	DownloadType int // 1:mp3, 2:PDF文档, 3:markdown文档
	ID           int
}

type EBookDownload struct {
	DownloadType int // 1:html, 2:PDF文档, 3:epub
	ID           int
}

func (d *CourseDownload) Download() error {
	course, err := CourseInfo(d.ID)
	if err != nil {
		return err
	}
	articles, err := ArticleList(d.ID, "")
	if err != nil {
		return err
	}

	switch d.DownloadType {
	case 1: // mp3
		downloadData := extractDownloadData(course, articles, d.AID, 1)
		errors := make([]error, 0)

		path, err := utils.Mkdir(OutputDir, utils.FileName(course.ClassInfo.Name, ""), "MP3")
		if err != nil {
			return err
		}

		for _, datum := range downloadData.Data {
			if !datum.IsCanDL {
				continue
			}
			stream := datum.Enid
			if err := downloader.Download(datum, stream, path); err != nil {
				errors = append(errors, err)
			}

		}
		if len(errors) > 0 {
			return errors[0]
		}
	case 2:
		// 下载 PDF
		downloadData := extractDownloadData(course, articles, d.AID, 2)
		errors := make([]error, 0)

		path, err := utils.Mkdir(OutputDir, utils.FileName(course.ClassInfo.Name, ""), "PDF")
		if err != nil {
			return err
		}

		cookies := LoginedCookies()
		for _, datum := range downloadData.Data {
			if err := downloader.PrintToPDF(datum, cookies, path); err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) > 0 {
			return errors[0]
		}
	case 3:
		// 下载 Markdown
		path, err := utils.Mkdir(OutputDir, utils.FileName(course.ClassInfo.Name, ""), "MD")
		if err != nil {
			return err
		}
		if err := DownloadMarkdown(CateCourse, d.ID, d.AID, path); err != nil {
			return err
		}
	}
	return nil

}

func (d *OdobDownload) Download() error {
	fileName := "每天听本书"
	switch d.DownloadType {
	case 1:
		downloadData := downloader.Data{
			Title: fileName,
		}
		downloadData.Type = "audio"
		downloadData.Data = extractOdobDownloadData(d.ID)
		errors := make([]error, 0)
		path, err := utils.Mkdir(OutputDir, utils.FileName(fileName, ""), "MP3")
		if err != nil {
			return err
		}
		for _, datum := range downloadData.Data {
			if !datum.IsCanDL {
				continue
			}
			stream := datum.Enid
			if err := downloader.Download(datum, stream, path); err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) > 0 {
			return errors[0]
		}
	case 2:
		err := errors.New("得到 Web 端暂未开放每天听本书，PDF 无法下载。")
		return err
	case 3:
		// 下载 Markdown
		path, err := utils.Mkdir(OutputDir, utils.FileName(fileName, ""), "MD")
		if err != nil {
			return err
		}
		if err := DownloadMarkdown(CateAudioBook, d.ID, 0, path); err != nil {
			return err
		}
	}
	return nil
}

func (d *EBookDownload) Download() error {
	detail, err := EbookDetail(d.ID)
	if err != nil {
		return err
	}

	title := strconv.Itoa(d.ID) + "_"
	if detail.Title != "" {
		title += detail.Title
	} else if detail.OperatingTitle != "" {
		title += detail.OperatingTitle
	}

	title += "_" + detail.BookAuthor
	info, svgContent, err := EbookPage(detail.Enid)
	if err != nil {
		return err
	}
	sort.Sort(svgContent)

	switch d.DownloadType {
	case 1:
		var toc []*utils.EbookToc
		for _, ebookToc := range info.BookInfo.Toc {
			toc = append(toc, &utils.EbookToc{
				Href:      ebookToc.Href,
				Level:     ebookToc.Level,
				PlayOrder: ebookToc.PlayOrder,
				Offset:    ebookToc.Offset,
				Text:      ebookToc.Text,
			})
		}
		if err = utils.Svg2Html(title, svgContent, toc); err != nil {
			return err
		}

	case 2:
		if err = utils.Svg2Pdf(title, svgContent); err != nil {
			return err
		}

	case 3:
		var opts utils.EpubOptions
		opts.Title = title
		opts.Author = detail.BookAuthor
		opts.Description = detail.BookIntro

		if err = utils.Svg2Epub(title, svgContent, opts); err != nil {
			return err
		}

		return err
	}

	return nil
}

func Download(downloader DeDaoDownloader) error {
	return downloader.Download()
}

// 生成下载数据
func extractDownloadData(course *services.CourseInfo, articles *services.ArticleList, aid int, flag int) downloader.Data {

	downloadData := downloader.Data{
		Title: course.ClassInfo.Name,
	}

	if course.HasAudio() {
		downloadData.Type = "audio"
		downloadData.Data = extractCourseDownloadData(articles, aid, flag)
	}

	return downloadData
}

// 生成课程下载数据
func extractCourseDownloadData(articles *services.ArticleList, aid int, flag int) []downloader.Datum {
	data := downloader.EmptyData
	audioIds := map[int]string{}

	audioData := make([]*downloader.Datum, 0)
	for _, article := range articles.List {
		if aid > 0 && article.ID != aid {
			continue
		}

		if article.VideoStatus == 0 && len(article.AudioAliasIds) > 0 {
			audioIds[article.ID] = article.Audio.AliasID

			var urls []downloader.URL
			key := article.Enid
			streams := map[string]downloader.Stream{
				key: {
					URLs:    urls,
					Size:    article.Audio.Size,
					Quality: key,
				},
			}
			isCanDL := true
			if len(article.Audio.AliasID) == 0 {
				isCanDL = false
			}
			datum := &downloader.Datum{
				ID:        article.ID,
				Enid:      article.Enid,
				ClassEnid: article.ClassEnid,
				ClassID:   article.ClassID,
				Title:     article.Title,
				IsCanDL:   isCanDL,
				M3U8URL:   article.Audio.Mp3PlayURL,
				Streams:   streams,
				Type:      "audio",
			}

			audioData = append(audioData, datum)
		}

	}

	if flag == 1 {
		handleStreams(audioData, audioIds)
	}

	for _, d := range audioData {
		data = append(data, *d)
	}
	return data
}

// 生成 AudioBook 下载数据
func extractOdobDownloadData(aid int) []downloader.Datum {
	data := downloader.EmptyData
	audioIds := map[int]string{}

	audioData := make([]*downloader.Datum, 0)
	article := config.Instance.GetIDMap(CateAudioBook, aid)
	aliasID := article["audio_alias_id"].(string)
	if aliasID == "" {
		list, err := CourseList(CateAudioBook)
		if err != nil {
			return nil
		}
		for _, course := range list.List {
			if aid > 0 && course.ID == aid {
				article = GetCourseIDMap(&course)
				break
			}
		}
	}

	audioIds[aid] = article["audio_alias_id"].(string)

	var urls []downloader.URL
	key := article["enid"].(string)
	streams := map[string]downloader.Stream{
		key: {
			URLs:    urls,
			Size:    int(article["audio_size"].(float64)),
			Quality: key,
		},
	}
	isCanDL := true
	if !article["has_play_auth"].(bool) {
		isCanDL = false
	}
	datum := &downloader.Datum{
		ID:      aid,
		Enid:    article["enid"].(string),
		ClassID: int(article["class_id"].(float64)),
		Title:   article["title"].(string),
		IsCanDL: isCanDL,
		M3U8URL: article["audio_mp3_play_url"].(string),
		Streams: streams,
		Type:    "audio",
	}

	audioData = append(audioData, datum)
	handleStreams(audioData, audioIds)

	for _, d := range audioData {
		data = append(data, *d)
	}
	return data
}

func handleStreams(audioData []*downloader.Datum, audioIds map[int]string) {
	wgp := utils.NewWaitGroupPool(10)
	for _, datum := range audioData {
		wgp.Add()
		go func(datum *downloader.Datum, streams map[int]string) {
			defer func() {
				wgp.Done()
			}()
			if datum.IsCanDL {
				if urls, err := utils.M3u8URLs(datum.M3U8URL); err == nil {
					key := datum.Enid
					stream := datum.Streams[key]
					for _, url := range urls {
						stream.URLs = append(stream.URLs, downloader.URL{
							URL: url,
							Ext: "ts",
						})
					}
					datum.Streams[key] = stream
				}
				for k, v := range datum.Streams {
					if len(v.URLs) == 0 {
						delete(datum.Streams, k)
					}
				}
			}
		}(datum, audioIds)
	}
	wgp.Wait()
}

func ContentsToMarkdown(contents []services.Content) (res string) {
	for _, content := range contents {
		switch content.Type {
		case "audio":
			title := strings.TrimRight(content.Title, ".mp3")
			res += getMdHeader(1) + title + "\r\n\r\n"
		case "header":
			content.Text = strings.Trim(content.Text, " ")
			if len(content.Text) > 0 {
				res += getMdHeader(content.Level) + content.Text + "\r\n\r\n"
			}
		case "blockquote":
			texts := strings.Split(content.Text, "\n")
			for _, text := range texts {
				res += "> " + text + "\r\n"
				res += "> \r\n"
			}
			res = strings.TrimRight(res, "> \r\n")
			res += "\r\n\r\n"
		case "paragraph":
			// map 转结构体
			tmpJson, err := jsoniter.Marshal(content.Contents)
			if err != nil {
				return
			}
			cont := services.Contents{}
			err = jsoniter.Unmarshal(tmpJson, &cont)
			if err != nil {
				return ""
			}
			for _, item := range cont {
				subContent := strings.Trim(item.Text.Content, " ")
				switch item.Type {
				case "text":
					if item.Text.Bold {
						res += " **" + subContent + "** "
					} else if item.Text.Highlight {
						res += " *" + subContent + "* "
					} else {
						res += subContent
					}
				}
			}
			res = strings.Trim(res, " ")
			res = strings.Trim(res, "\r\n")
			res += "\r\n\r\n"
		case "list":
			tmpJson, err := jsoniter.Marshal(content.Contents)
			if err != nil {
				return
			}
			var cont []services.Contents
			err = jsoniter.Unmarshal(tmpJson, &cont)
			if err != nil {
				return ""
			}

			for _, item := range cont {
				for _, item := range item {
					subContent := strings.Trim(item.Text.Content, " ")
					switch item.Type {
					case "text":
						if item.Text.Bold {
							res += "* **" + subContent + "** "
						} else if item.Text.Highlight {
							res += "* *" + subContent + "* "
						} else {
							res += "* " + subContent
						}
					}
				}
				res += "\r\n\r\n"
			}
		case "elite": // 划重点
			res += getMdHeader(2) + "划重点\r\n\r\n" + content.Text + "\r\n\r\n"

		case "image":
			res += "![" + content.URL + "](" + content.URL + ")" + "\r\n\r\n"
		case "label-group":
			res += getMdHeader(2) + "`" + content.Text + "`" + "\r\n\r\n"
		}
	}

	res += "---\r\n"
	return
}

func articleCommentsToMarkdown(contents []services.ArticleComment) (res string) {
	res = getMdHeader(2) + "热门留言\r\n\r\n"
	for _, content := range contents {
		res += content.NotesOwner.Name + "：" + content.Note + "\r\n\r\n"
		if content.CommentReply != "" {
			res += "> " + content.CommentReplyUser.Name + "(" + content.CommentReplyUser.Role + ") 回复：" + content.CommentReply + "\r\n\r\n"
		}
	}
	res += "---\r\n"
	return
}

func getMdHeader(level int) string {
	heads := map[int]string{
		1: "# ",
		2: "## ",
		3: "### ",
		4: "#### ",
		5: "##### ",
		6: "###### ",
	}
	if s, ok := heads[level]; ok {
		return s
	}
	return ""
}

func DownloadMarkdown(cType string, id, aid int, path string) error {
	switch cType {
	case CateCourse:
		list, err := ArticleList(id, "")
		if err != nil {
			return err
		}
		for _, v := range list.List {
			if aid > 0 && v.ID != aid {
				continue
			}
			detail, enId, err := ArticleDetail(id, v.ID)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			var content []services.Content
			err = jsoniter.UnmarshalFromString(detail.Content, &content)
			if err != nil {
				return err
			}

			name := utils.FileName(v.Title, "md")
			fileName := filepath.Join(path, name)
			fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", name)
			_, exist, err := utils.FileSize(fileName)

			if err != nil {
				fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
				return err
			}

			if exist {
				fmt.Printf("\033[33;1m%s\033[0m\n", "已存在")
				return nil
			}

			res := ContentsToMarkdown(content)
			// 添加留言
			commentList, err := ArticleCommentList(enId, "like", 1, 20)
			if err == nil {
				res += articleCommentsToMarkdown(commentList.List)
			}

			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
				return err
			}
			_, err = f.WriteString(res)
			if err != nil {
				fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
				return err
			}
			if err = f.Close(); err != nil {
				if err != nil {
					return err
				}
			}
			fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
		}
	case CateAudioBook:
		info := config.Instance.GetIDMap(CateAudioBook, id)
		aliasID := info["audio_alias_id"].(string)
		if aliasID == "" {
			list, err := CourseList(cType)
			if err != nil {
				return err
			}
			for _, v := range list.List {
				if v.AudioDetail.SourceID == id {
					aliasID = v.AudioDetail.AliasID
					break
				}
			}
		}
		detail, err := OdobArticleDetail(aliasID)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		var content []services.Content
		err = jsoniter.UnmarshalFromString(detail.Content, &content)
		if err != nil {
			return err
		}

		name := utils.FileName(info["title"].(string), "md")
		fileName := filepath.Join(path, name)
		fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", name)
		_, exist, err := utils.FileSize(fileName)

		if err != nil {
			fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
			return err
		}

		if exist {
			fmt.Printf("\033[33;1m%s\033[0m\n", "已存在")
			return nil
		}

		res := ContentsToMarkdown(content)

		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
			return err
		}
		_, err = f.WriteString(res)
		if err != nil {
			fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
			return err
		}
		if err = f.Close(); err != nil {
			if err != nil {
				return err
			}
		}
		fmt.Printf("\033[32;1m%s\033[0m\n", "完成")

	case CateAce:

	}

	return nil

}
