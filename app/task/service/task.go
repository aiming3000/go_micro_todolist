package service

import (
	"context"
	"encoding/json"
	"go_micro_todolist/app/task/repository/db/dao"
	"go_micro_todolist/app/task/repository/db/model"
	"go_micro_todolist/app/task/repository/mq"
	"go_micro_todolist/idl/pb"
	"go_micro_todolist/pkg/e"
	log "go_micro_todolist/pkg/logger"
	"sync"
)

var TaskSrvIns *TaskSrv
var TaskSrvOnce sync.Once

type TaskSrv struct {
}

func GetTaskSrv() *TaskSrv {
	TaskSrvOnce.Do(func() {
		TaskSrvIns = &TaskSrv{}
	})
	return TaskSrvIns
}

// CreateTask 创建备忘录，将备忘录信息生产，放到rabbitMQ消息队列中
func (t *TaskSrv) CreateTask(ctx context.Context, in *pb.TaskRequest, out *pb.TaskDetailResponse) error {
	body, _ := json.Marshal(in) //title,content
	out.Code = e.SUCCESS
	err := mq.SendMessage2MQ(body)
	if err != nil {
		out.Code = e.ERROR
		return nil
	}
	return nil
}

func TaskMQ2MySQL(ctx context.Context, req *pb.TaskRequest) error {
	m := &model.Task{
		Uid:       uint(req.Uid),
		Title:     req.Title,
		Status:    int(req.Status),
		Content:   req.Content,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	return dao.NewTaskDao(ctx).CreateTask(m)
}

func (t *TaskSrv) GetTasksList(ctx context.Context, in *pb.TaskRequest, out *pb.TaskListResponse) error {
	out.Code = e.SUCCESS
	if in.Limit == 0 {
		in.Limit = 10
	}
	//查找备忘录
	r, count, err := dao.NewTaskDao(ctx).ListTaskByUserId(in.Uid, int(in.Start), int(in.Limit))
	if err != nil {
		out.Code = e.ERROR
		log.LogrusObj.Error("ListTaskByUserId err:%v", err)
		return err
	}
	// 返回proto里面定义的类型
	var taskRes []*pb.TaskModel
	for _, item := range r {
		taskRes = append(taskRes, BuildTask(item))
	}
	out.TaskList = taskRes
	out.Count = uint32(count)
	return nil
}

// 获取备忘录详情
func (t *TaskSrv) GetTask(ctx context.Context, in *pb.TaskRequest, out *pb.TaskDetailResponse) error {
	out.Code = e.SUCCESS
	r, err := dao.NewTaskDao(ctx).GetTaskByTaskIdAndUserId(in.Id, in.Uid)
	if err != nil {
		out.Code = e.ERROR
		log.LogrusObj.Error("GetTask err:%v", err)
		return err
	}
	taskRes := BuildTask(r)
	out.TaskDetail = taskRes
	return nil
}

func (t *TaskSrv) UpdateTask(ctx context.Context, in *pb.TaskRequest, out *pb.TaskDetailResponse) error {
	out.Code = e.SUCCESS
	r, err := dao.NewTaskDao(ctx).UpdateTask(in)
	if err != nil {
		out.Code = e.ERROR
		log.LogrusObj.Error("UpdateTask err:%v", err)
		return err
	}
	out.TaskDetail = BuildTask(r)

	return nil
}

func (t *TaskSrv) DeleteTask(ctx context.Context, in *pb.TaskRequest, out *pb.TaskDetailResponse) error {
	out.Code = e.SUCCESS
	err := dao.NewTaskDao(ctx).DeleteTaskByIdAndUserId(in.Id, in.Uid)
	if err != nil {
		out.Code = e.ERROR
		log.LogrusObj.Error("UpdateTask err:%v", err)
		return err
	}

	return nil
}

func BuildTask(item *model.Task) *pb.TaskModel {
	taskModel := pb.TaskModel{
		Id:         uint64(item.ID),
		Uid:        uint64(item.Uid),
		Title:      item.Title,
		Content:    item.Content,
		StartTime:  item.StartTime,
		EndTime:    item.EndTime,
		Status:     int64(item.Status),
		CreateTime: item.CreatedAt.Unix(),
		UpdateTime: item.UpdatedAt.Unix(),
	}
	return &taskModel
}
