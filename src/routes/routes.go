package routes

import (
	"OptiOJ/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/auth/userRegister", controllers.RegisterUser)
	r.POST("/auth/userLogin", controllers.LoginUser)
	r.POST("/user/uploadAvatar", controllers.UploadAvatar)
	r.PUT("/user/updateProfile", controllers.UpdateProfile)
	r.GET("/auth/refreshToken", controllers.RefreshToken)
	r.GET("/user/globalData", controllers.GetGlobalData)
	r.GET("/user/getAvatar", controllers.GetAvatar)
	r.GET("/user/getProvinces", controllers.GetProvinces)
	r.GET("/user/getCities", controllers.GetCities)
	r.DELETE("/user/removeAvatar", controllers.RemoveAvatar)
	r.POST("/verification/sendVerificationCode", controllers.RequestVerification)
	r.POST("/verification/validateCaptcha", controllers.ValidateGeetest)

	r.POST("/admin/addAdmin", controllers.AddAdmin)
	r.DELETE("/admin/removeAdmin", controllers.RemoveAdmin)
	r.GET("/admin/listAdmin", controllers.GetAdminList)

	r.GET("/admin/users", controllers.GetUserList)
	r.PUT("/admin/users/:id", controllers.UpdateUser)
	r.POST("/admin/users/:id/ban", controllers.BanUser)
	r.POST("/admin/users/:id/unban", controllers.UnbanUser)
}
