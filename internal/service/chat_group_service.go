package service

import (
	"context"
	"errors"
	"realtime-chat/internal/dto/request"
	"realtime-chat/internal/dto/response"
	"realtime-chat/internal/model"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/repository/mysql"
)

type ChatGroupService interface {
	Create(ctx context.Context, userID int64, req *request.CreateGroup) error      // done
	AddMember(ctx context.Context, userID, req *request.AddMemberToGroup) error    // done
	Detail(ctx context.Context, userID, roomID int64) (*response.ChatGroup, error) // done
	Delete(ctx context.Context, userID, roomID int64) error
}

type chatGroupServiceImpl struct {
	roomRepo   mysql.RoomRepository
	memberRepo mysql.RoomMemberRepository
}

// add member to group
// - cek room id ada atau tidak
// - cek apakah userid pada room_member adalah admin

func (cgSvc *chatGroupServiceImpl) AddMember(ctx context.Context, userID int64, req *request.AddMemberToGroup) error {
	exist, err := cgSvc.roomRepo.IsExist(ctx, req.RoomID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("not found group")
	}

	memberResult, err := cgSvc.memberRepo.FindByUserIDAndRoomID(ctx, userID, req.RoomID)
	if err != nil {
		return err
	}

	if memberResult.Role != pkg.RoleAdmin {
		return errors.New("anda bukan admin")
	}

	if err := cgSvc.memberRepo.AddMembers(ctx, req.MemberUserIDs); err != nil {
		return err
	}
	return nil
}

func (cgSvc *chatGroupServiceImpl) Create(ctx context.Context, userID int64, req *request.CreateGroup) error {
	admin := model.RoomMemberModel{
		Role:   pkg.RoleAdmin,
		UserID: userID,
	}

	// belum bikin gambarnya
	room := model.RoomModel{
		Type: pkg.RoomTypeGroup,
		Name: req.Name,
		Members: []model.RoomMemberModel{
			admin,
		},
	}

	if err := cgSvc.roomRepo.Create(ctx, &room); err != nil {
		return err
	}
	return nil
}

func (cgSvc *chatGroupServiceImpl) Detail(ctx context.Context, userID, roomID int64) (*response.ChatGroup, error) {
	groupResult, err := cgSvc.roomRepo.FindGroupByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	members := make([]response.GroupMember, 0)
	for _, member := range groupResult.Members {
		members = append(members, response.GroupMember{
			MemberID: member.ID,
			UserID:   member.UserID,
			UserName: member.User.UserName,
			Role:     member.Role,
		})
	}

	return &response.ChatGroup{
		RoomID:    groupResult.ID,
		Type:      groupResult.Type,
		Name:      groupResult.Name,
		ImagePath: pkg.ImagePath + groupResult.Image,
		Members:   members,
		CreatedAt: groupResult.CreatedAt.Format(pkg.DateTimeLayout),
		UpdatedAt: groupResult.UpdatedAt.Format(pkg.DateTimeLayout),
	}, nil
}

func (cgSvc *chatGroupServiceImpl) Delete(ctx context.Context, userID, roomID int64) error {
	memberResult, err := cgSvc.memberRepo.FindByUserIDAndRoomID(ctx, userID, roomID)
	if err != nil {
		return err
	}

	if memberResult.Role != pkg.RoleAdmin {
		return errors.New("anda tidak diizinkan hapus group")
	}

	if err := cgSvc.roomRepo.DeleteByID(ctx, roomID); err != nil {
		return err
	}

	return nil
}
