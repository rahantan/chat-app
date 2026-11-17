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
)

type RoomService interface {
	// ListMyRooms(ctx context.Context, userID int64) ([]dto.RoomResponse, error)

	// GetRoomDetail(ctx context.Context, userID, roomID int64) (*dto.RoomDetailResponse, error)

	// CreateGroupRoom(ctx context.Context, userID int64, req dto.CreateGroupRoomRequest) (*dto.RoomResponse, error)

	ListMyRooms(ctx context.Context, userID int64) ([]response.Room, error)                           // done
	GetRoomDetail(ctx context.Context, userID, roomID int64) (*response.Room, error)                  // done
	CreateGroupRoom(ctx context.Context, userID int64, req *request.CreateGroup) error                // done
	AddMembersToGroup(ctx context.Context, userID, roomID int64, req *request.AddMemberToGroup) error // done

	GetOrCreatePrivateRoom(ctx context.Context, userID int64, req *request.GetOrCreatePrivateRoom) (int64, error)
	// DeleteGroupRoom(ctx context.Context, userID, roomID int64) error // BELUM

	// GetOrCreatePrivateRoom(ctx context.Context, userID, targetUserID int64) (*model.RoomModel, error)
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
	roomResults, err := roomSvc.roomRepo.FindAll(ctx, userID)
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

	roomResult, err := roomSvc.roomRepo.FindDetailByID(ctx, roomID)
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
func (roomSvc *roomServiceImpl) AddMembersToGroup(ctx context.Context, userID, roomID int64, req *request.AddMemberToGroup) error {
	_, err := roomSvc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return err
	}

	userIDs := append([]int64{userID}, req.UserIDs...)
	memberResults, err := roomSvc.memberRepo.FindByRoomIDAndUserIDs(ctx, roomID, userIDs)
	if err != nil {
		return err
	}

	membersMap := make(map[int64]model.RoomMemberModel)
	for _, member := range memberResults {
		membersMap[member.UserID] = member
	}

	if owner, ok := membersMap[userID]; ok {
		if owner.Role != pkg.GroupRoleAdmin {
			return pkg.ErrForbidden
		}
	} else {
		return pkg.ErrForbidden
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

	if err := roomSvc.memberRepo.AddMembers(ctx, roomID, req.UserIDs); err != nil {
		return err
	}

	return nil
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

	members := []model.RoomMemberModel{
		{UserID: userID},
		{UserID: req.TargetUserID},
	}

	room := model.RoomModel{
		Type:    pkg.RoomTypePrivate,
		Members: members,
	}

	fmt.Println("MASUK")
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
