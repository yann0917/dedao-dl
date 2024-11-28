package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
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
	IsMerge      bool
	IsComment    bool
	IsOrder      bool
	ClassName    string
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

	switch d.DownloadType {
	case 1: // mp3
		articles, err := ArticleList(d.ID, "")
		if err != nil {
			return err
		}
		downloadData := extractDownloadData(course, articles, d.AID, 1, d.IsOrder)
		errs := make([]error, 0)

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
				errs = append(errs, err)
			}

		}
		if len(errs) > 0 {
			return errs[0]
		}
	case 2:
		// 下载 PDF
		// downloadData := extractDownloadData(course, articles, d.AID, 2)
		errs := make([]error, 0)

		path, err := utils.Mkdir(OutputDir, utils.FileName(course.ClassInfo.Name, ""), "PDF")
		if err != nil {
			return err
		}
		d.ClassName = course.ClassInfo.Name
		if err := DownloadPdfCourse(d, path); err != nil {
			return err
		}
		// cookies := LoginedCookies()
		// for _, datum := range downloadData.Data {
		// 	if err := downloader.PrintToPDF(datum, cookies, path); err != nil {
		// 		errs = append(errs, err)
		// 	}
		// }
		if len(errs) > 0 {
			return errs[0]
		}
	case 3:
		// 下载 Markdown
		path, err := utils.Mkdir(OutputDir, utils.FileName(course.ClassInfo.Name, ""), "MD")
		if err != nil {
			return err
		}
		d.ClassName = course.ClassInfo.Name
		if err := DownloadMarkdownCourse(d, path); err != nil {
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
		errs := make([]error, 0)
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
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return errs[0]
		}
	case 2:
		path, err := utils.Mkdir(OutputDir, utils.FileName(fileName, ""), "PDF")
		if err != nil {
			return err
		}
		info, content, err2 := getArticleDetail(d.ID)
		if err2 != nil {
			return err2
		}
		res := ContentsToMarkdown(content)
		return utils.Md2Pdf(path, info["title"].(string), []byte(res))

	case 3:
		// 下载 Markdown
		path, err := utils.Mkdir(OutputDir, utils.FileName(fileName, ""), "MD")
		if err != nil {
			return err
		}
		if err := DownloadMarkdownAudioBook(d.ID, path); err != nil {
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
		if err = utils.Svg2Html(title, svgContent, info.BookInfo.Toc); err != nil {
			return err
		}

	case 2:
		if err = utils.Svg2Pdf(title, svgContent, info.BookInfo.Toc); err != nil {
			return err
		}

	case 3:
		var opts utils.EpubOptions
		opts.Title = title
		opts.Author = detail.BookAuthor
		opts.Description = detail.BookIntro
		opts.Toc = info.BookInfo.Toc

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
func extractDownloadData(course *services.CourseInfo, articles *services.ArticleList, aid int, flag int, isOrder bool) downloader.Data {

	downloadData := downloader.Data{
		Title: course.ClassInfo.Name,
	}

	if course.HasAudio() {
		downloadData.Type = "audio"
		downloadData.Data = extractCourseDownloadData(articles, aid, flag, isOrder)
	}

	return downloadData
}

// 生成课程下载数据
func extractCourseDownloadData(articles *services.ArticleList, aid int, flag int, isOrder bool) []downloader.Datum {
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
			name := article.Title
			if isOrder {
				name = fmt.Sprintf("%03d.%s", article.OrderNum, name)
			}
			datum := &downloader.Datum{
				ID:        article.ID,
				Enid:      article.Enid,
				ClassEnid: article.ClassEnid,
				ClassID:   article.ClassID,
				Title:     name,
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
	detail, err := getService().AudioDetailAlias(aliasID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	datum := &downloader.Datum{
		ID:      aid,
		Enid:    article["enid"].(string),
		ClassID: int(article["class_id"].(float64)),
		Title:   article["title"].(string),
		IsCanDL: isCanDL,
		M3U8URL: detail.Mp3PlayURL,
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
			resP, err := paragraphToMarkDown(content.Contents)
			if err != nil {
				return
			}
			res += resP
		case "list":
			resL, err := listToMarkdown(content.Contents)
			if err != nil {
				return
			}
			res += resL
		case "elite": // 划重点
			res += getMdHeader(2) + "划重点\r\n\r\n" + content.Text + "\r\n\r\n"

		case "image":
			res += "![" + content.URL + "](" + content.URL
			if content.Legend != "" {
				res += " \"" + content.Legend + "\""
			}
			res += ")" + "\r\n\r\n"
		case "label-group":
			res += getMdHeader(2) + "`" + content.Text + "`" + "\r\n\r\n"
		}
	}

	res += "---\r\n"
	return
}

func paragraphToMarkDown(content interface{}) (res string, err error) {
	tmpJson, err := jsoniter.Marshal(content)
	if err != nil {
		return
	}
	cont := services.Contents{}
	err = jsoniter.Unmarshal(tmpJson, &cont)
	if err != nil {
		return
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
	return
}

func listToMarkdown(content interface{}) (res string, err error) {
	tmpJson, err := jsoniter.Marshal(content)
	if err != nil {
		return
	}
	var cont []services.Contents
	err = jsoniter.Unmarshal(tmpJson, &cont)
	if err != nil {
		return
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

func DownloadMarkdownCourse(d *CourseDownload, path string) error {
	list, err := ArticleList(d.ID, "")
	if err != nil {
		return err
	}
	name, fileName := "", ""
	mName, mFileName := "", ""
	if d.IsMerge {
		mName = utils.FileName(d.ClassName+"-合集", "md")
		mFileName = filepath.Join(path, mName)
		fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】\n", mFileName)
	}
	for _, v := range list.List {
		if d.AID > 0 && v.ID != d.AID {
			continue
		}
		detail, enId, err := ArticleDetail(d.ID, v.ID)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		// fmt.Printf("%#v\n", detail)

		var content []services.Content
		err = jsoniter.UnmarshalFromString(detail.Content, &content)
		if err != nil {
			return err
		}

		name = utils.FileName(v.Title, "md")
		if d.IsOrder {
			name = fmt.Sprintf("%03d.%s", v.OrderNum, name)
		}
		fileName = filepath.Join(path, name)
		fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", name)
		_, exist, err := utils.FileSize(fileName)

		if err != nil {
			fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
			return err
		}

		if exist {
			fmt.Printf("\033[33;1m%s\033[0m\n", "已存在")
			continue
		}

		res := ContentsToMarkdown(content)
		if d.IsComment {
			// 添加留言
			commentList, err := ArticleCommentList(enId, "like", 1, 20)
			if err == nil {
				res += articleCommentsToMarkdown(commentList.List)
			}
		}

		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
			return err
		}
		fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
		if d.IsMerge {
			f, err := os.OpenFile(mFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Printf("\033[31;1m%s\033[0m\n", "合集失败"+err.Error())
				return err
			}
			_, err = f.WriteString(res)
			if err != nil {
				fmt.Printf("\033[31;1m%s\033[0m\n", "合集失败"+err.Error())
				return err
			}
			if err = f.Close(); err != nil {
				return err
			}
		}
	}
	if d.IsMerge {
		fmt.Printf("\033[32;1m%s\033[0m\n", "合集完成")
	}
	return nil
}

func DownloadPdfCourse(d *CourseDownload, path string) error {
	list, err := ArticleList(d.ID, "")
	if err != nil {
		return err
	}
	name, fileName := "", ""
	// mName, mFileName := "", ""
	// if d.IsMerge {
	// 	mName = utils.FileName(d.ClassName+"-合集", "pdf")
	// 	mFileName = filepath.Join(path, mName)
	// 	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】\n", mFileName)
	// }
	for _, v := range list.List {
		if d.AID > 0 && v.ID != d.AID {
			continue
		}
		detail, enId, err := ArticleDetail(d.ID, v.ID)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		// fmt.Printf("%#v\n", detail)

		var content []services.Content
		err = jsoniter.UnmarshalFromString(detail.Content, &content)
		if err != nil {
			return err
		}

		name = utils.FileName(v.Title, "pdf")
		if d.IsOrder {
			name = fmt.Sprintf("%03d.%s", v.OrderNum, name)
		}
		fileName = filepath.Join(path, name)
		fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", name)
		_, exist, err := utils.FileSize(fileName)

		if err != nil {
			fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
			return err
		}

		if exist {
			fmt.Printf("\033[33;1m%s\033[0m\n", "已存在")
			continue
		}

		res := ContentsToMarkdown(content)
		if d.IsComment {
			// 添加留言
			commentList, err := ArticleCommentList(enId, "like", 1, 20)
			if err == nil {
				res += articleCommentsToMarkdown(commentList.List)
			}
		}
		err = utils.Md2Pdf(path, name, []byte(res))
		if err != nil {
			return err
		}
		// 	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		// 	if err != nil {
		// 		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
		// 		return err
		// 	}
		// 	_, err = f.WriteString(res)
		// 	if err != nil {
		// 		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
		// 		return err
		// 	}
		// 	if err = f.Close(); err != nil {
		// 		return err
		// 	}
		// 	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
		// 	if d.IsMerge {
		// 		f, err := os.OpenFile(mFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		// 		if err != nil {
		// 			fmt.Printf("\033[31;1m%s\033[0m\n", "合集失败"+err.Error())
		// 			return err
		// 		}
		// 		_, err = f.WriteString(res)
		// 		if err != nil {
		// 			fmt.Printf("\033[31;1m%s\033[0m\n", "合集失败"+err.Error())
		// 			return err
		// 		}
		// 		if err = f.Close(); err != nil {
		// 			return err
		// 		}
		// 	}
		// }
		// if d.IsMerge {
		// 	fmt.Printf("\033[32;1m%s\033[0m\n", "合集完成")
	}
	return nil
}

func DownloadMarkdownAudioBook(id int, path string) error {
	info, content, err2 := getArticleDetail(id)
	if err2 != nil {
		return err2
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
		return err
	}
	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
	return nil
}

func getArticleDetail(id int) (info map[string]interface{}, content []services.Content, err error) {
	info = config.Instance.GetIDMap(CateAudioBook, id)
	aliasID := info["audio_alias_id"].(string)
	if aliasID == "" {
		list, err1 := CourseList(CateAudioBook)
		if err1 != nil {
			return nil, nil, err1
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
		return nil, nil, err
	}

	err = jsoniter.UnmarshalFromString(detail.Content, &content)
	if err != nil {
		return
	}
	return
}
