package util

import "net/http"

// 拷贝 http 头信息

func CloneHeader(src http.Header, dst *http.Header) {
	for k, v := range src {
		dst.Set(k, v[0])
	}
}
