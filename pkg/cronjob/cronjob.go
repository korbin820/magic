package cronjob

import (
	"fmt"
	"magic/basic/logs"
	"time"

	"github.com/robfig/cron"
)

// 定时任务demo
func Demo() {
	c := cron.New()

	// 每5秒钟执行(格式: 秒 分 时
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		logs.DefaultConsoleLog.Info("cronjob Demo", "开始执行!")
		start := time.Now()

		elapsed := time.Since(start)
		logs.DefaultConsoleLog.Info("cronjob Demo", "执行结束!")
		fmt.Println("cronjob Demo耗时: ", elapsed)
	})
	logs.DefaultConsoleLog.Info("cronjob Demo", "启用定时任务!")
	c.Start()
}
