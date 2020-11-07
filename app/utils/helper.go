//小助手函数，辅助业务逻辑的函数
package utils

import "regexp"

var (
	regAdr = regexp.MustCompile("(a|A)ndroid|dr")
	regIos = regexp.MustCompile("i(p|P)(hone|ad|od)|(m|M)ac")
)

// GetDeviceByUa 根据ua获取设备名称
func GetDeviceByUa(ua string) string {
	plat := "web"
	if regAdr.MatchString(ua) {
		plat = "android"
	} else if regIos.MatchString(ua) {
		plat = "ios"
	}

	return plat
}
