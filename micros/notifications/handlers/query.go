package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/pkg/parser"
	"github.com/red-gold/telar-core/types"
	utils "github.com/red-gold/telar-core/utils"
	"github.com/red-gold/telar-web/micros/notifications/database"
	models "github.com/red-gold/telar-web/micros/notifications/models"
	service "github.com/red-gold/telar-web/micros/notifications/services"
)

type UserProfileQueryModel struct {
	Limit int64 `query:"limit"`
	Page  int64 `query:"page"`
}

// GetNotificationsByUserIdHandle godoc
// @Summary Get notifications by user ID
// @Description Handles querying notifications by user ID
// @Tags Notification
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} domain.Notification
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /notifications [get]
func GetNotificationsByUserIdHandle(c *fiber.Ctx) error {

	// Create service
	notificationService, serviceErr := service.NewNotificationService(database.Db)
	if serviceErr != nil {
		log.Error("[GetNotificationsByUserIdHandle.NewNotificationService] %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/notificationService", "Error happened while creating notificationService!"))
	}

	query := new(UserProfileQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[GetNotificationsByUserIdHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		log.Error("[GetNotificationsByUserIdHandle] Can not get current user")
		return c.Status(http.StatusUnauthorized).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	notificationList, err := notificationService.GetNotificationByUserId(&currentUser.UserID, "created_date", query.Page, query.Limit)

	if err != nil {
		log.Error("[GetNotificationsByUserIdHandle.GetNotificationByUserId] %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getNotificationByUserId", "Error happened while reading notification!"))
	}

	return c.JSON(notificationList)

}

// GetNotificationHandle godoc
// @Summary Get a notification
// @Description Handles retrieving a notification by its ID
// @Tags Notification
// @Accept json
// @Produce json
// @Param notificationId path string true "Notification ID"
// @Success 200 {object} models.NotificationModel
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /notifications/id/{notificationId} [get]
func GetNotificationHandle(c *fiber.Ctx) error {

	// Create service
	notificationService, serviceErr := service.NewNotificationService(database.Db)
	if serviceErr != nil {
		log.Error("[GetNotificationHandle.NewNotificationService] %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/notificationService", "Error happened while creating notificationService!"))
	}
	notificationId := c.Params("notificationId")
	notificationUUID, uuidErr := uuid.FromString(notificationId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("Notification Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("notificationIdRequired", "Notification id is required!"))

	}

	foundNotification, err := notificationService.FindById(notificationUUID)
	if err != nil {
		log.Error("[GetNotificationHandle.notificationService.FindById] %s - %s", notificationUUID.String(), serviceErr.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("findNotification", "Error happened while finding notification!"))
	}

	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		log.Error("[GetNotificationHandle] Can not get current user")
		return c.Status(http.StatusUnauthorized).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	notificationModel := models.NotificationModel{
		ObjectId:             foundNotification.ObjectId,
		OwnerUserId:          currentUser.UserID,
		OwnerDisplayName:     currentUser.DisplayName,
		OwnerAvatar:          currentUser.Avatar,
		Title:                foundNotification.Title,
		Description:          foundNotification.Description,
		URL:                  foundNotification.URL,
		NotifyRecieverUserId: foundNotification.NotifyRecieverUserId,
		TargetId:             foundNotification.TargetId,
		IsSeen:               foundNotification.IsSeen,
		Type:                 foundNotification.Type,
		EmailNotification:    foundNotification.EmailNotification,
	}

	return c.JSON(notificationModel)

}
