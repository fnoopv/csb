package csb

import (
	"net/url"
	"strings"
)

// appendParams 拼接参数
func appendParams(reqUrl string, params map[string]string) string {
	u := url.Values{}
	for k, v := range params {
		u.Set(k, v)
	}
	if strings.Contains(reqUrl, "?") {
		return reqUrl + "&" + u.Encode()
	} else {
		return reqUrl + "?" + u.Encode()
	}
}
