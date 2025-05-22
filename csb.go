package csb

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// CSBClient CSBClient
type CSBClient struct {
	url         string              // csb地址
	accessKey   string              // ak
	secretKey   string              // sk
	ApiName     string              // 接口名称
	ApiMethod   string              // 接口请求方法
	ApiVersion  string              // 接口版本
	ContentType string              // 请求的content-type
	Headers     map[string][]string // 请求头
	QueryParam  map[string]string   // query参数
	FormParam   map[string]string   // 表单数据
	Body        []byte              // 请求体,文件、表单、JSON等
	client      *resty.Client
}

const (
	apiNameKey       = "_api_name"
	apiVersionKey    = "_api_version"
	accessKey        = "_api_access_key"
	secretKey        = "_api_secret_key"
	signatureKey     = "_api_signature"
	timestampKey     = "_api_timestamp"
	defaultUserAgent = "csbBroker"
)

// NewCSBClient 返回新的CSB客户端
func NewCSBClient(url, accessKey, secretKey string) *CSBClient {

	return &CSBClient{
		client:    resty.New(),
		url:       url,
		accessKey: accessKey,
		secretKey: secretKey,
	}
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
func (c *CSBClient) SetHeaders(headers map[string][]string) *CSBClient {
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
func (c *CSBClient) Do(ctx context.Context) (res *resty.Response, body []byte, err error) {

	// 参数验证
	if err := c.validate(); err != nil {
		return nil, nil, err
	}

	// 表单数据设置
	formData := url.Values{}
	if c.FormParam != nil {
		for k, v := range c.FormParam {
			formData.Set(k, v)
		}
	}

	requestBody := c.Body
	if c.Body == nil {
		requestBody = []byte(formData.Encode())
	}

	// merge params
	params := make(map[string]string)
	if c.QueryParam != nil {
		params = c.QueryParam
	}
	if c.FormParam != nil {
		for k, v := range c.FormParam {
			params[k] = v
		}
	}

	// add request header
	signHeaders := signParams(params, c.ApiName, c.ApiVersion, c.accessKey, c.secretKey)

	r := c.client.R().
		SetQueryParams(c.QueryParam).
		SetBody(requestBody).
		SetHeader("Content-Type", c.ContentType).
		SetHeaderMultiValues(signHeaders).
		SetHeaderMultiValues(c.Headers)

	method := strings.ToLower(c.ApiMethod)

	if method == "get" {
		res, err = r.Get(c.url)
	} else if method == "post" {
		res, err = r.Post(c.url)
	}
	if err != nil {
		return res, nil, err
	}

	return res, res.Body(), nil
}

// signParams 对参数进行签名
func signParams(
	params map[string]string,
	api string,
	version string,
	ak string,
	sk string,
) (headMaps map[string][]string) {
	headMaps = make(map[string][]string)

	params[apiNameKey] = api
	headMaps[apiNameKey] = []string{api}

	params[apiVersionKey] = version
	headMaps[apiVersionKey] = []string{version}

	v := time.Now().UnixNano() / 1000000
	params[timestampKey] = strconv.FormatInt(v, 10)
	headMaps[timestampKey] = []string{strconv.FormatInt(v, 10)}

	params[accessKey] = ak
	headMaps[accessKey] = []string{ak}

	delete(params, secretKey)
	delete(params, signatureKey)

	signValue := doSign(params, sk)

	headMaps[signatureKey] = []string{signValue}

	return headMaps
}

// validate 验证参数
func (c *CSBClient) validate() error {
	method := strings.ToLower(c.ApiMethod)
	if method != "get" && method != "post" {
		return errors.New("bad method, only support 'get' or 'post'")
	}
	if c.accessKey == "" || c.secretKey == "" {
		return errors.New("bad request params, accessKey and secretKey must defined together")
	}
	if c.ApiName == "" || c.ApiVersion == "" {
		return errors.New("bad request params, api or version not defined")
	}
	if c.ContentType == "" {
		return errors.New("content-type must defined")
	}

	return nil
}
