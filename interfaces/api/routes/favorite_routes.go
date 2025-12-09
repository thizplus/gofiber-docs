package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupFavoriteRoutes(api fiber.Router, h *handlers.Handlers) {
	favorites := api.Group("/favorites")
	favorites.Use(middleware.Protected())

	favorites.Post("/", h.FavoriteHandler.AddFavorite)
	favorites.Get("/", h.FavoriteHandler.GetFavorites)
	favorites.Delete("/:id", h.FavoriteHandler.RemoveFavorite)
	favorites.Get("/check", h.FavoriteHandler.CheckFavorite)
	favorites.Post("/toggle", h.FavoriteHandler.ToggleFavorite)
}
