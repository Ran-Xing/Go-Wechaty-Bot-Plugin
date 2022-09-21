package Health

import (
	"fmt"
	. "github.com/XRSec/Go-Wechaty-Bot/General"
	"github.com/XRSec/Go-Wechaty-Bot/Plug/Jinrishici"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	. "github.com/wechaty/go-wechaty/wechaty"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"time"
)

/*
	Health()
	健康监测
*/

var v int

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	nyc, _ := time.LoadLocation("Asia/Shanghai")

	plug.OnLogin(func(context *Context, user *user.ContactSelf) {
		if v != 0 {
			return
		}
		c := cron.New(cron.WithLocation(nyc))

		if _, err := c.AddFunc("0 23 * * *", func() {
			fmt.Println("========================每日一句👇========================")
			var roomID _interface.IRoom
			//if roomID = bot.Room().Find(&schemas.RoomQueryFilter{Id: "roomID@chatroom"}); roomID == nil {
			if roomID = NewGlobleService().GetBot().Room().Find("Debug"); roomID == nil {
				DingSend(viper.GetString("Bot.AdminID"), "RoomID Find Error")
				log.Infof("RoomID Find Error")
				return
			}

			if _, err := roomID.Say(Jinrishici.Do()); err != nil {
				DingSend(viper.GetString("Bot.AdminID"), "failed to send messages")
				log.Errorf("onHeartbeat Say Error: [%v]", err)
				return
			}
			log.Infof("Heartbeat Say Success")
		}); err != nil {
			DingSend(viper.GetString("Bot.AdminID"), "Heartbeat Cron Add Error: "+err.Error())
			log.Errorf("Heartbeat Cron Add Error: [%v]", err)
		}
		log.Infof("Health Cron Start")
		// 启动定时器
		c.Start()
		v += 1
	})
	return plug
}
