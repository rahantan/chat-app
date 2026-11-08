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

type ChatService interface {
	List(ctx context.Context, userID, roomID int64)
	History(ctx context.Context, userID, roomID int64) ([]response.Message, error) // done
	SendMessage(ctx context.Context, userID int64, req *request.AddMessage) error  // done
	DeleteMessage(ctx context.Context, userID, roomID, msgID int64) error          // done
}

type chatServiceImpl struct {
	roomRepo   mysql.RoomRepository
	memberRepo mysql.RoomMemberRepository
	msgRepo    mysql.MessageRepository
}

// cek user bisa akses atau ga
// ambil data
func (cgSvc *chatServiceImpl) History(ctx context.Context, userID, roomID int64) ([]response.Message, error) {

	exist, err := cgSvc.roomRepo.IsExist(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("not found room")
	}

	_, err = cgSvc.memberRepo.FindByUserIDAndRoomID(ctx, userID, roomID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, errors.New("tidak ada hak akses")
		}
		return nil, err
	}

	msgResults, err := cgSvc.msgRepo.FindAllByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	messages := make([]response.Message, 0)
	for _, msg := range msgResults {
		sender := response.Sender{
			SenderID:  msg.SenderID,
			UserName:  msg.Sender.UserName,
			ImagePath: msg.Sender.Image,
		}

		message := response.Message{
			MsgID:       msg.ID,
			MsgType:     msg.MsgType,
			Content:     msg.Content,
			Sender:      &sender,
			CreatedAt:   msg.CreatedAt.Format(pkg.DateTimeLayout),
			UpdatedAtAt: msg.UpdatedAt.Format(pkg.DateTimeLayout),
		}
		if msg.MsgType != pkg.MsgTypeText {
			attachments := make([]response.Attachment, 0)
			for _, attchmnt := range msg.Attachments {
				attachments = append(attachments, response.Attachment{
					AttachmentID: attchmnt.ID,
					FileType:     attchmnt.FileType,
					FileSize:     attchmnt.FileSize,
					FilePath:     pkg.MsgFilePath + attchmnt.FileName,
				})
			}
			message.Attachments = append(message.Attachments, attachments...)
		}

		messages = append(messages, message)
	}
	return messages, nil

}
func (cgSvc *chatServiceImpl) SendMessage(ctx context.Context, userID int64, req *request.AddMessage) error {
	msg := model.MessageModel{
		SenderID: userID,
		RoomID:   req.RoomID,
		MsgType:  req.MsgType,
		Content:  req.Content,
	}

	exist, err := cgSvc.roomRepo.IsExist(ctx, req.RoomID)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("not found room")
	}

	_, err = cgSvc.memberRepo.FindByUserIDAndRoomID(ctx, userID, req.RoomID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return errors.New("tidak diizinkan mengirim pesan")
		}
		return err
	}

	if err := cgSvc.msgRepo.Add(ctx, &msg); err != nil {
		return err
	}

	return nil
}

func (cgSvc *chatServiceImpl) DeleteMessage(ctx context.Context, userID, roomID, msgID int64) error {
	if err := cgSvc.msgRepo.Delete(ctx, userID, roomID, msgID); err != nil {
		return err
	}
	return nil
}
