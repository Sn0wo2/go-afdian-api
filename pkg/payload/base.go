package payload

type Base struct {
	EC int    `json:"ec,omitempty"`
	EM string `json:"em,omitempty"`
}

// Checker 的创建是因为 afdian 官方API违背了RESTful API设计原则:
// 在Response Body中返回业务状态码EC, 却让HTTP Status Code始终返回200 OK。
// 因此需要通过Checker接口在doRequest泛型函数中读取EC和EM，避免进行两次Unmarshal。
//
// Thanks afdian open api for this "wonderful" design.
type Checker interface {
	GetEC() int
	GetEM() string
}

func (b *Base) GetEC() int {
	return b.EC
}

func (b *Base) GetEM() string {
	return b.EM
}
