package callsvc

import (
	"fmt"
	"io/ioutil"

	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/model/vo/orgvo"

	"github.com/gin-gonic/gin"
)

type CardCallBack struct {
	SpaceId    string `json:"spaceId"`
	Extension  string `json:"extension"`
	CorpId     string `json:"corpId"`
	SpaceType  string `json:"spaceType"`
	UserIdType string `json:"userIdType"`
	OutTrackId string `json:"outTrackId"`
	Type       string `json:"type"`
	UserId     string `json:"userId"`
	Value      string `json:"value"`
	Content    string `json:"content"`
}

type ValueData struct {
	CardPrivateData *CardPrivateData `json:"cardPrivateData"`
}

type CardPrivateData struct {
	ActionIds []string          `json:"actionIds"`
	Params    map[string]string `json:"params"`
}

func DingTalkCardHandler(c *gin.Context) {
	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	log.Infof("card callback:%v, err :%v", string(body), err)

	clickData := &CardCallBack{}
	jsonErr := json.Unmarshal(body, &clickData)
	if jsonErr != nil {
		log.Errorf("[DingTalkCardHandler] Unmarshal err:%v", jsonErr)
		return
	}
	valueData := &ValueData{}
	jsonErr = json.FromJson(clickData.Value, valueData)
	if jsonErr != nil {
		log.Errorf("[DingTalkCardHandler] Unmarshal err:%v", jsonErr)
		return
	}
	if valueData.CardPrivateData != nil && valueData.CardPrivateData.Params["openConversationId"] != "" {
		cardDataResp := orgfacade.GetUpdateCoolAppTopCardData(orgvo.GetCoolAppTopCardDataReq{Input: &orgvo.CoolAppBaseData{OpenConversationId: valueData.CardPrivateData.Params["openConversationId"]}})
		if cardDataResp.Failure() {
			log.Errorf("[DingTalkCardHandler] GetUpdateCoolAppTopCardData err:%v", cardDataResp.Error())
			return
		}

		fmt.Fprintln(c.Writer, cardDataResp.Data)
	}
	//fmt.Fprintln(c.Writer, `{"cardData":{"cardParamMap":{"today":"1","needMeConfirm": "2","mine":"3","myLate":"4","updateTime": "2/27 17:00"}}}`)
}
