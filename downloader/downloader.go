package downloader

import (
	"fmt"
	"path/filepath"

	"github.com/yann0917/dedao-dl/utils"
)

//PrintToPDF print to pdf
func PrintToPDF(v Datum, cookies map[string]string, path string) error {

	name := utils.FileName(v.Title, "pdf")

	filename := filepath.Join(path, name)
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", name)

	_, exist, err := utils.FileSize(filename)

	if err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败")
		return err
	}

	if exist {
		fmt.Printf("\033[33;1m%s\033[0m\n", "已存在")
		return nil
	}

	err = utils.ColumnPrintToPDF(v.Enid, filename, cookies)

	if err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败")
		return err
	}

	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")

	return nil
}
