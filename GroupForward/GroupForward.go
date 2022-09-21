package GroupForward

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	. "github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"io/ioutil"
	"os"
	"strings"
	"time"
	. "github.com/XRSec/Go-Wechaty-Bot/General"
)

/*
 目前已知问题：
 1. 微信默认的服务号（文件传输助手）会识别成好友
 2. 消息频率 风控问题
*/
var (
	listResult []_interface.IContact
	listsAll   []_interface.IContact
	listsAllID map[string]bool
	listsExist map[string]bool
	err        error
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(func(context *Context, message *user.Message) {
		if message.Self() || message.Talker().ID() == viper.GetString("BOT.ADMINID") {
		} else {
			return
		}
		if message.Type() != schemas.MessageTypeText {
			return
		}
		if !strings.Contains(message.Text(), "节日祝福") {
			return
		}
		if message.Text()[0:13] == "节日祝福 " {
			if _, err = os.Stat("friendExist.txt"); err != nil {
				getAllToFile(message.GetWechaty().Contact())
				listResult = listsAll
				SayMessage(context, message, fmt.Sprintf("群发开始,共计%v人", len(listResult)))
			} else {
				getAllToFile(message.GetWechaty().Contact())
				readFromFile()
				listResult = DiffArray(listsAll, listsExist)
				SayMessage(context, message, fmt.Sprintf("群发继续, 剩余%v人", len(listResult)))
			}

			//if msg == "" {
			//	if message.Text()[0:7] == "群发 " {
			//		msg = message.Text()[7:]
			//	}
			//	if message.Text()[0:8] == "forward " {
			//		msg = message.Text()[8:]
			//	}
			//}
			for i := 1; i < len(listResult); i++ {
				msg := message.Text()[13:]
				if strings.Contains(message.Text(), "%v") {
					msg = fmt.Sprintf(message.Text()[13:], listResult[i].Name())
				}
				if _, err = listResult[i].Say(fmt.Sprintf(msg)); err != nil {
					SayMessage(context, message, fmt.Sprintf("群发失败, 剩余%v人未发送成功", len(listResult)-i))
					writeToFile("friendExist.txt", listsExist)
					return
				}
				listsExist[listResult[i].ID()] = true
				writeToFile("friendExist.txt", listsExist)
				time.Sleep(time.Second * 8)
			}
			SayMessage(context, message, "群发完毕")
			//if err := os.Remove("friend.json"); err != nil {
			//	log.Errorf("os.Remove Error: [%v]", err)
			//	return
			//}
		}
		if message.Text()[0:19] == "节日祝福测试 " {
			SayMessage(context, message, fmt.Sprintf("嗨, 亲爱的%v, %v", message.Talker().Name(), message.Text()[19:]))
		}
	})
	//Do(plug.Wechaty)
	return plug
}

func getAllToFile(c _interface.IContactFactory) {
	listsAllID = make(map[string]bool)
	for _, v := range c.FindAll(nil) {
		if v.Type() != schemas.ContactTypePersonal {
			continue
		}
		if !v.Friend() {
			continue
		}
		listsAll = append(listsAll, v)
		listsAllID[v.ID()] = false
	}
	log.Infof("ContactList: 加载成功, 检测到 %v 位好友", len(listsAll))
}

func readFromFile() {
	result, err := ioutil.ReadFile("friendExist.txt")
	if err != nil {
		log.Errorf("ioutil.ReadFile Error: [%v]", err)
		return
	}
	listsExist = make(map[string]bool)
	err = json.Unmarshal(result, &listsExist)
	if err != nil {
		log.Errorf("json.Unmarshal Error: [%v]", err)
		return
	}
}

func DiffArray(a []_interface.IContact, b map[string]bool) (c []_interface.IContact) {
	for _, v := range a {
		if !b[v.ID()] {
			c = append(c, v)
		}
	}
	return c
}

func writeToFile(fileName string, v map[string]bool) {
	result, err := json.Marshal(v)
	if err != nil {
		log.Errorf("json.Marshal Error: [%v]", err)
		return
	}
	_ = ioutil.WriteFile(fileName, result, 0644)
}
