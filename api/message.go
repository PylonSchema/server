package api

type MessageDatabase interface {
}

type MessageAPI struct {
	DB MessageDatabase
}
