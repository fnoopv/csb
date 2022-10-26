package csb

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// CSBClient CSBClient
type CSBClient struct {
	Url         string            // csb地址
	AccessKey   string            // ak
	SecretKey   string            // sk
	ApiName     string            // 接口名称
	ApiMethod   string            // 接口请求方法
	ApiVersion  string            // 接口版本
	ContentType string            // 请求的content-type
	Headers     map[string]string // 请求头
	QueryParam  map[string]string // query参数
	FormParam   map[string]string // 表单数据
	Body        []byte            // 请求体,文件、表单、JSON等
}

const (
	CSB_SDK_VERSION = "1.1.0"

	API_NAME_KEY               = "_api_name"
	VERSION_KEY                = "_api_version"
	ACCESS_KEY                 = "_api_access_key"
	SECRET_KEY                 = "_api_secret_key"
	SIGNATURE_KEY              = "_api_signature"
	TIMESTAMP_KEY              = "_api_timestamp"
	RESTFUL_PATH_SIGNATURE_KEY = "csb_restful_path_signature_key" //TODO: fix the terrible key name!
)

// NewCSBClient 返回新的CSB客户端
func NewCSBClient() *CSBClient {
	return &CSBClient{}
}

// SetUrl 设置CSB地址
func (c *CSBClient) SetUrl(url string) *CSBClient {
	c.Url = url
	return c
}

// SetAccessKey 设置ak
func (c *CSBClient) SetAccessKey(accessKey string) *CSBClient {
	c.AccessKey = accessKey
	return c
}

// SetSecretKey 设置sk
func (c *CSBClient) SetSecretKey(secretKey string) *CSBClient {
	c.SecretKey = secretKey
	return c
}

// SetApiName 设置请求接口的名称
func (c *CSBClient) SetApiName(apiName string) *CSBClient {
	c.ApiName = apiName
	return c
}

// SetApiMethod 设置请求接口的方法,只支持get或post
func (c *CSBClient) SetApiMethod(apiMethod string) *CSBClient {
	c.ApiMethod = apiMethod
	return c
}

// SetApiVersion 设置请求接口的版本
func (c *CSBClient) SetApiVersion(apiVersion string) *CSBClient {
	c.ApiVersion = apiVersion
	return c
}

// SetContentType 设置请求content-type
func (c *CSBClient) SetContentType(contentType string) *CSBClient {
	c.ContentType = contentType
	return c
}

// SetHeaders 设置请求头
func (c *CSBClient) SetHeaders(headers map[string]string) *CSBClient {
	c.Headers = headers
	return c
}

// SetQueryParam 设置query参数对
func (c *CSBClient) SetQueryParam(queryParam map[string]string) *CSBClient {
	c.QueryParam = queryParam
	return c
}

// SetFormParam 设置表单数据
func (c *CSBClient) SetFormParam(formParam map[string]string) *CSBClient {
	c.FormParam = formParam
	return c
}

// SetBody 设置请求体
func (c *CSBClient) SetBody(body []byte) *CSBClient {
	c.Body = body
	return c
}

// Do 执行请求
func (c *CSBClient) Do(ctx context.Context, result interface{}) *CSBError {
	// 参数验证
	if err := c.validate(); err != nil {
		return err
	}
	client := &http.Client{}

	// 表单数据设置
	formData := url.Values{}
	if c.FormParam != nil {
		for k, v := range c.FormParam {
			formData.Set(k, v)
		}
	}

	// merge query param to url
	link := c.Url
	if c.QueryParam != nil {
		_, err := url.Parse(link)
		if err != nil {
			return &CSBError{Message: "bad request url"}
		}
		link = appendParams(link, c.QueryParam)
	}

	// merge body
	requestBody := c.Body
	if c.Body == nil {
		requestBody = []byte(formData.Encode())
	}

	req, err := http.NewRequest(strings.ToLower(c.ApiMethod), link, bytes.NewReader(requestBody))
	if err != nil {
		return &CSBError{Message: "failed to construct http post request", CauseErr: err}
	}
	req = req.WithContext(ctx)

	// merge params
	params := c.QueryParam
	if c.FormParam != nil {
		for k, v := range c.FormParam {
			params[k] = v
		}
	}

	// add request header
	signHeaders := signParams(params, c.ApiName, c.ApiVersion, c.AccessKey, c.SecretKey)
	if c.Headers != nil {
		for k, v := range c.Headers {
			req.Header.Add(k, v)
		}
	}
	for k, v := range signHeaders {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &CSBError{Message: "failed to request http post", CauseErr: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &CSBError{Message: "read response body failed", CauseErr: err}
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.Unmarshal(body, result); err != nil {
			return &CSBError{Message: "json unmarshal failed", CauseErr: err}
		}
	} else if strings.Contains(contentType, "text/xml") {
		if err := xml.Unmarshal(body, result); err != nil {
			return &CSBError{Message: "xml unmarshal failed", CauseErr: err}
		}
	} else {
		result = string(body)
	}

	return nil
}

// signParams 对参数进行签名
func signParams(params map[string]string, api string, version string, ak string, sk string) (headMaps map[string]string) {
	headMaps = make(map[string]string)

	params[API_NAME_KEY] = api
	headMaps[API_NAME_KEY] = api

	params[VERSION_KEY] = version
	headMaps[VERSION_KEY] = version

	v := time.Now().UnixNano() / 1000000
	params[TIMESTAMP_KEY] = strconv.FormatInt(v, 10)
	headMaps[TIMESTAMP_KEY] = strconv.FormatInt(v, 10)

	if ak != "" {
		params[ACCESS_KEY] = ak
		headMaps[ACCESS_KEY] = ak

		delete(params, SECRET_KEY)
		delete(params, SIGNATURE_KEY)

		signValue := doSign(params, sk)

		headMaps[SIGNATURE_KEY] = signValue
	}

	return headMaps
}

// validate 验证参数
func (c *CSBClient) validate() *CSBError {
	method := strings.ToLower(c.ApiMethod)
	if method != "get" && method != "post" {
		return &CSBError{Message: "bad method, only support 'get' or 'post'"}
	}
	if c.AccessKey == "" || c.SecretKey == "" {
		return &CSBError{Message: "bad request params, accessKey and secretKey must defined together"}
	}
	if c.ApiName == "" || c.ApiVersion == "" {
		return &CSBError{Message: "bad request params, api or version not defined"}
	}
	if c.ContentType == "" {
		return &CSBError{Message: "content-type must defined"}
	}

	return nil
}
