package routes

import (
	"OptiOJ/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.RegisterUser)
	r.POST("/verification/sendVerificationCode", controllers.RequestVerification)
	r.POST("/verification/validateCaptcha", controllers.ValidateGeetest)
}
