package rsi

import (
	"mcp_service/pkg/cron"
)

func StartRsiTask() {
	cronTask := cron.CronTask{}
	for _, symbol := range SymbolList {
		//添加定时任务
		for interval, cronExpression := range IntervalList {
			cronTask.CronExpression = cronExpression
			cronTask.Task = func() {
				GetKline(symbol, interval)
			}
			cronTask.AddCronTask()

		}
	}
	cronTask.Start()
}
