package csb

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"hash"
	"io"
	"sort"
)

// paramSorter 排序容器
type paramSorter struct {
	Keys   []string
	Values []string
}

// params map 转换为paramsSorter格式
func newParamsSorter(m map[string]string) *paramSorter {
	hs := &paramSorter{
		Keys:   make([]string, 0, len(m)),
		Values: make([]string, 0, len(m)),
	}

	for k, v := range m {
		hs.Keys = append(hs.Keys, k)
		hs.Values = append(hs.Values, v)
	}
	return hs
}

// 进行字典顺序排序 sort required method
func (hs *paramSorter) Sort() {
	sort.Sort(hs)
}

// Additional function for function  sort required method
func (hs *paramSorter) Len() int {
	return len(hs.Values)
}

// Additional function for function  sort required method
func (hs *paramSorter) Less(i, j int) bool {
	return bytes.Compare([]byte(hs.Keys[i]), []byte(hs.Keys[j])) < 0
}

// Additional function for function paramsSorter.
func (hs *paramSorter) Swap(i, j int) {
	hs.Values[i], hs.Values[j] = hs.Values[j], hs.Values[i]
	hs.Keys[i], hs.Keys[j] = hs.Keys[j], hs.Keys[i]
}

// 做签名处理
func doSign(params map[string]string, secretKey string) string {
	hs := newParamsSorter(params)

	// Sort the temp by the Ascending Order
	hs.Sort()

	// Get the CanonicalizedOSSHeaders
	canonicalizedParams := ""
	for i := range hs.Keys {
		if i > 0 {
			canonicalizedParams += "&"
		}
		canonicalizedParams += hs.Keys[i] + "=" + hs.Values[i]
	}

	signStr := canonicalizedParams

	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(secretKey))
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedStr
}
