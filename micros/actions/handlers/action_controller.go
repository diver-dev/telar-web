package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
	serviceConfig "github.com/red-gold/telar-web/micros/actions/config"
)

// DispatchHandle handles the creation of a new actionRoom
// @Summary Create a new actionRoom
// @Description Handles the creation of a new actionRoom by dispatching a request to the websocket server
// @Tags actions
// @Produce json
// @Security JWT
// @Param Authorization header string true "Authentication" default(Bearer <Add_token_here>)
// @Param roomId path string true "ActionRoom ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} utils.TelarError "Bad request"
// @Failure 500 {object} utils.TelarError "Internal server error"
// @Router /dispatch/{roomId} [post]
func DispatchHandle(c *fiber.Ctx) error {

	actionConfig := serviceConfig.ActionConfig

	// params from /actions/dispatch/:roomId
	actionRoomId := c.Params("roomId")
	if actionRoomId == "" {
		errorMessage := fmt.Sprintf("ActionRoom Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("actionRoomIdRequired", "ActionRoom Id is required!"))
	}

	bodyReader := bytes.NewBuffer(c.Body())
	URL := fmt.Sprintf("%s/api/dispatch/%s", actionConfig.WebsocketServerURL, actionRoomId)
	log.Info(" Dispatch URL: %s", URL)

	httpReq, httpErr := http.NewRequest(http.MethodPost, URL, bodyReader)
	if httpErr != nil {
		errorMessage := fmt.Sprintf("Error while creating dispatch request! %s", httpErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("createDispatchRequest", "Error while creating dispatch request!"))
	}

	xCloudSignature := c.Get(types.HeaderHMACAuthenticate)
	httpReq.Header.Add(types.HeaderHMACAuthenticate, xCloudSignature)
	httpReq.Header.Add("ORIGIN", *config.AppConfig.Gateway)

	httpClient := http.Client{}
	res, reqErr := httpClient.Do(httpReq)
	if reqErr != nil {
		errorMessage := fmt.Sprintf("Error while sending dispatch request to websocket server! %s", httpErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("createDispatchRequest", "Error while creating dispatch request!"))
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	return c.SendStatus(http.StatusOK)

}
