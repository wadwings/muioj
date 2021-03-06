package question

import (
	"MuiOJ-backEnd/controllers/auth"
	"MuiOJ-backEnd/models/forms"
	"MuiOJ-backEnd/services/question"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Edit(c *gin.Context){
	authObject, err := auth.GetAuthObjWithAdmin(c)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}
	questionEditForm := forms.QuestionEditForm{}
	if err := c.BindJSON(&questionEditForm); err != nil {
		c.JSON(400, gin.H{
			"code": 403,
			"message": err.Error(),
		})
		return
	}
	tid, err := strconv.Atoi(c.Param("tid"))
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}
	err = question.Edit(
		authObject.Uid,
		uint32(tid),
		&questionEditForm,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
	})
}