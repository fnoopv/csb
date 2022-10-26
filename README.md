# csb
阿里专有云csb sdk
> 仅支持http调用

根据阿里云官方sdk修改，支持`context`和`go module`,官方`sdk`不支持

使用方法
```go
package main

import (
    "context"
    "fmt"

    "github.com/fnoopv/csb"
)

type User struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Avatar string `json:"avatar"`
}

func main() {
    client := csb.NewCSBClient()
    client.SetUrl("http://xxx.xx.xx.xx:8888")
    client.SetAccessKey("asdasdasdsa")
    client.SetSecretKey("asfasfafasfas==")
    client.SetApiName("user")
    client.SetApiMethod("get")
    client.SetApiVersion("1.0")
    client.SetContentType("application/json; charset=utf8")
    // 添加其他请求头
    headers := map[string]string
    headers["abc"] = "abcdef"
    client.SetHeaders(headers)
    // 添加查询参数
    queryParams := map[string]string
    queryParams["id"] = "123"
    client.SetQueryParam(queryParams)

    //发起请求
    user := User{}
    if err := client.Do(context.Background(), &user); err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(user)
}
```