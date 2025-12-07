package consts

type contextKeysStruct struct {
	Uuid  string
	Email string
}

var ContextKeys = contextKeysStruct{
	Uuid:  "uuid",
	Email: "email",
}
