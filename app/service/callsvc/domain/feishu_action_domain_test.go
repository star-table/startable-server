package callsvc

import (
	"html/template"
	"testing"
)

//func TestCreateIssueFeishu(t *testing.T) {
//	reqObj := "{\"token\":\"UDbcxblq3u84EQ0xYmkoPgq0ykwiCqw7\",\"type\":\"application_action_execute\",\"data\":{\"header\":{\"user_id\":\"ou_22f5529e4a1bc5c3a1271ca609b24ed5\",\"tenant_id\":\"2ed263bf32cf1651\",\"happen_at\":1603881044000},\"namespace\":\"polaris\",\"domain_id\":\"polarisTask\",\"action_id\":\"createIssue\",\"serial_id\":\"abc\",\"params\":{\"projectId\":\"1602\",\"projectObjectTypeId\":\"2005\",\"title\":\"捷径3\",\"remark\":\"ssssss\",\"priorityId\":\"1027\",\"planStartTime\":\"1603941439\",\"planEndTime\":\"1603941459\",\"owner\":\"1023\",\"followerIds\":[\"1024\"],\"tagIds\":[\"1001\"]}}}"
//	req := &bo.ActionReq{}
//	_ = json.FromJson(reqObj, req)
//	data := req.Data.Params
//	if followerIdInterface, ok := data["followerIds"]; ok {
//		followIds := &[]string{}
//		followIdStr, ok1 := followerIdInterface.(string)
//		fmt.Println(11111, followIdStr, followerIdInterface)
//		fmt.Println(reflect.TypeOf(followerIdInterface))
//		err := copyer.Copy(followerIdInterface, followIds)
//		fmt.Println(222, err)
//		fmt.Println(json.ToJsonIgnoreError(followIds))
//		if ok1 {
//			//jsonErr := json.FromJson(followIdStr, followIds)
//			//if jsonErr != nil {
//			//	log.Error(jsonErr)
//			//} else {
//			//	fmt.Println("res", json.ToJsonIgnoreError(followIds))
//			//}
//		}
//	}
//}

func TestMatchIsUserName(t *testing.T) {
	text := `<at open_id="yyy">@张三</at>`
	isUserName := CheckIsUserName(text)
	t.Log(isUserName)
	text = `<at open_id="12323123at>`
	isUserName = CheckIsUserName(text)
	t.Log(isUserName)
}

func TestCheckIsUserNameWithIssueTitle(t *testing.T) {
	text := `<at open_id="ou_83776f77689ab1a5f399d5aafa27a0a4">@苏汉13</at> <br>`
	hasTitle := CheckIsUserNameWithIssueTitle(text)
	if hasTitle {
		t.Error("“<br>” 应该被过滤掉，不能作为标题的一部分。")
	}
	{
		text := `你是<at open_id="ou_83776f77689ab1a5f399d5aafa27a0a4">@苏汉13</at> <br>`
		hasTitle := CheckIsUserNameWithIssueTitle(text)
		if !hasTitle {
			t.Error("此处应该匹配出标题。[1005]")
		}
	}
	text = `<at open_id="ou_83776f77689ab1a5f399d5aafa27a0a4">@苏汉13</at> 指令创建任务002`
	isOk := CheckIsUserNameWithIssueTitle(text)
	if !isOk {
		t.Error(isOk)
	}
	isOk = CheckIsUserName(text)
	if isOk {
		t.Error("这里不应该匹配成功001。")
	}
	atUserName, atUserOpenId, issueTitle := MatchUserInsUserNameWithIssueTitle(text)
	t.Log(atUserOpenId, atUserName, issueTitle)
	text = `<at open_id="ou_83776f77689ab1a5f399d5aafa27a0a4">@苏汉13</at> <br/><script> &nbsp;`
	atUserName, atUserOpenId, issueTitle = MatchUserInsUserNameWithIssueTitle(text)
	// 对标题中的标签进行转义
	issueTitle = template.HTMLEscapeString(issueTitle)
	t.Log(atUserOpenId, atUserName, issueTitle)
	text = `<at open_id="ou_83776f77689ab1a5f399d5aafa27a0a4">@苏汉13</at> <br/><script> test1-&nbsp;&amp`
	atUserName, atUserOpenId, issueTitle = MatchUserInsUserNameWithIssueTitle(text)
	t.Log(atUserOpenId, atUserName, issueTitle)
	{
		text = `<at open_id="ou_83776f77689ab1a5f399d5aafa27a0a4">@苏汉13</at> <br标题1`
		atUserName, atUserOpenId, issueTitle = MatchUserInsUserNameWithIssueTitle(text)
		t.Log(atUserOpenId, atUserName, issueTitle)
	}
	{
		text = `标题1001<at open_id="ou_83776f77689ab1a5f399d5aafa27a0a4">@苏汉13</at>`
		atUserName, atUserOpenId, issueTitle = MatchUserInsUserNameWithIssueTitle(text)
		t.Log(atUserOpenId, atUserName, issueTitle)
	}
	{
		text = `<at open_id="ou_83776f77689ab1a5f399d5aafa27a0a4">@苏汉13</at> >ddsqewad>`
		atUserName, atUserOpenId, issueTitle = MatchUserInsUserNameWithIssueTitle(text)
		t.Log(atUserOpenId, atUserName, issueTitle)
	}
}
