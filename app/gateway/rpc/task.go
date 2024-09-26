package rpc

import (
	"context"
	"fmt"

	"go_micro_todolist/idl/pb"
	"go_micro_todolist/pkg/e"
)

func TaskCreate(ctx context.Context, req *pb.TaskRequest) (resp *pb.TaskDetailResponse, err error) {
	r, err := TaskService.CreateTask(ctx, req)
	if err != nil {
		return
	}
	if r.Code != e.SUCCESS {
		return
	}

	return r, nil
}

func TaskUpdate(ctx context.Context, req *pb.TaskRequest) (resp *pb.TaskDetailResponse, err error) {
	fmt.Println(22222)
	r, err := TaskService.UpdateTask(ctx, req)
	fmt.Println(err)
	if err != nil {
		return
	}
	if r.Code != e.SUCCESS {
		return
	}

	return r, nil
}

func TaskDelete(ctx context.Context, req *pb.TaskRequest) (resp *pb.TaskDetailResponse, err error) {
	r, err := TaskService.DeleteTask(ctx, req)
	if err != nil {
		return
	}
	if r.Code != e.SUCCESS {
		return
	}

	return r, nil
}

func TaskList(ctx context.Context, req *pb.TaskRequest) (resp *pb.TaskListResponse, err error) {
	r, err := TaskService.GetTasksList(ctx, req)
	if err != nil {
		return
	}
	if r.Code != e.SUCCESS {
		return
	}

	return r, nil
}

func TaskGet(ctx context.Context, req *pb.TaskRequest) (resp *pb.TaskDetailResponse, err error) {
	r, err := TaskService.GetTask(ctx, req)
	if err != nil {
		return
	}
	if r.Code != e.SUCCESS {
		return
	}

	return r, nil
}
