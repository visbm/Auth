package app

import (
	"auth/internal/composite"

	"github.com/julienschmidt/httprouter"
)

func ConfigureRouter(router *httprouter.Router, authcomposite composite.AuthComposite){
	authcomposite.Handler.Register(router)

}