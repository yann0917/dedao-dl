package downloader

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// URL for url information
type URL struct {
	URL  string `json:"url"`
	Size int    `json:"size"`
	Ext  string `json:"ext"`
}

// Stream data
type Stream struct {
	URLs    []URL  `json:"urls"`
	Size    int    `json:"size"`
	Quality string `json:"quality"`
	name    string
}

// Datum download information
type Datum struct {
	ID        int    `json:"id"`
	ClassID   int    `json:"class_id"`
	Enid      string `json:"enid"`
	ClassEnid string `json:"class_enid"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	IsCanDL   bool   `json:"is_can_dl"`
	M3U8URL   string `json:"m3u8_url"`

	Streams       map[string]Stream `json:"streams"`
	sortedStreams []Stream
}

// Data 课程信息
type Data struct {
	Title string  `json:"title"`
	Type  string  `json:"type"`
	Data  []Datum `json:"articles"`
}

// VideoMediaMap 视频大小信息
type VideoMediaMap struct {
	Size int `json:"size"`
}

// AudioMediaMap 音频频大小信息
type AudioMediaMap struct {
	Size int `json:"size"`
}

// EmptyData empty data list
var EmptyData = make([]Datum, 0)

// PrintInfo print info
func (data *Data) PrintInfo() {
	if len(data.Data) == 0 {
		fmt.Println(data.Type + "目录为空")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)

	header := []string{"#", "ID", "类型", "名称"}
	for key := range data.Data[0].Streams {
		header = append(header, key)
	}
	header = append(header, "下载")

	table.Header(header)
	i := 0
	for _, p := range data.Data {
		reg, _ := regexp.Compile(" \\| ")
		title := reg.ReplaceAllString(p.Title, " ")

		isCanDL := ""
		if p.IsCanDL {
			isCanDL = " ✔"
		}

		value := []string{strconv.Itoa(i), strconv.Itoa(p.ID), p.Type, title}

		if len(p.Streams) > 0 {
			for _, stream := range p.Streams {
				value = append(value, fmt.Sprintf("%.2fMB (%d Bytes)\n", float64(stream.Size)/(1024*1024), stream.Size))
			}
		} else {
			for range data.Data[0].Streams {
				value = append(value, " -")
			}
		}

		value = append(value, isCanDL)

		table.Append(value)
		i++
	}
	table.Render()
}

func (v *Datum) genSortedStreams() {
	for k, data := range v.Streams {
		if data.Size == 0 {
			data.calculateTotalSize()
		}
		data.name = k
		v.Streams[k] = data
		v.sortedStreams = append(v.sortedStreams, data)
	}
	if len(v.Streams) > 1 {
		sort.Slice(
			v.sortedStreams, func(i, j int) bool { return v.sortedStreams[i].Size > v.sortedStreams[j].Size },
		)
	}
}

func (stream *Stream) calculateTotalSize() {

	if stream.Size > 0 {
		return
	}

	size := 0
	for _, url := range stream.URLs {
		size += url.Size
	}
	stream.Size = size
}
