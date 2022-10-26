package csb

import "fmt"

type CSBError struct {
	Message   string `xml:"Message"`   // csb返回的错误信息
	RequestID string `xml:"RequestId"` // csb返回的请求uuid
	CauseErr  error  //具体的错误信息
}

// Error 返回错误信息，实现Error接口
func (e CSBError) Error() string {
	return fmt.Sprintf("csb: service returned error: ErrorMessage=%s, RequestId=%s, CauseErr=%v", e.Message, e.RequestID, e.CauseErr)
}
