// 获取目录全路径接口

package RootDir

import (
	"errors"
	"path/filepath"
)

func GetRootDir() (rootdir string, err error) {
	rootdir, err = filepath.Abs("../")
	if err != nil {
		return "", errors.New("错误:获取根路径失败")
	}
	return rootdir, nil
}
