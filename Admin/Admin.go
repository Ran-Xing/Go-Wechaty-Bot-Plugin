package Admin

/*
该插件仅提供：
	群聊 添加机器人好友
	群聊踢人
	改名字
*/

import (
	"fmt"
	. "github.com/XRSec/Go-Wechaty-Bot/General"
	. "github.com/XRSec/Go-Wechaty-Bot/Plug"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	err error
)

/*
	Admin()
	管理员
*/
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
	if m.Pass {
		log.Infof("Pass CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if m.Reply {
		log.Infof("Reply CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !m.Status {
		log.Infof("Room Status CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	//if !m.AtMe {
	//	log.Infof("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
	//	return
	//}
	if message.Type() != schemas.MessageTypeText {
		log.Infof("Type: [%v] CoptRight: [%v]", message.Type().String(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Infof("Age: [%v] CoptRight: [%v]", message.Age()/(60*time.Second), Copyright(make([]uintptr, 1)))
		return
	}

	if message.Self() || m.UserID == viper.GetString("BOT.ADMINID") {
	} else {
		return
	}

	if message.MentionText() == "add" || message.MentionText() == "加" { // 添加好友
		var (
			addUser = message.MentionList()[0]
			//member  _interface.IContact
		)
		//if member, err = message.Room().Member(addUserName); err != nil && member != nil {
		//	log.Errorf(fmt.Sprintf("搜索用户名ID失败, 用户名: [%v], 用户信息: [%v]", addUserName, member.String()), err)
		//}
		if message.GetWechaty().Contact().Load(addUser.ID()).Friend() {
			log.Infof("用户已经是好友, 用户名: [%v]", addUser.Name())
			SayMessage(context, message, fmt.Sprintf("用户: [%v] 已经是好友了", addUser.Name()))
			m.PassResult = fmt.Sprintf("用户: [%v] 已经是好友了", addUser.Name())
			m.Pass = true
			goto end
		}

		// TODO 目前这个模块有问题
		if err = message.GetWechaty().Friendship().Add(addUser, fmt.Sprintf("你好,我是%v,以后请多多关照!", viper.GetString("Bot.Name"))); err != nil {
			log.Errorf("添加好友失败, 用户名: [%v], Error: [%v]", addUser.Name(), err)
			SayMessage(context, message, fmt.Sprintf("添加好友失败, 用户名: [%v]\n 是不是没开权限哦[旺柴]", addUser.Name()))
			//return
		} else {
			log.Infof("添加好友成功, 用户名: [%v]", addUser.Name())
			SayMessage(context, message, fmt.Sprintf("好友申请发送成功, 用户: [%v]", addUser.Name()))
		}

		m.PassResult = fmt.Sprintf("好友申请, 用户: [%v]", addUser.Name())
		m.Pass = true
		goto end
	}

	if message.MentionText() == "del" || message.MentionText() == "踢" { // 从群聊中移除用户
		var (
			delUser = message.MentionList()[0]
		)
		if err = message.Room().Del(delUser); err != nil {
			log.Errorf("移除用户失败, 用户: [%v], Error: [%v]", delUser.Name(), err)
			SayMessage(context, message, fmt.Sprintf("移除用户失败, 用户: [%v]", delUser.Name()))
			//return
		} else {
			log.Infof("移除用户成功, 用户: [%v]", delUser.Name())
		}

		m.PassResult = fmt.Sprintf("从群聊中移除用户: [%v]", delUser.Name())
		m.Pass = true
		goto end
	}

	if message.MentionText() == "quit" || message.MentionText() == "退" { // 退群
		SayMessage(context, message, "我走了, 拜拜👋🏻, 记得想我哦 [大哭]")
		if err = message.Room().Quit(); err != nil {
			log.Errorf("退出群聊失败, 群聊名称: [%v], Error: [%v]", message.Room().Topic(), err)
			SayMessage(context, message, fmt.Sprintf("退出群聊失败, 群聊名称: [%v]", message.Room().Topic()))
			//return
		} else {
			log.Infof("退出群聊成功, 群聊名称: [%v]", message.Room().Topic())
		}

		m.PassResult = fmt.Sprintf("退出群聊, 群聊名称: [%v]", message.Room().Topic())
		m.Pass = true
		goto end
	}

	if strings.Contains(message.MentionText(), "gmz") {
		var (
			newName = strings.Replace(message.MentionText(), "gmz ", "", 1)
		)

		if err = message.GetPuppet().SetContactSelfName(newName); err != nil {
			log.Errorf("修改用户名失败, Error: [%v]", err)
			SayMessage(context, message, fmt.Sprintf("修改用户名失败, 新用户名: [%v]", newName))
			//return
		} else {
			log.Infof("修改用户名成功, 新用户名: [%v]", newName)
			SayMessage(context, message, fmt.Sprintf("修改用户名成功, 新用户名: [%v]", newName))
		}

		m.PassResult = fmt.Sprintf("改名字: [%v]", newName)
		m.Pass = true
		goto end
	}
end:
	context.SetData("msgInfo", m)
}
