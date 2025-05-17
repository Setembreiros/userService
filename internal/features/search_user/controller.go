package search_user

import (
	"strconv"
	"userservice/internal/api"
	"userservice/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type SearchUserController struct {
	service ControllerService
}

type ControllerService interface {
	SearchUserProfileSnippets(query, lastUsername string, limit int) ([]*model.UserProfileSnippet, string, error)
}

type GetUserProfileSnippetsResponse struct {
	Users        []*model.UserProfileSnippet `json:"users"`
	LastUsername string                      `json:"lastUsername"`
}

func NewSearchUserController(service ControllerService) *SearchUserController {
	return &SearchUserController{
		service: service,
	}
}

func (controller *SearchUserController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/userprofile-snippets", controller.SearchUser)
}

func (controller *SearchUserController) SearchUser(c *gin.Context) {
	log.Info().Msg("Handling Request GET UserProfile Snippets")

	query, lastUsername, limit, err := getQueryParameters(c)
	if err != nil || limit <= 0 {
		return
	}

	userProfileSnippets, lastUsername, err := controller.service.SearchUserProfileSnippets(query, lastUsername, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetUserProfileSnippetsResponse{
		Users:        userProfileSnippets,
		LastUsername: lastUsername,
	})
}

func getQueryParameters(c *gin.Context) (string, string, int, error) {
	query := c.DefaultQuery("query", "")

	lastUsername := c.DefaultQuery("lastUsername", "")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Invalid pagination parameters, limit must be greater than 0")
		return "", "", 0, err
	}

	return query, lastUsername, limit, nil
}
