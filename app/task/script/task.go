package script

import (
	"context"

	"go_micro_todolist/app/task/repository/mq/task"
	log "go_micro_todolist/pkg/logger"
)

func TaskCreateSync(ctx context.Context) {
	tSync := new(task.SyncTask)
	err := tSync.RunTaskCreate(ctx)
	if err != nil {
		log.LogrusObj.Infof("RunTaskCreate:%s", err)
	}
}
