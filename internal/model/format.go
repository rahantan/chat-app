package model

import "errors"

var (
	ErrContactNotFound      = errors.New("contact not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrRoomNotFound         = errors.New("room not found")
	ErrMemberNotFound       = errors.New("member not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrDuplicatePrivateRoom = errors.New("duplicate private room")
)
