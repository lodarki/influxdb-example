package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"sort"
)

// @Description generate param sorted string
func GenerateSource(mapParams map[string]string, timestamp uint32) string {
	var params []string
	for k, _ := range mapParams {
		params = append(params, k)
	}

	sort.Strings(params)
	var s string
	for _, p := range params {
		s += p + mapParams[p]
	}
	return fmt.Sprintf("timestamp%d%stimestamp%d", timestamp, s, timestamp)
}

// @Description generate double md5 by params
func GenerateDoubleMD5ByParams(mapParams map[string]string, salt string, timestamp uint32) (md5str string, err error) {
	return GenerateDoubleMD5(GenerateSource(mapParams, timestamp), salt)
}

// @Description double md5 sum to source with salt.
func GenerateDoubleMD5(source, salt string) (md5str string, err error) {
	m := md5.New()
	_, err = io.WriteString(m, source)
	if err != nil {
		logs.Error("m.Write err", err)
		return
	}

	tmp := fmt.Sprintf("%x", m.Sum(nil))
	m.Reset()
	_, _ = io.WriteString(m, tmp)
	_, _ = io.WriteString(m, salt)
	md5str = fmt.Sprintf("%x", m.Sum(nil))
	return
}

// @Description compare src's md5 to dst.
func CompareMd5(dst, src, salt string) bool {
	tmp, err := GenerateDoubleMD5(src, salt)
	if err != nil {
		return false
	}

	if tmp != dst {
		return false
	}

	return true
}
