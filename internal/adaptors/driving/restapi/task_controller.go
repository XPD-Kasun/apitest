package restapi

import "github.com/gin-gonic/gin"

type TaskController struct {
}

func NewTaskCtrl() *TaskController {
	return &TaskController{}
}

func (t *TaskController) Fn(c *gin.Context) {

}
