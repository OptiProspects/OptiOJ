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
	r.POST("/admin/users/generate", controllers.GenerateUsers)

	// 题目管理相关路由
	problems := r.Group("/problems")
	{
		problems.POST("", controllers.CreateProblem)                                   // 创建题目
		problems.PUT("/:id", controllers.UpdateProblem)                                // 更新题目
		problems.DELETE("/:id", controllers.DeleteProblem)                             // 删除题目
		problems.GET("/:id", controllers.GetProblemDetail)                             // 获取题目详情
		problems.GET("", controllers.GetProblemList)                                   // 获取题目列表
		problems.POST("/switch-difficulty-system", controllers.SwitchDifficultySystem) // 切换难度等级系统
		problems.GET("/difficulty-system", controllers.GetDifficultySystem)            // 获取难度等级系统
	}

	// 标签管理相关路由
	tags := r.Group("/tags")
	{
		tags.POST("", controllers.CreateTag)       // 创建标签
		tags.PUT("/:id", controllers.UpdateTag)    // 更新标签
		tags.DELETE("/:id", controllers.DeleteTag) // 删除标签
		tags.GET("", controllers.GetTagList)       // 获取标签列表
	}

	// 测试用例管理相关路由
	testcases := r.Group("/testcases")
	{
		testcases.POST("", controllers.UploadTestCase)                  // 上传测试用例
		testcases.DELETE("/:id", controllers.DeleteTestCase)            // 删除测试用例
		testcases.GET("/problem/:problem_id", controllers.GetTestCases) // 获取题目的测试用例列表
		testcases.GET("/:id/content", controllers.GetTestCaseContent)   // 获取测试用例内容
	}

	// 判题相关路由
	submissions := r.Group("/submissions")
	{
		submissions.POST("", controllers.SubmitCode)             // 提交代码
		submissions.GET("", controllers.GetSubmissionList)       // 获取提交记录列表
		submissions.GET("/:id", controllers.GetSubmissionDetail) // 获取提交记录详情
		submissions.POST("/debug", controllers.Debug)            // 在线调试代码
	}

	r.GET("/user/:id/activity", controllers.GetUserActivity) // 获取用户活跃度
}
