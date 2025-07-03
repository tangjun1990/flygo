package xsign

import (
	"crypto/md5"
	"fmt"
	"sort"
)

func GenSign(secret string, params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//拼接
	var dataParams string
	for _, k := range keys {
		// 排除空值
		if "" != params[k] {
			dataParams = dataParams + k + "=" + params[k] + "&"
		}
	}

	ff := dataParams[0 : len(dataParams)-1]

	return MD5(ff + secret)
}

func MD5(s string) string {
	has := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", has)
}
