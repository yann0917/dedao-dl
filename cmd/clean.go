package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/utils"
)

const (
	cleanTargetOutput = "output"
	cleanTargetCache  = "cache"
)

type cleanResult struct {
	Target string `json:"target"`
	Path   string `json:"path"`
	Status string `json:"status"`
}

var cleanCmd = &cobra.Command{
	Use:       "clean [output|cache]",
	Short:     "清理 output 或 .cache 目录",
	Long:      "使用 dedao-dl clean output 或 dedao-dl clean cache 清理工作目录下的输出或缓存文件夹。",
	Example:   "dedao-dl clean output\ndedao-dl clean cache",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{cleanTargetOutput, cleanTargetCache},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runClean(args[0])
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func runClean(target string) error {
	path, err := cleanTargetPath(target)
	if err != nil {
		return err
	}

	if target == cleanTargetCache {
		if err := utils.CloseBadgerDB(); err != nil {
			return fmt.Errorf("关闭缓存数据库失败: %w", err)
		}
	}

	if err := resetDir(path); err != nil {
		return err
	}

	if outputJSON {
		return printJSON(cleanResult{
			Target: target,
			Path:   path,
			Status: "cleaned",
		})
	}

	fmt.Printf("已清理 %s 目录: %s\n", cleanTargetLabel(target), path)
	return nil
}

func cleanTargetPath(target string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前目录失败: %w", err)
	}

	switch target {
	case cleanTargetOutput:
		return filepath.Join(cwd, utils.OutputDir), nil
	case cleanTargetCache:
		return filepath.Join(cwd, ".cache"), nil
	default:
		return "", fmt.Errorf("不支持的清理目标: %s", target)
	}
}

func cleanTargetLabel(target string) string {
	switch target {
	case cleanTargetOutput:
		return "output"
	case cleanTargetCache:
		return ".cache"
	default:
		return target
	}
}

func resetDir(path string) error {
	if err := ensureCleanPath(path); err != nil {
		return err
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("删除目录失败: %w", err)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("重建目录失败: %w", err)
	}

	return nil
}

func ensureCleanPath(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("解析目录路径失败: %w", err)
	}

	base := filepath.Base(absPath)
	if base != utils.OutputDir && base != ".cache" {
		return fmt.Errorf("拒绝清理非预期目录: %s", absPath)
	}

	root := filepath.Dir(absPath)
	if root == absPath {
		return fmt.Errorf("拒绝清理根目录: %s", absPath)
	}

	return nil
}
