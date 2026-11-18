package service

import (
	"context"
	"errors"
	"fmt"
	"realtime-chat/internal/dto/request"
	"realtime-chat/internal/dto/response"
	"realtime-chat/internal/model"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/repository/mysql"
	"slices"
)

type RoomService interface {
	ListMyRooms(ctx context.Context, userID int64) ([]response.Room, error)          // done
	GetRoomDetail(ctx context.Context, userID, roomID int64) (*response.Room, error) // done

	GetOrCreatePrivateRoom(ctx context.Context, userID int64, req *request.GetOrCreatePrivateRoom) (int64, error) // done

	AddMembersToGroup(ctx context.Context, userID, roomID int64, req *request.AddMembersToGroup) error           // done
	CreateGroupRoom(ctx context.Context, userID int64, req *request.CreateGroup) error                           // done
	RemoveMembersFromGroup(ctx context.Context, userID, roomID int64, req *request.RemoveMembersFromGroup) error // done
	DeleteGroupRoom(ctx context.Context, userID, roomID int64) error                                             // done
}

type roomServiceImpl struct {
	roomRepo   mysql.RoomRepository
	memberRepo mysql.RoomMemberRepository
	userRepo   mysql.UserRepository
}

func NewRoomService(room mysql.RoomRepository, member mysql.RoomMemberRepository, user mysql.UserRepository) RoomService {
	return &roomServiceImpl{
		roomRepo:   room,
		memberRepo: member,
		userRepo:   user,
	}
}
func (roomSvc *roomServiceImpl) ListMyRooms(ctx context.Context, userID int64) ([]response.Room, error) {
	roomResults, err := roomSvc.roomRepo.FindByUserIDWithDetail(ctx, userID)
	if err != nil {
		return nil, err
	}

	rooms := make([]response.Room, 0, len(roomResults))
	for _, room := range roomResults {

		roomName := room.Name
		roomImage := pkg.GroupImageUrl + room.Image

		if room.Type == pkg.RoomTypePrivate {
			for _, member := range room.Members {
				if member.UserID != userID {
					roomName = member.User.UserName
					roomImage = pkg.AvatarImageUrl + member.User.Avatar
				}
			}
		}

		var lastMsg *response.RoomLastMsg

		if room.LastMessage == nil {
			if room.Type == pkg.RoomTypePrivate {
				continue
			}
			lastMsg = nil
		} else {
			lastMsg = &response.RoomLastMsg{
				UserID:      room.LastMessage.Sender.ID,
				UserName:    room.LastMessage.Sender.UserName,
				Content:     room.LastMessage.Content,
				TypeContent: room.LastMessage.MsgType,
				CreatedAt:   room.LastMessage.CreatedAt.Format(pkg.DateTimeLayout),
			}
		}

		rooms = append(rooms, response.Room{
			RoomID:    room.ID,
			Type:      room.Type,
			Name:      roomName,
			ImageUrl:  roomImage,
			LastMsg:   lastMsg,
			UpdatedAt: room.UpdatedAt.Format(pkg.DateTimeLayout),
		})
	}
	return rooms, nil
}

func (roomSvc *roomServiceImpl) GetRoomDetail(ctx context.Context, userID, roomID int64) (*response.Room, error) {

	roomResult, err := roomSvc.roomRepo.FindByIDWithMembers(ctx, roomID)
	if err != nil {
		return nil, err
	}

	_, err = roomSvc.memberRepo.FindByUserIDAndRoomID(ctx, userID, roomID)
	if err != nil {
		if errors.Is(err, model.ErrMemberNotFound) {
			return nil, pkg.ErrAccessDenied
		}
		return nil, err
	}

	members := make([]response.Member, 0, len(roomResult.Members))

	roomName := roomResult.Name
	roomImage := pkg.GroupImageUrl + roomResult.Image

	for _, member := range roomResult.Members {
		if roomResult.Type == pkg.RoomTypePrivate {
			if member.UserID != userID {
				roomName = member.User.UserName
				roomImage = pkg.AvatarImageUrl + member.User.Avatar
			}
		}

		members = append(members, response.Member{
			UserID:    member.UserID,
			UserName:  member.User.UserName,
			Role:      member.Role,
			AvatarUrl: pkg.AvatarImageUrl + member.User.Avatar,
		})
	}

	return &response.Room{
		RoomID:    roomResult.ID,
		Type:      roomResult.Type,
		Name:      roomName,
		ImageUrl:  roomImage,
		Members:   members,
		UpdatedAt: roomResult.UpdatedAt.Format(pkg.DateTimeLayout),
	}, nil
}

func (roomSvc *roomServiceImpl) GetOrCreatePrivateRoom(ctx context.Context, userID int64, req *request.GetOrCreatePrivateRoom) (int64, error) {
	roomID, err := roomSvc.roomRepo.FindPrivateRoomByUsers(ctx, userID, req.TargetUserID)
	if err == nil {
		return roomID, nil
	}

	if !errors.Is(err, model.ErrRoomNotFound) {
		return 0, err
	}

	_, err = roomSvc.userRepo.FindByID(ctx, req.TargetUserID)
	if err != nil {
		return 0, err
	}

	room := model.RoomModel{
		Type: pkg.RoomTypePrivate,
		Members: []model.RoomMemberModel{
			{UserID: userID},
			{UserID: req.TargetUserID},
		},
	}

	if err := roomSvc.roomRepo.Create(ctx, &room); err != nil {
		if !errors.Is(err, model.ErrDuplicatePrivateRoom) {
			return 0, err
		}

		roomID, findErr := roomSvc.roomRepo.FindPrivateRoomByUsers(ctx, userID, req.TargetUserID)
		if findErr != nil {
			return 0, findErr
		}

		return roomID, nil
	}

	return room.ID, nil
}

// GROUP AREA
func (roomSvc *roomServiceImpl) CreateGroupRoom(ctx context.Context, userID int64, req *request.CreateGroup) error {
	admin := model.RoomMemberModel{
		Role:   pkg.GroupRoleAdmin,
		UserID: userID,
	}

	imagename := ""
	if req.Image == nil {
		imagename = pkg.DefaultGroupImage
	}
	// belum bikin gambarnya
	room := model.RoomModel{
		Type:  pkg.RoomTypeGroup,
		Name:  req.Name,
		Image: imagename,
		Members: []model.RoomMemberModel{
			admin,
		},
	}

	if err := roomSvc.roomRepo.Create(ctx, &room); err != nil {
		return err
	}

	return nil
}

func (roomSvc *roomServiceImpl) validateOnlyAdmin(ctx context.Context, userID int64, roomID int64) error {
	roomResult, err := roomSvc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return err
	}

	if roomResult.Type == pkg.RoomTypePrivate {
		return pkg.ErrRoomNotGroup
	}

	me, err := roomSvc.memberRepo.FindByUserIDAndRoomID(ctx, userID, roomID)
	if err != nil {
		if errors.Is(err, model.ErrMemberNotFound) {
			return pkg.ErrAccessDenied
		}
		return err
	}

	if me.Role != pkg.GroupRoleAdmin {
		return pkg.ErrForbidden
	}
	return nil
}

func (roomSvc *roomServiceImpl) getMembersWithMap(ctx context.Context, roomID int64, userIDs []int64) (map[int64]model.RoomMemberModel, error) {
	memberResults, err := roomSvc.memberRepo.FindByRoomIDAndUserIDs(ctx, roomID, userIDs)
	if err != nil {
		return nil, err
	}

	membersMap := make(map[int64]model.RoomMemberModel)
	for _, member := range memberResults {
		membersMap[member.UserID] = member
	}

	return membersMap, nil
}

func (roomSvc *roomServiceImpl) AddMembersToGroup(ctx context.Context, userID, roomID int64, req *request.AddMembersToGroup) error {

	if err := roomSvc.validateOnlyAdmin(ctx, userID, roomID); err != nil {
		return err
	}

	membersMap, err := roomSvc.getMembersWithMap(ctx, roomID, req.UserIDs)
	if err != nil {
		return err
	}

	memberExists := []int64{}
	for _, id := range req.UserIDs {
		if _, ok := membersMap[id]; ok {
			memberExists = append(memberExists, id)
		}
	}

	if len(memberExists) > 0 {
		return errors.New("salah satu user sudah ada dalam group")
	}

	userResults, err := roomSvc.userRepo.FindByIDs(ctx, req.UserIDs)
	if err != nil {
		return err
	}

	if len(userResults) != len(req.UserIDs) {
		return errors.New("user not found, gunakan user id yang valid")
	}

	// BELUM BUAT CEK USER APAKAH SUDAH ADA DI GROUP ATAU BELUM

	members := make([]model.RoomMemberModel, 0, len(req.UserIDs))
	for _, id := range req.UserIDs {
		members = append(members, model.RoomMemberModel{
			Role:   pkg.GroupRoleMember,
			RoomID: roomID,
			UserID: id,
		})
	}

	if err := roomSvc.memberRepo.AddMembers(ctx, roomID, members); err != nil {
		return err
	}

	return nil
}

func (roomSvc *roomServiceImpl) RemoveMembersFromGroup(ctx context.Context, userID, roomID int64, req *request.RemoveMembersFromGroup) error {
	if slices.Contains(req.UserIDs, userID) {
		return pkg.NewErr(pkg.ErrCodeBadRequest, "you cannot remove yourself, use leave group instead")
	}

	if err := roomSvc.validateOnlyAdmin(ctx, userID, roomID); err != nil {
		return err
	}

	membersMap, err := roomSvc.getMembersWithMap(ctx, roomID, req.UserIDs)
	if err != nil {
		return err
	}

	var missing []int64
	for _, id := range req.UserIDs {
		if _, ok := membersMap[id]; !ok {
			missing = append(missing, id)
		}
	}

	if len(missing) > 0 {
		return pkg.NewErr(pkg.ErrCodeNotFound, fmt.Sprintf("users are not members of this room: %v", missing))
	}

	if err := roomSvc.memberRepo.DeleteMembersGroup(ctx, roomID, req.UserIDs); err != nil {
		return err
	}
	return nil
}

func (roomSvc *roomServiceImpl) DeleteGroupRoom(ctx context.Context, userID, roomID int64) error {
	if err := roomSvc.validateOnlyAdmin(ctx, userID, roomID); err != nil {
		return err
	}

	if err := roomSvc.roomRepo.DeleteByID(ctx, roomID); err != nil {
		return err
	}
	return nil
}
