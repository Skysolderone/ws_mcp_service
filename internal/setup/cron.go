package setup

import (
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

var CronTimer *cron.Cron

type CronTask struct {
	CronID         cron.EntryID
	CronExpression string
	Task           func()
}

func InitCronTimer() {
	CronTimer = cron.New()
}

func (c *CronTask) AddCronTask() {
	var err error
	c.CronID, err = CronTimer.AddFunc(c.CronExpression, c.Task)
	if err != nil {
		logx.Errorw("添加定时任务失败", logx.Field("错误", err.Error()), logx.Field("定时任务ID", c.CronID), logx.Field("定时任务表达式", c.CronExpression))
		return
	}
	logx.Infow("添加定时任务成功", logx.Field("定时任务ID", c.CronID), logx.Field("定时任务表达式", c.CronExpression))
}

func (c *CronTask) Start() {
	CronTimer.Start()
}

func (c *CronTask) Stop() {
	CronTimer.Stop()
}
