# csb
阿里专有云csb sdk
> 仅支持http调用

根据阿里云官方sdk修改，支持`context`和`go module`,官方`sdk`不支持

使用方法
```go
package cron

import (
	"context"
	"fmt"

	"github.com/fnoopv/csb"
)

type Result struct {
	Data []struct {
		UserID          string `json:"USER_ID"`
		Mobile          string `json:"MOBILE"`
		Email           string `json:"EMAIL,omitempty"`
    }
	DataSize    int    `json:"dataSize"`
	Total       int    `json:"total"`
	ResultCode  int    `json:"resultCode"`
	ResultMsg   string `json:"resultMsg"`
	HasNextPage bool   `json:"hasNextPage"`
}

func SyncAvicUser() {
	c := csb.NewCSBClient("http://dadsa.com:8888/CSB", "dasdsa", "dadadadas=")
	c.SetApiName("users")
	c.SetApiMethod("get")
	c.SetApiVersion("1.0.0")
	c.SetContentType("application/json;charset=utf-8")

    // 添加query参数
	queryParam := make(map[string]string)
	queryParam["startTime"] = "0"
	c.SetQueryParam(queryParam)
    // 添加请求头
	headers := make(map[string]string)
	headers["appKey"] = "dasdadasd"
	c.SetHeaders(headers)

	result := Result{}
	err := c.Do(context.Background(), &result)
	if err != nil {
		fmt.Printf("request error: %s\n", err)
		return
	}
    fmt.Println(result)
}

```