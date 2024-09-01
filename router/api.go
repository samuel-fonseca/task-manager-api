package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samuel-fonseca/task-manager-api/controllers"
	"github.com/samuel-fonseca/task-manager-api/middleware"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	registerAuthenticationRoutes(app)
	registerTaskRoutes(app)

	app.Get("/api/user", middleware.DeserializeUser, controllers.GetUserDetails)
}

func registerTaskRoutes(app *fiber.App) {
	taskGroup := app.Group("api").Group("tasks", middleware.DeserializeUser)

	taskGroup.Get("/", controllers.ListTasks)
	taskGroup.Post("/", controllers.CreateTask)
	taskGroup.Get("/:id", controllers.ShowTask)
	taskGroup.Put("/:id", controllers.UpdateTask)
	taskGroup.Delete("/:id", controllers.DeleteTask)
}

func registerAuthenticationRoutes(app *fiber.App) {
	api := app.Group("api")

	api.Post("/login", controllers.LoginUser)
	api.Post("/register", controllers.RegisterUser)
}
