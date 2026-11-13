package model

// 服务端保存和返回的用户数据
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u User) TableName() string {
	return "user"
}
