package service

import (
	"context"
	"realtime-chat/internal/dto/request"
	"realtime-chat/internal/dto/response"
	"realtime-chat/internal/repository/mysql"
)

type RoomService interface {
	ListPrivateRoom(ctx context.Context)
	// CreatePrivateRoom(ctx context.Context, userID int64, req *request.CreateRoom) error

	ListGroupRoom(ctx context.Context, userID int64)
	CreateGroup(ctx context.Context, userID int64, req *request.CreateGroup) error     // done
	AddMemberToGroup(ctx context.Context, userID, req *request.AddMemberToGroup) error // done
	ListMemberGroup(ctx context.Context, userID, roomID int64) ([]response.GroupMember, error)
	DeleteGroup(ctx context.Context, userID, roomID int64) error
}

type roomServiceImpl struct {
	roomRepo   mysql.RoomRepository
	memberRepo mysql.RoomMemberRepository
}

/*
1.list room
- ambil semua room
- ambil last message
- join members.User WHERE type roomnya adalah private
*/

// func (roomSvc *roomServiceImpl) ListGroupRoom(ctx context.Context, userID int64) ([]response.Room, error) {
// 	roomResults, err := roomSvc.roomRepo.List(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rooms := make([]response.Room, 0)
// 	for _, room := range roomResults {
// 		var (
// 			image string
// 			name  string
// 		)

// 		if room.Type == pkg.RoomTypePrivate {
// 			user := room.Members[0].User
// 			image = user.Image
// 			name = user.UserName
// 		} else {
// 			image = room.Image
// 			name = room.Title
// 		}

// 		rooms = append(rooms, response.Room{
// 			RoomID:    room.ID,
// 			Name:      name,
// 			ImagePath: pkg.ImagePath + image,
// 			UpdatedAt: room.Image,
// 		})
// 	}
// }
