package callsvc

import (
	"regexp"
	"strings"

	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util/json"
)

type GroupChatUserAtBotHandler struct{}

func (GroupChatUserAtBotHandler) Handle(data string) (string, errs.SystemErrorInfo) {
	reqData := &req_vo.UserAtBotReq{}
	_ = json.FromJson(data, reqData)
	log.Infof("用户 @Bot 请求处理：%s", data)
	insReg := NewUserInstructionReg()
	insReg.HandleInstruction(&reqData.Event)

	return "", nil
}

/// 对处理用户的指令做一层抽象
/// 注册指令处理对象，其中包含指令，以及对应的 handler

// / 指令注册器。Ins 代表 Instruction
type UserInstructionReg struct {
	InsHandlers map[string]*UserInsIf
}

// todo
func NewUserInstructionReg() *UserInstructionReg {
	return &UserInstructionReg{}
}

func (reg *UserInstructionReg) AddIns(flag string, ins *UserInsIf) {
	reg.InsHandlers[flag] = ins
}

func (reg *UserInstructionReg) GetIns(flag string) *UserInsIf {
	var ins *UserInsIf
	if tmpIns, ok := reg.InsHandlers[flag]; ok {
		ins = tmpIns
	}
	return ins
}

// / 通过用户输入的文本，处理指令。
// / 通过用户 at 机器人的文本，判断要执行的指令（操作/行为），并对指令进行对应的回复处理。
func (reg *UserInstructionReg) HandleInstruction(reqEvent *req_vo.UserAtBotReqEvent) errs.SystemErrorInfo {
	var insText string
	var err errs.SystemErrorInfo
	// 通过文本解析出指令
	insText = reg.ParseInsText(reqEvent.TextWithoutAtBot)
	switch insText {
	case "UserInsProIssue":
		insHandle := instruction.NewUserInsProIssue()
		err = insHandle.Handler(reqEvent)
	}
	return err
}

func (reg *UserInstructionReg) ParseInsText(text string) string {
	// 指令名
	insText := ""
	text = strings.Trim(text, " ")
	regExp1 := regexp.MustCompile(`(项目任务|项目进展|项目动态设置)$`)
	matchRes := regExp1.FindString(text)
	// 1.固定的指令 `项目任务|项目进展|项目动态设置`
	if len(matchRes) > 0 {
		switch matchRes {
		case "项目任务":
			insText = "UserInsProIssue"
		case "项目进展":
			insText = "UserInsProProgress"
		}
	} else {
		// 2.@负责人姓名
		// 3.@负责人姓名 任务标题
		tiles := strings.Split(text, " ")
		if len(tiles) == 1 {

		} else if len(tiles) == 2 {

		} else {
			// 不支持的指令
			insText = ""
		}
	}
	return insText
}

// / 用户指令抽象
type UserInsIf interface {
	GetInsName() string
	Handler(reqEvent *req_vo.UserAtBotReqEvent) errs.SystemErrorInfo
}
