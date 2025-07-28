package callsvc

type CustomUserEvent struct {
	CustomEvent
	ToUserName string `xml:"ToUserName"`
	UserID     string `xml:"UserID"`
	NewUserID  string `xml:"NewUserID"`
	Department []int  `xml:"Department"`
}

type CustomDeptEvent struct {
	CustomEvent
	ToUserName string `xml:"ToUserName"`
	Id         int    `xml:"Id"`
	ParentId   int    `xml:"ParentId"`
	Name       string `xml:"Name"`
}

// DataUserEvent 上架应用
type DataUserEvent struct {
	DataEvent
	UserID     string `xml:"UserID"`
	NewUserID  string `xml:"NewUserID"`
	Department []int  `xml:"Department"`
	Name       string `xml:"Name"`
	OpenUserID string `xml:"OpenUserID"`
}

type DataDeptEvent struct {
	DataEvent
	Id       int    `xml:"Id"`
	ParentId int    `xml:"ParentId"`
	Name     string `xml:"Name"`
}
