// 获取目录全路径接口

package RootDir

import (
	"errors"
	"os"
	"path/filepath"
)

func GetRootDir() (rootdir string, err error) {
	ProgramFiles := os.Getenv("ProgramFiles")
	naseClient := filepath.Join(ProgramFiles, "NaseForXP\\NaseClient")

	rootdir, err = filepath.Abs(naseClient)
	if err != nil {
		return "", errors.New("错误:获取根路径失败")
	}
	return rootdir, nil
}
