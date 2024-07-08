package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/pkg/parser"
	utils "github.com/red-gold/telar-core/utils"
	"github.com/red-gold/telar-web/micros/profile/database"
	service "github.com/red-gold/telar-web/micros/profile/services"
)

type UserProfileQueryModel struct {
	Search     string      `query:"search"`
	Page       int64       `query:"page"`
	NotInclude []uuid.UUID `query:"nin"`
}

// @Summary Query user profiles
// @Description Query user profiles by search query
// @Tags profiles
// @Accept  json
// @Produce  json
// @Security JWT
// @Param Authorization header string true "Authentication" default(Bearer <Add_token_here>)
// @Param   query  body     UserProfileQueryModel  true "User profile query model"
// @Success 200 {array} dto.UserProfile
// @Failure 400 {object} utils.TelarError
// @Failure 500 {object} utils.TelarError
// @Router / [get]
func QueryUserProfileHandle(c *fiber.Ctx) error {

	// Create service
	userService, serviceErr := service.NewUserProfileService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserProfileService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userProfileService", "Error happened while creating userProfileService!"))
	}

	query := new(UserProfileQueryModel)

	if err := parser.QueryParser(c, query); err != nil {
		log.Error("[QueryUserProfileHandle] QueryParser %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("queryParser", "Error happened while parsing query!"))
	}

	userList, err := userService.QueryUserProfile(query.Search, "created_date", query.Page, query.NotInclude)
	if err != nil {
		log.Error("[QueryUserProfile] %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userProfileService", "Error happened while creating userProfileService!"))
	}

	return c.JSON(userList)
}
