package CodeAuth

import (
	"fmt"
	. "github.com/XRSec/Go-Wechaty-Bot/General"
	. "github.com/XRSec/Go-Wechaty-Bot/Plug"
	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"strings"
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !strings.Contains(message.MentionText(), "auth") {
		return
	}
	parts := strings.Split(message.Text(), " ")
	if len(parts) != 3 {
		return
	}

	if parts[1] == "" {
		return
	}

	if len(parts[2]) != 6 {
		return
	}

	log.Debugln("查找微信")
	msg := ""
	var userIndex = -1
	userInfo := message.GetWechaty().Contact().FindAll("初见.")
	for i := 0; i < len(userInfo); i++ {
		tempID := userInfo[i].ID()
		if tempID == "" {
			continue
		}
		if userIndex != -1 {
			msg = fmt.Sprintf("存在多个微信为: %s 的用户", parts[1])
			goto end
		}
		userIndex = i
	}

	if _, err := userInfo[userIndex].Say(parts[2]); err != nil {
		msg = fmt.Sprintf("发送失败: %s", err.Error())
		goto end
	}
	msg = "发送成功"
	goto end

end:
	SayMessage(context, message, msg)
	m.ReplyResult = msg
	m.Pass = true
	m.PassResult = "登录认证"
	context.SetData("msgInfo", m)
}
