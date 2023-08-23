package ugserver

import (
	"context"
	"user_growth/models"
	"user_growth/pb"
	"user_growth/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UgCoinServer struct {
	pb.UnimplementedUserCoinServer
}

// 获取所有的积分任务列表
func (s *UgCoinServer) ListTasks(ctx context.Context, in *pb.ListTasksRequest) (*pb.ListTasksReply, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method ListTasks not implemented")
	coinTaskSvc := service.NewCoinTaskService(ctx)
	datalist, err := coinTaskSvc.FindAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	dlist := make([]*pb.TbCoinTask, 0, len(datalist))
	for _, v := range datalist {
		dlist = append(dlist, models.CoinTaskToMessage(&v))
	}
	out := &pb.ListTasksReply{Datalist: dlist}
	return out, nil
}

// 获取用户的积分信息
func (s *UgCoinServer) UserCoinInfo(ctx context.Context, in *pb.UserCoinInfoRequest) (*pb.UserCoinInfoReply, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method UserCoinInfo not implemented")
	coinUserSvc := service.NewCoinUserService(ctx)
	uid := int(in.Uid)
	data, err := coinUserSvc.GetByUid(uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	d := models.CoinUserToMessage(data)
	out := &pb.UserCoinInfoReply{Data: d}
	return out, nil
}

// 获取用户的积分明细列表
func (s *UgCoinServer) UserDetails(ctx context.Context, in *pb.UserDetailsRequest) (*pb.UserDetailsReply, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method UserDetails not implemented")
	uid := int(in.Uid)
	page := int(in.Page)
	size := int(in.Size)
	coinDetailSvc := service.NewCoinDetailService(ctx)
	datalist, total, err := coinDetailSvc.FindByUid(uid, page, size)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	dlist := make([]*pb.TbCoinDetail, 0, len(datalist))
	for _, v := range datalist {
		dlist = append(dlist, models.CoinDetailToMessage(&v))
	}
	return &pb.UserDetailsReply{
		Datalist: dlist,
		Total:    int32(total),
	}, nil
}

// 调整用户积分-奖励和惩罚都是用这个接口
func (s *UgCoinServer) UserCoinChange(ctx context.Context, in *pb.UserCoinChangeRequest) (*pb.UserCoinChangeReply, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method UserCoinChange not implemented")
	uid := int(in.Uid)
	task := in.Task
	coin := int(in.Coin)
	// 先查询任务是否存在
	taskInfo, err := service.NewCoinTaskService(ctx).GetByTask(task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if taskInfo == nil {
		return nil, status.Errorf(codes.NotFound, "任务不存在")
	}
	if coin == 0 {
		coin = taskInfo.Coin
	}
	// 插入积分详情
	coinDetail := models.TbCoinDetail{
		Uid:    uid,
		TaskId: int(taskInfo.Id),
		Coin:   coin,
	}
	err = service.NewCoinDetailService(ctx).Save(&coinDetail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	// 更新用户信息
	coinUserSvc := service.NewCoinUserService(ctx)
	coinUser, err := coinUserSvc.GetByUid(uid)
	if err != nil {
		return nil, err
	}
	if coinUser == nil {
		coinUser = &models.TbCoinUser{
			Uid:   uid,
			Coins: coin,
		}
	} else {
		coinUser.Coins += coin
		// 时间置为空，防止再次写入
		coinUser.SysCreated = nil
		coinUser.SysUpdated = nil
	}
	err = coinUserSvc.Save(coinUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UserCoinChangeReply{User: models.CoinUserToMessage(coinUser)}, nil
}
