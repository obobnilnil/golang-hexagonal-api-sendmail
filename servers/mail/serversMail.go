package mail

import (
	"sendMail_git/mail/handlers"
	"sendMail_git/mail/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutesMail(router *gin.Engine) {

	s := services.NewServiceAdapter()
	h := handlers.NewHanerhandlerAdapter(s)

	router.POST("/api/mailChicCRM", h.MailChicCRMHandlers)
}
