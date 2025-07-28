package vo

type TodoUrgeReq struct {
	TodoId int64  `json:"todoId,string"`
	Msg    string `json:"msg"`
}
