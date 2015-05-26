package serial

import (
	"fmt"
	"testing"
)

func Test_getSysInfo(t *testing.T) {
	info, err := GetSysInfo()
	if err != nil {
		fmt.Println(err, info)
	}

}

func Test_ClientGetRegInfo(t *testing.T) {
	sn, err := ClientGetRegInfo()
	fmt.Println(err, sn)
}

func Test_ClientVerifySn(t *testing.T) {
	sn := "MDEwMDIwMTUwNTE5Pw+iC9ew/2Hy6H8jemybSS6yGABmV7r/vHusX7+2hwU="
	err := ClientVerifySn(sn)
	if err == nil {
		ClientSaveLicense(sn)
	}
}

func Test_ClientgetValidDate(t *testing.T) {
	date, err := ClientgetValidDate()
	fmt.Println(err, date)
}
