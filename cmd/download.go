package cmd

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/downloader"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// OutputDir OutputDir
var OutputDir = "output"

var downloadType int
var downloadCmd = &cobra.Command{
	Use:   "dl",
	Short: "下载已购买课程，并转换成 PDF & 音频",
	Long: `使用 dedao-dl dl 下载已购买课程, 并转换成 PDF & 音频 & markdown
-t 指定下载格式, 1:mp3, 2:PDF文档, 3:markdown文档, 默认 mp3`,
	Example: "dedao-dl dl 123 -t 1",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("课程ID错误")
		}
		aid := 0
		if len(args) > 1 {
			aid, err = strconv.Atoi(args[1])
			if err != nil {
				return errors.New("文章ID错误")
			}
		}
		err = download(app.CateCourse, id, aid)
		return err
	},
}

var dlOdobCmd = &cobra.Command{
	Use:   "dlo",
	Short: "下载每天听本书音频 & 文稿",
	Long: `使用 dedao-dl dlo 下载每天听本书音频, 并转换成 PDF(未实现) & 音频 & markdown
-t 指定下载格式, 1:mp3, 2:PDF文档, 3:markdown文档, 默认 mp3`,
	Example: "dedao-dl dlo 123 -t 1",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("听书ID错误")
		}
		aid := 0
		if len(args) > 1 {
			return errors.New("参数错误")
		}
		err = download(app.CateAudioBook, id, aid)
		return err
	},
}

var dlEbookCmd = &cobra.Command{
	Use:     "dle",
	Short:   "下载电子书",
	Long:    `使用 dedao-dl dle 下载电子书`,
	Example: "dedao-dl dle 123",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("电子书ID错误")
		}
		aid := 0
		if len(args) > 1 {
			return errors.New("参数错误")
		}
		err = download(app.CateEbook, id, aid)
		return err
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(dlOdobCmd)
	rootCmd.AddCommand(dlEbookCmd)
	downloadCmd.PersistentFlags().IntVarP(&downloadType, "downloadType", "t", 1, "下载格式, 1:mp3, 2:PDF文档, 3:markdown文档")
	dlOdobCmd.PersistentFlags().IntVarP(&downloadType, "downloadType", "t", 1, "下载格式, 1:mp3, 2:PDF文档, 3:markdown文档")
}

func download(cType string, id, aid int) error {
	switch cType {
	case app.CateCourse:
		course, err := app.CourseInfo(id)
		if err != nil {
			return err
		}
		articles, err := app.ArticleList(id, "")
		if err != nil {
			return err
		}

		switch downloadType {
		case 1: // mp3
			downloadData := extractDownloadData(course, articles, aid, 1)
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
			downloadData := extractDownloadData(course, articles, aid, 2)
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
			if err := DownloadMarkdown(app.CateCourse, id, path); err != nil {
				return err
			}
		}

	case app.CateAudioBook:
		fileName := "每天听本书"
		switch downloadType {
		case 1:
			downloadData := downloader.Data{
				Title: fileName,
			}
			downloadData.Type = "audio"
			downloadData.Data = extractOdobDownloadData(id)
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
			if err := DownloadMarkdown(app.CateAudioBook, id, path); err != nil {
				return err
			}
		}

	case app.CateEbook:
		detail, err := app.EbookDetail(id)
		if err != nil {
			return err
		}

		title := strconv.Itoa(id) + "_"
		if detail.Title != "" {
			title += detail.Title
		} else if detail.OperatingTitle != "" {
			title += detail.OperatingTitle
		}

		title += "_" + detail.BookAuthor
		if _, err := app.EbookPage(title, detail.Enid); err != nil {
			return err
		}

	}
	return nil
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
	article := config.Instance.GetIDMap(app.CateAudioBook, aid)
	aliasID := article["audio_alias_id"].(string)
	if aliasID == "" {
		list, err := app.CourseList(app.CateAudioBook)
		if err != nil {
			return nil
		}
		for _, course := range list.List {
			if aid > 0 && course.ID == aid {
				article = app.GetCourseIDMap(&course)
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
