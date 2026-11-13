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

type MessageService interface {
	ListMessage(ctx context.Context, userID, roomID int64) ([]response.Message, error) // done
	SendMessage(ctx context.Context, userID int64, req *request.AddMessage) error      // done
	DeleteMessage(ctx context.Context, userID, roomID, msgID int64) error              // done
	EditMessage(ctx context.Context, userID, roomID, msgID int64) error
}

type messageServiceImpl struct {
	roomRepo   mysql.RoomRepository
	memberRepo mysql.RoomMemberRepository
	msgRepo    mysql.MessageRepository
}

func (msgSvc *messageServiceImpl) roomValidate(ctx context.Context, userID, roomID int64) error {
	exist, err := msgSvc.roomRepo.IsExist(ctx, roomID)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("not found room")
	}

	_, err = msgSvc.memberRepo.FindByUserIDAndRoomID(ctx, userID, roomID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return errors.New("tidak ada hak akses")
		}
		return err
	}

	return nil
}

// cek user bisa akses atau ga
// ambil data
func (msgSvc *messageServiceImpl) ListMessage(ctx context.Context, userID, roomID int64) ([]response.Message, error) {

	if err := msgSvc.roomValidate(ctx, userID, roomID); err != nil {
		return nil, err
	}

	msgResults, err := msgSvc.msgRepo.FindAllByRoomID(ctx, roomID)
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

func (msgSvc *messageServiceImpl) SendMessage(ctx context.Context, userID int64, req *request.AddMessage) error {

	if err := msgSvc.roomValidate(ctx, userID, req.RoomID); err != nil {
		return err
	}

	msg := model.MessageModel{
		SenderID: userID,
		RoomID:   req.RoomID,
		MsgType:  req.MsgType,
		Content:  req.Content,
	}

	if err := msgSvc.msgRepo.Add(ctx, &msg); err != nil {
		return err
	}

	return nil
}

func (msgSvc *messageServiceImpl) DeleteMessage(ctx context.Context, userID, roomID, msgID int64) error {
	if err := msgSvc.msgRepo.Delete(ctx, userID, roomID, msgID); err != nil {
		return err
	}
	return nil
}
