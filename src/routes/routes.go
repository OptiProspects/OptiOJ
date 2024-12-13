package routes

import (
	"OptiOJ/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/auth/userRegister", controllers.RegisterUser)
	r.POST("/auth/userLogin", controllers.LoginUser)
	r.GET("/auth/refreshToken", controllers.RefreshToken)
	r.GET("/user/globalData", controllers.GetGlobalData)
	r.POST("/verification/sendVerificationCode", controllers.RequestVerification)
	r.POST("/verification/validateCaptcha", controllers.ValidateGeetest)
}
