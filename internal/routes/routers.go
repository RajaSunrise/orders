package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rajasunsire/orders/internal/handlers"
)

func Routes(app *fiber.App) {
	// Health endpoints
	app.Get("/health", handlers.HealthHandler)

	// CRUD endpoints for users
	app.Get("/users", handlers.GetUsers)
	app.Get("/users/:id", handlers.GetUser)
	app.Post("/users", handlers.CreateUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)

	// CRUD endpoints for products
	app.Get("/products", handlers.GetProducts)
	app.Get("/products/:id", handlers.GetProduct)
	app.Post("/products", handlers.CreateProduct)
	app.Put("/products/:id", handlers.UpdateProduct)
	app.Delete("/products/:id", handlers.DeleteProduct)

	// CRUD endpoints for orders
	app.Get("/orders", handlers.GetOrders)
	app.Get("/orders/:id", handlers.GetOrder)
	app.Post("/orders", handlers.CreateOrder)
	app.Put("/orders/:id", handlers.UpdateOrder)
	app.Delete("/orders/:id", handlers.DeleteOrder)
}
