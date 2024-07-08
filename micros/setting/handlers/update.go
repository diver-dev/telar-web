package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	"github.com/red-gold/telar-web/micros/setting/database"
	domain "github.com/red-gold/telar-web/micros/setting/dto"
	models "github.com/red-gold/telar-web/micros/setting/models"
	service "github.com/red-gold/telar-web/micros/setting/services"
)

// @Summary Update user settings
// @Description Update user settings for the current user
// @Tags user-settings
// @Accept  json
// @Produce  json
// @Security JWT
// @Param Authorization header string true "Authentication" default(Bearer <Add_token_here>)
// @Param   body  body     models.UpdateSettingGroupModel  true "Update setting group model"
// @Success 200
// @Failure 400 {object} utils.TelarError
// @Failure 500 {object} utils.TelarError
// @Router / [put]
func UpdateUserSettingHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.UpdateSettingGroupModel)
	if err := c.BodyParser(model); err != nil {
		log.Error("[UpdateUserSettingHandle] %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	userSettingService, serviceErr := service.NewUserSettingService(database.Db)
	if serviceErr != nil {
		errorMessage := fmt.Sprintf("userSetting Error %s", serviceErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userSettingService", "Error happened while creating userSettingService!"))
	}

	if model.Type == "" {
		errorMessage := fmt.Sprintf("Setting type can not be empty Error")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("settingTypeEmptyError", "Setting type can not be empty Error"))
	}

	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		log.Error("[UpdateUserSettingHandle] Can not get current user")
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	var userSettings []domain.UserSetting
	for _, setting := range model.List {

		// TODO: Remove temporary function
		if setting.ObjectId == uuid.Nil && setting.Name == "send_email_app_news" {
			go func() {
				newUserSetting := &domain.UserSetting{
					OwnerUserId: currentUser.UserID,
					Name:        setting.Name,
					Value:       setting.Value,
					Type:        model.Type,
					IsSystem:    false,
				}

				if err := userSettingService.SaveUserSetting(newUserSetting); err != nil {
					errorMessage := fmt.Sprintf("Save UserSetting Error %s", err.Error())
					log.Error(errorMessage)
				}
			}()
		}
		updatedUserSetting := domain.UserSetting{
			ObjectId:    setting.ObjectId,
			OwnerUserId: currentUser.UserID,
			CreatedDate: utils.UTCNowUnix(),
			Name:        setting.Name,
			Value:       setting.Value,
			Type:        model.Type,
			IsSystem:    false,
		}
		userSettings = append(userSettings, updatedUserSetting)
	}

	if !(len(userSettings) > 0) {
		errorMessage := fmt.Sprintf("No setting added for update Error")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("noSettingForUpdate", "Can not find setting for update!"))
	}

	if err := userSettingService.UpdateUserSettingsById(currentUser.UserID, userSettings); err != nil {
		errorMessage := fmt.Sprintf("Update UserSetting Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("updateUserSetting", "Can not update user setting!"))
	}

	return c.SendStatus(http.StatusOK)

}
