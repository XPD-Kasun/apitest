package restapi

import (
	"apitest/internal/core/app/ports"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	*BaseController
	authUseCase ports.AuthUseCase
}

func NewAuthCtrl(ac ports.AuthUseCase) *AuthController {

	return &AuthController{
		authUseCase:    ac,
		BaseController: BaseHandlerObj,
	}
}

func (h *AuthController) Login(c *gin.Context) {

	c.Request.ParseForm()
	uname := c.Request.PostFormValue("username")
	password := c.Request.PostFormValue("password")

	token, err := h.authUseCase.LoginFromPassword(uname, password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  err.Error(),
			"status": "Unauthorized",
		})
		log.Err(err).Msg("LoginFromPassword failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *AuthController) Fn(c *gin.Context) {
	c.String(200, "Thsi is fn")
}

func (a *AuthController) OAuth2Login(c *gin.Context) {

}
