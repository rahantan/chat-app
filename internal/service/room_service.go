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

type RoomService interface {
	// ListMyRooms(ctx context.Context, userID int64) ([]dto.RoomResponse, error)

	// GetRoomDetail(ctx context.Context, userID, roomID int64) (*dto.RoomDetailResponse, error)

	// CreateGroupRoom(ctx context.Context, userID int64, req dto.CreateGroupRoomRequest) (*dto.RoomResponse, error)

	ListMyRooms(ctx context.Context, userID int64) ([]response.Room, error)                   // done
	GetRoomDetail(ctx context.Context, userID, roomID int64) (*response.Room, error)          // done
	CreateGroupRoom(ctx context.Context, userID int64, req *request.CreateGroup) error        // done
	AddMembersToGroup(ctx context.Context, userID int64, req *request.AddMemberToGroup) error // done

	GetOrCreatePrivateRoom(ctx context.Context, userID, targetUserID int64) (*model.RoomModel, error)
	DeleteGroupRoom(ctx context.Context, userID, roomID int64) error
}

type roomServiceImpl struct {
	roomRepo   mysql.RoomRepository
	memberRepo mysql.RoomMemberRepository
	userRepo   mysql.UserRepository
}

func (roomSvc *roomServiceImpl) ListMyRooms(ctx context.Context, userID int64) ([]response.Room, error) {
	roomResults, err := roomSvc.roomRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	rooms := make([]response.Room, 0, len(roomResults))
	for _, room := range roomResults {

		roomName := room.Name
		roomImage := room.Image

		if room.Type == pkg.RoomTypePrivate {
			for _, member := range room.Members {
				if member.UserID != userID {
					roomName = member.User.UserName
					roomImage = member.User.Image
				}
			}
		}

		var lastMsg response.RoomLastMsg
		if room.LastMessage != nil {
			lastMsg = response.RoomLastMsg{
				SenderName:  room.LastMessage.Sender.UserName,
				Content:     room.LastMessage.Content,
				TypeContent: room.LastMessage.MsgType,
				CreatedAt:   room.LastMessage.CreatedAt.Format(pkg.DateTimeLayout),
			}
		}

		rooms = append(rooms, response.Room{
			RoomID:    room.ID,
			Type:      room.Type,
			Name:      roomName,
			ImagePath: roomImage,
			LastMsg:   lastMsg,
		})
	}
	return rooms, nil
}

func (roomSvc *roomServiceImpl) GetRoomDetail(ctx context.Context, userID, roomID int64) (*response.Room, error) {

	roomResult, err := roomSvc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	_, err = roomSvc.memberRepo.FindByUserIDAndRoomID(ctx, userID, roomID)
	if err != nil {
		if errors.Is(err, model.ErrRoomNotFound) {
			return nil, errors.New("tidak diizinkan")
		}
		return nil, err
	}

	members := make([]response.Member, 0, len(roomResult.Members))

	roomName := roomResult.Name
	roomImage := roomResult.Image

	for _, member := range roomResult.Members {
		if roomResult.Type == pkg.RoomTypePrivate {
			if member.UserID != userID {
				roomName = member.User.UserName
				roomImage = member.User.Image
			}
		}

		members = append(members, response.Member{
			UserID:   member.UserID,
			UserName: member.User.UserName,
			Avatar:   member.User.Image,
		})
	}

	return &response.Room{
		RoomID:    roomResult.ID,
		Type:      roomResult.Type,
		Name:      roomName,
		ImagePath: roomImage,
		Members:   members,
		UpdatedAt: roomResult.UpdatedAt.Format(pkg.DateTimeLayout),
	}, nil
}

func (roomSvc *roomServiceImpl) CreateGroupRoom(ctx context.Context, userID int64, req *request.CreateGroup) error {
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

	if err := roomSvc.roomRepo.Create(ctx, &room); err != nil {
		return err
	}

	return nil
}
func (roomSvc *roomServiceImpl) AddMembersToGroup(ctx context.Context, userID int64, req *request.AddMemberToGroup) error {
	roomResult, err := roomSvc.roomRepo.FindByID(ctx, req.RoomID)
	if err != nil {
		return err
	}

	userIDMap := make(map[int64]bool)
	for _, userID := range req.UserIDs {
		userIDMap[userID] = true
	}

	isAdmin := false
	for _, member := range roomResult.Members {
		if userIDMap[member.UserID] {
			return errors.New("user sudah ada di group")
		}

		if member.Role == pkg.RoleAdmin && member.UserID == userID {
			isAdmin = true
		}
	}

	if !isAdmin {
		return errors.New("anda bukan admin")
	}

	userResults, err := roomSvc.userRepo.FindByIDs(ctx, req.UserIDs)
	if err != nil {
		return err
	}

	if len(userResults) != len(req.UserIDs) {
		return errors.New("user not found, gunakan user id yang valid")
	}

	// BELUM BUAT CEK USER APAKAH SUDAH ADA DI GROUP ATAU BELUM

	if err := roomSvc.memberRepo.AddMembers(ctx, req.RoomID, req.UserIDs); err != nil {
		return err
	}

	return nil
}

// func (roomSvc *roomServiceImpl) GetOrCreatePrivateRoom(ctx context.Context, userID, targetUserID int64) (*response.Room, error) {
// 	room := response.Room
// }
