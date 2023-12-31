package main

import (
	"net/http"

	"github.com/SmileL1ne/web-mailing/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(lg handlers.Logic) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", lg.Home())
	mux.Post("/api/subscribe", lg.GetSubscriber())
	mux.Post("/api/send", lg.SendMail())
	mux.Post("/api/unsubscribe", lg.DeleteSubscriber())

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
