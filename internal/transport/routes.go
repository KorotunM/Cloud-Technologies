package transport

import (
	"net/http"

	"pragma/internal/cache"
	"pragma/internal/config"
	"pragma/internal/storage"
	"pragma/internal/transport/handlers"
)

func SetupRoutes(cfg config.Config, redisClient *cache.RedisClient, storageClient *storage.Client) {
	handlers.InitSessionStore(cfg.Session.Key)

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/healthz", handlers.HealthHandler)
	http.HandleFunc("/universities", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUniversitiesHandler(w, r, redisClient, storageClient)
	})
	http.HandleFunc("/register", handlers.RegistrationHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/compare", handlers.CompareHandler)
	http.HandleFunc("/add/favorites", handlers.AddToFavoritesHandler)
	http.HandleFunc("/remove/favorites", handlers.RemoveFromFavoritesHandler)
	http.HandleFunc("/favorites", handlers.FavoritesHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		handlers.HistoryHandler(w, r, redisClient)
	})
	http.HandleFunc("/user", handlers.GetUserHandler)

	http.HandleFunc("/auth/google/redirect", handlers.GoogleRedirectHandler)
	http.HandleFunc("/auth/google/callback", handlers.GoogleCallbackHandler)
	http.HandleFunc("/auth/telegram/callback", handlers.TelegramCallbackHandler)
	http.HandleFunc("/auth/vk/redirect", handlers.VKRedirectHandler)
	http.HandleFunc("/auth/vk/callback", handlers.VKCallbackHandler)

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/css/", fs)
	http.Handle("/assets/", fs)
	http.Handle("/js/", fs)
}
