package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/samuel-fonseca/task-manager-api/database"
	"github.com/samuel-fonseca/task-manager-api/model"
)

func ListTasks(c *fiber.Ctx) error {
	var tasks []model.Task
	user := c.Locals("user").(model.UserDetailsData)
	database.DB.Find(&tasks, "user_id = ?", user.ID)

	return c.JSON(fiber.Map{
		"tasks": tasks,
	})
}

// TODO: Authentication middleware to be added
func ShowTask(c *fiber.Ctx) error {
	var task model.Task
	id := c.Params("id")
	user := c.Locals("user").(model.UserDetailsData)

	database.DB.First(&task, "id = ? and user_id = ?", id, user.ID)

	if task.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Task not found",
		})
	}

	return c.JSON(fiber.Map{
		"task": task,
	})
}

func CreateTask(c *fiber.Ctx) error {
	var payload = new(model.TaskInput)
	err := c.BodyParser(&payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	user := c.Locals("user").(model.UserDetailsData)

	task := model.Task{
		Title:       payload.Title,
		Description: payload.Description,
		Status:      payload.Status,
		UserID:      user.ID,
	}

	result := database.DB.Create(&task)

	if result.Error != nil {
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": fmt.Sprintf("Could not create task: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Task created.",
		"task":    task,
	})
}

func UpdateTask(c *fiber.Ctx) error {
	reqId := c.Params("id")
	var task model.Task
	user := c.Locals("user").(model.UserDetailsData)

	database.DB.First(&task, "id = ? and user_id = ?", reqId, user.ID)

	if task.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Task not found.",
		})
	}

	var payload = new(model.TaskInput)
	err := c.BodyParser(&payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	task.Title = payload.Title
	task.Description = payload.Description
	task.Status = payload.Status

	database.DB.Save(&task)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Updated task.",
		"task":    task,
	})
}

func DeleteTask(c *fiber.Ctx) error {
	reqId := c.Params("id")
	var task model.Task
	user := c.Locals("user").(model.UserDetailsData)

	result := database.DB.Where("id = ? and user_id = ?", reqId, user.ID).Delete(&task)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted.",
	})
}
