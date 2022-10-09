package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/yann0917/dedao-dl/request"
)

var OutputDir = "output"

// MAXLENGTH Maximum length of file name
const MAXLENGTH = 80

// TimeFormat format
const TimeFormat = "2006-01-02 15:04:05"

// FileName filter invalid string
func FileName(name, ext string) string {
	rep := strings.NewReplacer("\n", " ", "/", " ", "|", "-", ": ", "：", ":", "：", "'", "’", "\t", " ")
	name = rep.Replace(name)

	if runtime.GOOS == "windows" {
		rep = strings.NewReplacer("\"", " ", "?", " ", "*", " ", "\\", " ", "<", " ", ">", " ", ":", " ", "：", " ")
		name = rep.Replace(name)
	}

	name = strings.TrimSpace(name)

	limitedName := LimitLength(name, MAXLENGTH)
	if ext != "" {
		return fmt.Sprintf("%s.%s", limitedName, ext)
	}
	return limitedName
}

// LimitLength cut string
func LimitLength(s string, length int) string {
	ellipses := "..."

	str := []rune(s)
	if len(str) > length {
		s = string(str[:length-len(ellipses)]) + ellipses
	}

	return s
}

// FilePath gen valid file path
func FilePath(name, ext string, escape bool) (string, error) {
	var outputPath string

	fileName := name
	if escape {
		fileName = FileName(name, ext)
	} else {
		if ext != "" {
			fileName = fmt.Sprintf("%s.%s", name, ext)
		}
	}
	outputPath = filepath.Join(fileName)
	return outputPath, nil
}

// Mkdir mkdir path
func Mkdir(elem ...string) (string, error) {
	path := filepath.Join(elem...)

	err := os.MkdirAll(path, os.ModePerm)

	return path, err
}

// FileSize return the file size of the specified path file
func FileSize(filePath string) (int, bool, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return int(file.Size()), true, nil
}

// M3u8URLs get all ts urls from m3u8 url
func M3u8URLs(uri string) (urls []string, err error) {
	if len(uri) == 0 {
		return nil, errors.New("M3u8地址为空")
	}

	html, err := request.HTTPGet(uri)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(html), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "http") {
				urls = append(urls, line)
			} else {
				base, err := url.Parse(uri)
				if err != nil {
					continue
				}
				u, err := url.Parse(line)
				if err != nil {
					continue
				}
				urls = append(urls, base.ResolveReference(u).String())
			}
		}
	}
	return
}

// CurrentDir CurrentDir
func CurrentDir(joinPath ...string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	p := strings.Replace(dir, "\\", "/", -1)
	whole := filepath.Join(joinPath...)
	whole = filepath.Join(p, whole)
	return whole, nil
}

// ResolveURL parse url
func ResolveURL(u *url.URL, p string) string {
	if strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "http://") {
		return p
	}
	var baseURL string
	if strings.Index(p, "/") == 0 {
		baseURL = u.Scheme + "://" + u.Host
	} else {
		tU := u.String()
		baseURL = tU[0:strings.LastIndex(tU, "/")]
	}
	return baseURL + path.Join("/", p)
}

// DrawProgressBar Draw ProgressBar
func DrawProgressBar(prefix string, proportion float32, width int, suffix ...string) {
	pos := int(proportion * float32(width))
	s := fmt.Sprintf("[%s] %s%*s %6.2f%% %s",
		prefix, strings.Repeat("■", pos), width-pos, "", proportion*100, strings.Join(suffix, ""))
	fmt.Print("\r" + s)
}

// Unix2String 时间戳[转换为]字符串 eg:(2019-09-09 09:09:09)
func Unix2String(stamp int64) string {
	str := time.Unix(stamp, 0).Format(TimeFormat)
	return str
}

// Contains int in array
func Contains(s []int, n int) bool {
	for _, a := range s {
		if a == n {
			return true
		}
	}
	return false
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func CheckFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func WriteFileWithTrunc(filename, content string) (err error) {

	var f *os.File
	if CheckFileExist(filename) {
		f, err = os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)

		if err != nil {
			return
		}
	} else {
		f, err = os.Create(filename)
		if err != nil {
			return
		}
	}
	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		return
	}

	err = f.Sync()
	return

}

func MD5str(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
