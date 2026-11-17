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
	ListMessage(ctx context.Context, userID, roomID int64) ([]response.Message, error)    // done
	SendMessage(ctx context.Context, userID, roomID int64, req *request.AddMessage) error // done
	DeleteMessage(ctx context.Context, userID, roomID, msgID int64) error                 // done
	// EditMessage(ctx context.Context, userID, roomID, msgID int64) error
}

type messageServiceImpl struct {
	roomRepo   mysql.RoomRepository
	memberRepo mysql.RoomMemberRepository
	msgRepo    mysql.MessageRepository
	tx         mysql.TransactorRepository
}

func NewMessageService(room mysql.RoomRepository, member mysql.RoomMemberRepository, msg mysql.MessageRepository, tx mysql.TransactorRepository) MessageService {
	return &messageServiceImpl{
		roomRepo:   room,
		memberRepo: member,
		msgRepo:    msg,
		tx:         tx,
	}
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
			UserID:    msg.SenderID,
			UserName:  msg.Sender.UserName,
			AvatarUrl: pkg.AvatarImageUrl + msg.Sender.Avatar,
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
					FileUrl:      pkg.MsgFileUrl + attchmnt.FileName,
				})
			}
			message.Attachments = append(message.Attachments, attachments...)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (msgSvc *messageServiceImpl) SendMessage(ctx context.Context, userID, roomID int64, req *request.AddMessage) error {

	if err := msgSvc.roomValidate(ctx, userID, roomID); err != nil {
		return err
	}

	msg := model.MessageModel{
		SenderID: userID,
		RoomID:   roomID,
		MsgType:  req.MsgType,
		Content:  req.Content,
	}

	return msgSvc.tx.WithInTransaction(ctx, func(ctxTx context.Context) error {
		if err := msgSvc.msgRepo.Add(ctx, &msg); err != nil {
			return err
		}

		if err := msgSvc.roomRepo.UpdateTime(ctxTx, msg.RoomID); err != nil {
			return err
		}
		return nil
	})

}

func (msgSvc *messageServiceImpl) DeleteMessage(ctx context.Context, userID, roomID, msgID int64) error {
	msgResult, err := msgSvc.msgRepo.FindByID(ctx, msgID)
	if err != nil {
		return err
	}
	if userID != msgResult.SenderID {
		return pkg.NewErr(pkg.ErrCodeForbidden, "you are not allowed to delete this message")
	}

	if err := msgSvc.msgRepo.Delete(ctx, userID, roomID, msgID); err != nil {
		return err
	}
	return nil
}
