package dto

import "github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"

type RoomDtoInterface interface {
	GetRoomInfo(room model.Room, userId string) RoomListResponse
	ResponseRoomList(rooms []model.Room, uuid string) []RoomListResponse
}

type RoomDtoStruct struct{}

func NewRoomDtoStruct() *RoomDtoStruct {
	return &RoomDtoStruct{}
}

type RoomListResponse struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	OwnerID     string `json:"OwnerID"`
	IsPrivate   bool   `json:"IsPrivate"`
	IsMember    bool   `json:"IsMember"`
	IsOwner     bool   `json:"IsOwner"`
	MemberCount int    `json:"MemberCount"`
	CreatedAt   string `json:"CreatedAt"`
}

func (s *RoomDtoStruct) contains(members []string, target string) bool {
	for _, v := range members {
		if v == target {
			return true
		}
	}
	return false
}

func (d *RoomDtoStruct) GetRoomInfo(room model.Room, userId string) RoomListResponse {
	return RoomListResponse{
		ID:          room.ID.Hex(),
		Name:        room.Name,
		OwnerID:     room.OwnerID,
		IsPrivate:   room.IsPrivate,
		IsMember:    d.contains(room.Members, userId),
		IsOwner:     room.OwnerID == userId,
		MemberCount: len(room.Members),
		CreatedAt:   room.CreatedAt.String(),
	}
}

func (d *RoomDtoStruct) ResponseRoomList(rooms []model.Room, uuid string) []RoomListResponse {
	var responses []RoomListResponse
	for _, room := range rooms {
		response := d.GetRoomInfo(room, uuid)
		responses = append(responses, response)
	}
	return responses
}
