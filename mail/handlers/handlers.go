package handlers

import (
	"net/http"
	"sendMail_git/mail/models"
	"sendMail_git/mail/services"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	MailChicCRMHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicePort
}

func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) MailChicCRMHandlers(c *gin.Context) {
	var mailRequest models.MailRequest
	if err := c.ShouldBind(&mailRequest); err != nil {
		c.JSON(400, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	form, _ := c.MultipartForm()
	files := form.File["attachment"]
	attachmentURLs, err := h.s.MailChicCRMServices(mailRequest, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": attachmentURLs})
}
