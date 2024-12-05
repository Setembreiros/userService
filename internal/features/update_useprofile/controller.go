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

type UpdateUserProfileImageResponse struct {
	PresignedUrl string `json:"presigned_url"`
}

func NewPutUserProfileController(repository Repository, bus *bus.EventBus) *PutUserProfileController {
	return &PutUserProfileController{
		service: NewUpdateUserProfileService(repository, bus),
	}
}

func (controller *PutUserProfileController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.PUT("/userprofile", controller.PutUserProfile)
	routerGroup.PUT("/userprofile/image", controller.PutUserProfileImage)
	routerGroup.PUT("/userprofile/confirm-updated-image", controller.ConfirmUserProfileImageUpdated)
}

func (controller *PutUserProfileController) PutUserProfile(c *gin.Context) {
	log.Info().Msg("Handling Request PUT UserProfile")
	var userProfile UserProfile

	if err := c.BindJSON(&userProfile); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
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

func (controller *PutUserProfileController) PutUserProfileImage(c *gin.Context) {
	log.Info().Msg("Handling Request PUT UserProfileImage")
	var userProfileImage UserProfileImage

	if err := c.BindJSON(&userProfileImage); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
		return
	}

	presignedUrl, err := controller.service.UpdateUserProfileImage(&userProfileImage)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &UpdateUserProfileImageResponse{
		PresignedUrl: presignedUrl,
	})
}

func (controller *PutUserProfileController) ConfirmUserProfileImageUpdated(c *gin.Context) {
	log.Info().Msg("Handling Request PUT ConfirmUserProfileImageUpdated")

	var updatedImageConfirm ConfirmUserProfileImageUpdated

	if err := c.BindJSON(&updatedImageConfirm); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
		return
	}

	err := controller.service.ConfirmUserProfileImageUpdated(&updatedImageConfirm)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOK(c)
}
