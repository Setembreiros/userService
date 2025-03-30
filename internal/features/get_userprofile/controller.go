package get_userprofile

import (
	"errors"
	"userservice/internal/api"
	"userservice/internal/bus"
	database "userservice/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GetUserProfileController struct {
	service *GetUserProfileService
}

func NewGetUserProfileController(repository Repository, bus *bus.EventBus) *GetUserProfileController {
	return &GetUserProfileController{
		service: NewGetUserProfileService(repository, bus),
	}
}

func (controller *GetUserProfileController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/userprofile/:username", controller.GetUserProfile)
	routerGroup.GET("/userprofile/image/:username", controller.GetUserProfileImage)
}

func (controller *GetUserProfileController) GetUserProfile(c *gin.Context) {
	log.Info().Msg("Handling Request GET UserProfile")
	username := c.Param("username")

	userProfile, err := controller.service.GetUserProfile(username)
	if err != nil {
		var notFoundError *database.NotFoundError
		if errors.As(err, &notFoundError) {
			message := "User Profile not found for username " + username
			api.SendNotFound(c, message)
		} else {
			api.SendInternalServerError(c, err.Error())
		}
		return
	}

	api.SendOKWithResult(c, userProfile)
}

func (controller *GetUserProfileController) GetUserProfileImage(c *gin.Context) {
	log.Info().Msg("Handling Request GET UserProfile")
	username := c.Param("username")

	userProfileImage, err := controller.service.GetUserProfileImage(username)
	if err != nil {
		var notFoundError *database.NotFoundError
		if errors.As(err, &notFoundError) {
			message := "User Profile not found for username " + username
			api.SendNotFound(c, message)
		} else {
			api.SendInternalServerError(c, err.Error())
		}
		return
	}

	api.SendOKWithResult(c, userProfileImage)
}
