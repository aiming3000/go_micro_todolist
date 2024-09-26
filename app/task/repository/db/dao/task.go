package dao

import (
	"context"
	"go_micro_todolist/app/task/repository/db/model"
	"go_micro_todolist/idl/pb"
	"gorm.io/gorm"
)

type TaskDao struct {
	*gorm.DB
}

func NewTaskDao(ctx context.Context) *TaskDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &TaskDao{NewDBClient(ctx)}
}

func (dao *TaskDao) CreateTask(in *model.Task) error {
	return dao.Model(&model.Task{}).Create(&in).Error
}

// 通过userId 获取任务列表
func (dao *TaskDao) ListTaskByUserId(userId uint64, start, limit int) (r []*model.Task, count int64, err error) {
	err = dao.Model(&model.Task{}).Offset(start).Limit(limit).
		Where("uid = ?", userId).Find(&r).Error
	err = dao.Model(&model.Task{}).Where("uid = ?", userId).Count(&count).Error
	return
}

func (dao *TaskDao) GetTaskByTaskIdAndUserId(taskId, userId uint64) (r *model.Task, err error) {
	err = dao.Model(&model.Task{}).
		Where("id = ? AND uid = ?", taskId, userId).
		First(&r).Error
	return
}

func (dao *TaskDao) UpdateTask(req *pb.TaskRequest) (r *model.Task, err error) {
	r = new(model.Task)
	err = dao.Model(&model.Task{}).Where("id = ? AND uid = ?", req.Id, req.Uid).
		First(&r).Error
	if err != nil {
		return
	}
	r.Title = req.Title
	r.Status = int(req.Status)
	r.Content = req.Content

	err = dao.Save(&r).Error
	return
}

func (dao *TaskDao) DeleteTaskByIdAndUserId(taskId, userId uint64) (err error) {
	return dao.Model(&model.Task{}).
		Where("id =? AND uid=?", taskId, userId).
		Delete(&model.Task{}).Error
}
