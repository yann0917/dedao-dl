package parse

// Partial reference https://github.com/grafov/m3u8/blob/master/reader.go
import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type (
	// PlaylistType Playlist Type
	PlaylistType string
	// CryptMethod Crypt Method
	CryptMethod string
)

const (
	// PlaylistTypeVOD VOD
	PlaylistTypeVOD PlaylistType = "VOD"
	// PlaylistTypeEvent EVENT
	PlaylistTypeEvent PlaylistType = "EVENT"
	// CryptMethodAES AES-128
	CryptMethodAES CryptMethod = "AES-128"
	// CryptMethodNONE NONE
	CryptMethodNONE CryptMethod = "NONE"
)

// regex pattern for extracting `key=value` parameters from a line
var linePattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)

// M3u8 M3U8
type M3u8 struct {
	MediaSequence  uint64 // Default 0, #EXT-X-MEDIA-SEQUENCE:sequence
	Segments       []*Segment
	MasterPlaylist []*MasterPlaylist
	Keys           map[int]*Key
	PlaylistType   PlaylistType // VOD or EVENT
	TargetDuration float64      // #EXT-X-TARGETDURATION:duration
	Version        int8         // EXT-X-VERSION:version
	EndList        bool         // #EXT-X-ENDLIST
}

// Segment EXTINF:
type Segment struct {
	URI      string
	KeyIndex int
	Title    string  // #EXTINF: duration,<title>
	Duration float32 // #EXTINF: duration,<title>
	Length   uint64  // #EXT-X-BYTERANGE: length[@offset]
	Offset   uint64  // #EXT-X-BYTERANGE: length[@offset]
}

// MasterPlaylist master playlist
//
// #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=240000,RESOLUTION=416x234,CODECS="avc1.42e00a,mp4a.40.2"
type MasterPlaylist struct {
	URI        string
	BandWidth  uint32
	Resolution string
	Codecs     string
	ProgramID  uint32
}

// Key EXT-X-KEY:
// #EXT-X-KEY:METHOD=AES-128,URI="key.key",KEYFORMAT="identity",IV=0xc2b6c3bc4cc3afc2962d737f2142c29d
type Key struct {
	// 'AES-128' or 'NONE'
	// If the encryption method is NONE, the URI and the IV attributes MUST NOT be present
	Method CryptMethod
	URI    string
	IV     string
}

func parse(reader io.Reader) (*M3u8, error) {
	s := bufio.NewScanner(reader)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	var (
		i     = 0
		count = len(lines)
		m3u8  = &M3u8{
			Keys: make(map[int]*Key),
		}
		keyIndex = 0

		key     *Key
		seg     *Segment
		extInf  bool
		extByte bool
	)

	for ; i < count; i++ {
		line := strings.TrimSpace(lines[i])
		if i == 0 {
			if "#EXTM3U" != line {
				return nil, fmt.Errorf("invalid m3u8, missing #EXTM3U in line 1")
			}
			continue
		}
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-PLAYLIST-TYPE:%s", &m3u8.PlaylistType); err != nil {
				return nil, err
			}
			isValid := m3u8.PlaylistType == "" || m3u8.PlaylistType == PlaylistTypeVOD || m3u8.PlaylistType == PlaylistTypeEvent
			if !isValid {
				return nil, fmt.Errorf("invalid playlist type: %s, line: %d", m3u8.PlaylistType, i+1)
			}
		case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &m3u8.TargetDuration); err != nil {
				return nil, err
			}
		case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &m3u8.MediaSequence); err != nil {
				return nil, err
			}
		case strings.HasPrefix(line, "#EXT-X-VERSION:"):
			if _, err := fmt.Sscanf(line, "#EXT-X-VERSION:%d", &m3u8.Version); err != nil {
				return nil, err
			}
		// Parse master playlist
		case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
			mp, err := parseMasterPlaylist(line)
			if err != nil {
				return nil, err
			}
			i++
			mp.URI = lines[i]
			if mp.URI == "" || strings.HasPrefix(mp.URI, "#") {
				return nil, fmt.Errorf("invalid EXT-X-STREAM-INF URI, line: %d", i+1)
			}
			m3u8.MasterPlaylist = append(m3u8.MasterPlaylist, mp)
			continue
		case strings.HasPrefix(line, "#EXTINF:"):
			if extInf {
				return nil, fmt.Errorf("duplicate EXTINF: %s, line: %d", line, i+1)
			}
			if seg == nil {
				seg = new(Segment)
			}
			var s string
			if _, err := fmt.Sscanf(line, "#EXTINF:%s", &s); err != nil {
				return nil, err
			}
			if strings.Contains(s, ",") {
				split := strings.Split(s, ",")
				seg.Title = split[1]
				s = split[0]
			}
			df, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return nil, err
			}
			seg.Duration = float32(df)
			seg.KeyIndex = keyIndex
			extInf = true
		case strings.HasPrefix(line, "#EXT-X-BYTERANGE:"):
			if extByte {
				return nil, fmt.Errorf("duplicate EXT-X-BYTERANGE: %s, line: %d", line, i+1)
			}
			if seg == nil {
				seg = new(Segment)
			}
			var b string
			if _, err := fmt.Sscanf(line, "#EXT-X-BYTERANGE:%s", &b); err != nil {
				return nil, err
			}
			if b == "" {
				return nil, fmt.Errorf("invalid EXT-X-BYTERANGE, line: %d", i+1)
			}
			if strings.Contains(b, "@") {
				split := strings.Split(b, "@")
				offset, err := strconv.ParseUint(split[1], 10, 64)
				if err != nil {
					return nil, err
				}
				seg.Offset = uint64(offset)
				b = split[0]
			}
			length, err := strconv.ParseUint(b, 10, 64)
			if err != nil {
				return nil, err
			}
			seg.Length = uint64(length)
			extByte = true
		// Parse segments URI
		case !strings.HasPrefix(line, "#"):
			if extInf {
				if seg == nil {
					return nil, fmt.Errorf("invalid line: %s", line)
				}
				seg.URI = line
				extByte = false
				extInf = false
				m3u8.Segments = append(m3u8.Segments, seg)
				seg = nil
				continue
			}
		// Parse key
		case strings.HasPrefix(line, "#EXT-X-KEY"):
			params := parseLineParameters(line)
			if len(params) == 0 {
				return nil, fmt.Errorf("invalid EXT-X-KEY: %s, line: %d", line, i+1)
			}
			method := CryptMethod(params["METHOD"])
			if method != "" && method != CryptMethodAES && method != CryptMethodNONE {
				return nil, fmt.Errorf("invalid EXT-X-KEY method: %s, line: %d", method, i+1)
			}
			keyIndex++
			key = new(Key)
			key.Method = method
			key.URI = params["URI"]
			key.IV = params["IV"]
			m3u8.Keys[keyIndex] = key
		case line == "#EndList":
			m3u8.EndList = true
		default:
			continue
		}
	}

	return m3u8, nil
}

func parseMasterPlaylist(line string) (*MasterPlaylist, error) {
	params := parseLineParameters(line)
	if len(params) == 0 {
		return nil, errors.New("empty parameter")
	}
	mp := new(MasterPlaylist)
	for k, v := range params {
		switch {
		case k == "BANDWIDTH":
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, err
			}
			mp.BandWidth = uint32(v)
		case k == "RESOLUTION":
			mp.Resolution = v
		case k == "PROGRAM-ID":
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, err
			}
			mp.ProgramID = uint32(v)
		case k == "CODECS":
			mp.Codecs = v
		}
	}
	return mp, nil
}

// parseLineParameters extra parameters in string `line`
func parseLineParameters(line string) map[string]string {
	r := linePattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
}
