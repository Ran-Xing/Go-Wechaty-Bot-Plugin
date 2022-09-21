package Cron

import (
	"fmt"
	. "github.com/XRSec/Go-Wechaty-Bot/General"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	. "github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"time"
)

var (
	err error
	v   int
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()

	plug.OnLogin(func(context *Context, user *user.ContactSelf) {
		if v != 0 {
			return
		}
		nyc, _ := time.LoadLocation("Asia/Shanghai")
		c := cron.New(cron.WithLocation(nyc))

		if _, err = c.AddFunc("0 0 * * *", func() {
			fmt.Println("========================别卷了提醒👇========================")
			var roomID _interface.IRoom
			// 得不到 roomID 所以出问题了
			if roomID = NewGlobleService().GetBot().Room().Find(&schemas.RoomQueryFilter{Id: "25436935928@chatroom"}); roomID == nil {
				DingSend(viper.GetString("Bot.AdminID"), "RoomID Find Error")
				log.Infof("RoomID Find Error")
				return
			}
			// 发送消息
			if _, err = roomID.Say("@\u2005Qian别卷了,赶紧睡觉!"); err != nil {
				DingSend(viper.GetString("Bot.AdminID"), "failed to send messages")
				log.Errorf("Bypass Say Error: [%v]", err)
				return
			}
			log.Infof("别卷了提醒 Say Success")
		}); err != nil {
			DingSend(viper.GetString("Bot.AdminID"), "Bypass Cron Add Error: "+err.Error())
			log.Errorf("Bypass Cron Add Error: [%v]", err)
		}
		log.Infof("Cron Start")
		// 启动定时器
		c.Start()
		v += 1
	})

	return plug
}
