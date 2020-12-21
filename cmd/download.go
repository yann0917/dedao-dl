package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/downloader"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

var downloadCmd = &cobra.Command{
	Use:     "dl",
	Short:   "`dedao-dl dl` 下载已购买课程，并转换成 PDF 或者音频",
	Long:    `使用 dedao-dl dl 下载已购买课程，并转换成 PDF 或者音频`,
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("download cmd", args)
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
		err = download(id, aid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

func download(id, aid int) error {

	course, err := app.CourseInfo(id)
	if err != nil {
		return err
	}
	articles, err := app.ArticleList(id)
	downloadData := extractDownloadData(course, articles, aid)
	errors := make([]error, 0)
	path, err := utils.Mkdir(utils.FileName(course.ClassInfo.Name, ""), "MP3")

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
	// 下载 PDF
	// path, err := utils.Mkdir(utils.FileName(course.ClassInfo.Name, ""), "PDF")
	// if err != nil {
	// 	return err
	// }

	// cookies := LoginedCookies()
	// for _, datum := range downloadData.Data {
	// 	if !datum.IsCanDL {
	// 		continue
	// 	}
	// 	if err := downloader.PrintToPDF(datum, cookies, path); err != nil {
	// 		errors = append(errors, err)
	// 	}
	// }
	// if len(errors) > 0 {
	// 	return errors[0]
	// }
	return nil
}

//生成下载数据
func extractDownloadData(course *services.CourseInfo, articles *services.ArticleList, aid int) downloader.Data {

	// TODO: odob ,course, ebook compass are diff
	downloadData := downloader.Data{
		Title: course.ClassInfo.Name,
	}

	if course.ClassInfo.LogType == "class" {
		downloadData.Type = "audio"
		downloadData.Data = extractCourseDownloadData(articles, aid)
	}

	return downloadData
}

//生成下载数据
func extractCourseDownloadData(articles *services.ArticleList, aid int) []downloader.Datum {
	data := downloader.EmptyData
	audioIds := map[int]string{}

	audioData := make([]*downloader.Datum, 0)
	fmt.Println(aid)
	for _, article := range articles.List {
		if aid > 0 && article.ID != aid {
			continue
		}
		audioIds[article.ID] = article.Aduio.AliasID

		urls := []downloader.URL{}
		key := article.Enid
		streams := map[string]downloader.Stream{
			key: {
				URLs:    urls,
				Size:    article.Aduio.Size,
				Quality: key,
			},
		}

		datum := &downloader.Datum{
			ID:        article.ID,
			Enid:      article.Enid,
			ClassEnid: article.ClassEnid,
			ClassID:   article.ClassID,
			Title:     article.Title,
			IsCanDL:   true,
			M3U8URL:   article.Aduio.Mp3PlayURL,
			Streams:   streams,
			Type:      "audio",
		}

		audioData = append(audioData, datum)
	}

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

	for _, d := range audioData {
		data = append(data, *d)
	}
	return data
}

func printExtractDownloadData(v interface{}) {
	jsonData, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", jsonData)
	}
}
