package projectvo

type TodoUrgeInput struct {
	TodoId int64  `json:"todoId"`
	Msg    string `json:"msg"`
}

// TodoUrgeReq .
type TodoUrgeReq struct {
	OrgId  int64          `json:"orgId"`
	UserId int64          `json:"userId"`
	Input  *TodoUrgeInput `json:"input"`
}
