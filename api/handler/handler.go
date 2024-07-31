package handler

type MainHandler interface {
	NewAutheticationHandler() AutheticationHandler
}
