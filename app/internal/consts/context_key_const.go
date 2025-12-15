package consts

type contextKeysStruct struct {
	Uuid      string
	Email     string
	RoomModel string
	IsAdmin   string
	IsMember  string
}

var ContextKeys = contextKeysStruct{
	Uuid:      "uuid",
	Email:     "email",
	RoomModel: "room_model",
	IsAdmin:   "is_admin",
	IsMember:  "is_member",
}
