package downloader

import (
	"sort"
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
