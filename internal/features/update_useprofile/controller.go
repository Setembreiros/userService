package update_userprofile

import (
	"errors"
	"userservice/internal/api"
	"userservice/internal/bus"
	database "userservice/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type PutUserProfileController struct {
	service *UpdateUserProfileService
}

func NewPutUserProfileController(repository Repository, bus *bus.EventBus) *PutUserProfileController {
	return &PutUserProfileController{
		service: NewUpdateUserProfileService(repository, bus),
	}
}

func (controller *PutUserProfileController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.PUT("/userprofile", controller.PutUserProfile)
}

func (controller *PutUserProfileController) PutUserProfile(c *gin.Context) {
	log.Info().Msg("Handling Request PUT UserProfile")
	var userProfile UserProfile

	if err := c.BindJSON(&userProfile); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		return
	}

	err := controller.service.UpdateUserProfile(&userProfile)
	if err != nil {
		var notFoundError *database.NotFoundError
		if errors.As(err, &notFoundError) {
			message := "User Profile not found for username " + userProfile.Username
			api.SendNotFound(c, message)
		} else {
			api.SendInternalServerError(c, err.Error())
		}
		return
	}

	api.SendOKWithResult(c, &userProfile)
}
