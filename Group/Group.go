package Group

import (
	"fmt"
	"strings"
	"time"
	. "github.com/XRSec/Go-Wechaty-Bot/General"
	. "github.com/XRSec/Go-Wechaty-Bot/Plug"

	"github.com/spf13/viper"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"

	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var err error

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.
		OnMessage(onMessage)
	//OnMessage(onMessage).
	//OnRoomJoin(onRoomJoin)
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
	if !m.AtMe {
		log.Infof("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !m.Status {
		log.Infof("Status: [%v] CoptRight: [%v]", m.Status, Copyright(make([]uintptr, 1)))
		return
	}
	if message.Type() != schemas.MessageTypeText {
		log.Infof("Type: [%v] CoptRight: [%v]", message.Type().String(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Self() || m.UserID == viper.GetString("BOT.ADMINID") {
	} else {
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Infof("Age: [%v] CoptRight: [%v]", message.Age()/(60*time.Second), Copyright(make([]uintptr, 1)))
		return
	}
	if !strings.Contains(message.Text(), "group") {
		log.Infof("Msg not Group: [%v]", message.Text())
		return
	}
	if strings.Contains(message.MentionText(), "group.pass.add") {
		if !viper.GetBool("GROUP.PASS." + m.RoomID) {
			viper.Set("GROUP.PASS."+m.RoomID, true)
		}
		SayMessage(context, message, "群服务关闭!")
	}
	if strings.Contains(message.MentionText(), "group.pass.remove") {
		if viper.GetBool("GROUP.PASS." + m.RoomID) {
			viper.Set("GROUP.PASS."+m.RoomID, false)
		}
		SayMessage(context, message, fmt.Sprintf("群服务开启!"))
	}
	if strings.Contains(message.MentionText(), "group.pass.clean") {
		viper.Set("GROUP.PASS.", "")
		SayMessage(context, message, fmt.Sprintf("群服务已全部开启!"))
	}

	if strings.Contains(message.MentionText(), "group.lists.add") {
		if viper.GetString("GROUP.LISTS."+m.RoomID) != strings.ToLower(m.RoomName) {
			viper.Set("GROUP.LISTS."+m.RoomID, m.RoomName)
		}
		SayMessage(context, message, fmt.Sprintf("智能加群已开启, 关键词: 「%v」", strings.ToLower(m.RoomName)))
	}
	if strings.Contains(message.MentionText(), "group.lists.remove") {
		if viper.GetString("GROUP.LISTS."+m.RoomID) != strings.ToLower(m.RoomName) {
			viper.Set("GROUP.LISTS."+m.RoomID, m.RoomName)
		}
		SayMessage(context, message, fmt.Sprintf("智能加群已开启,关键词: 「%v」", strings.ToLower(m.RoomName)))
	}
	if strings.Contains(message.MentionText(), "group.lists.clean") {
		viper.Set("GROUP.LISTS."+m.RoomID, "")
		SayMessage(context, message, fmt.Sprintf("智能加群已关闭!"))
	}

	if viper.GetBool("GROUP.PASS." + m.RoomID) {
		log.Printf("GROUP Pass, [%v] CoptRight: [%v]", message.Talker().Name(), Copyright(make([]uintptr, 1)))
		m.PassResult = "GROUPPASS"
		m.Pass = true
		context.SetData("msgInfo", m)
		return
	}
}

/*
	进入房间监听回调 room-群聊 inviteeList-受邀者名单 inviter-邀请者
	判断配置项群组id数组中是否存在该群聊id
*/
func onRoomJoin(context *wechaty.Context, room *user.Room, inviteeList []_interface.IContact, inviter _interface.IContact, date time.Time) {
	fmt.Println("========================onRoomJoin👇========================")
	newUser := inviteeList[0].Name()
	if inviteeList[0].Self() {
		log.Infof("机器人加入群聊, 群聊名称:[%v] ,邀请人: [%v], 时间: [%v]", room.Topic(), inviter.Name(), date)
		if _, err = room.Say(fmt.Sprintf("大家好呀.我是%v, 以后请多多关照!", newUser)); err != nil {
			log.Errorf("[onRoomJoin] 加入群聊自我介绍消息发送失败, Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
			return
		} else {
			log.Infof("[onRoomJoin] 加入群聊自我介绍消息发送成功")
			return
		}
	}
	log.Infof("群聊名称: [%v], 新人: [%v], 邀请人: [%v], 时间: [%v]", room.Topic(), newUser, inviter.Name(), date)
	//if !Plug.NightMode(inviter.ID()) {
	//	return
	//}
	welcomeString := fmt.Sprintf("@%v\u2005欢迎新人!", newUser)

	if room.ID() == "24633623445@chatroom" {
		welcomeString = fmt.Sprintf("@%v\u2005欢迎加入数藏手动党交流群，请仔细阅读群公告📢", newUser)
	}

	if _, err = room.Say(welcomeString); err != nil {
		log.Errorf("[onRoomJoin] 欢迎新人加入群聊消息发送失败, Error: %v CoptRight: [%s]", err, Copyright(make([]uintptr, 1)))
	} else {
		log.Infof("[onRoomJoin] 欢迎新人加入群聊消息发送成功")
	}
}
