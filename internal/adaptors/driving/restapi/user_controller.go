package restapi

import (
	"apitest/internal/core/app/dto"
	"apitest/internal/core/app/ports"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserController struct {
	*BaseController
	userUseCase ports.UserUseCase
}

func NewUserCtrl(uc ports.UserUseCase) *UserController {
	return &UserController{
		userUseCase:    uc,
		BaseController: BaseHandlerObj,
	}
}

func (a *UserController) CreateUser(c *gin.Context) {
	var user dto.UserDTO

	err := c.BindJSON(&user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		log.Error().Err(err)
		return
	}

	log.Debug().Any("user", user).Send()
	err = a.userUseCase.CreateUser(user.ToAppUser())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		log.Error().Err(err).Msg("userUseCase.CreateUser failed")
		return
	}

	c.Status(http.StatusOK)
}
