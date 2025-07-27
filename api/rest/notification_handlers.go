package rest

import (
	"hermes-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// GetNotifications retrieves all notifications
func GetNotifications(c *fiber.Ctx) error {
	options := response.SuccessResponse([]string{}, "Notifications retrieved successfully")
	return response.ApiResponse(c, options)
}

// CreateNotification creates a new notification
func CreateNotification(c *fiber.Ctx) error {
	options := response.CreatedResponse(nil, "Notification created successfully")
	return response.ApiResponse(c, options)
}

// GetNotificationByID retrieves a notification by ID
func GetNotificationByID(c *fiber.Ctx) error {
	id := c.Params("id")
	options := response.SuccessResponse(fiber.Map{"id": id}, "Notification retrieved successfully")
	return response.ApiResponse(c, options)
}

// UpdateNotification updates a notification
func UpdateNotification(c *fiber.Ctx) error {
	id := c.Params("id")
	options := response.SuccessResponse(fiber.Map{"id": id}, "Notification updated successfully")
	return response.ApiResponse(c, options)
}

// DeleteNotification deletes a notification
func DeleteNotification(c *fiber.Ctx) error {
	_ = c.Params("id") // ID parameter for future use
	options := response.NoContentResponse("Notification deleted successfully")
	return response.ApiResponse(c, options)
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	options := response.SuccessResponse(fiber.Map{"id": id, "read": true}, "Notification marked as read")
	return response.ApiResponse(c, options)
}
