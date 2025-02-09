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

	// 会话管理相关路由
	sessions := r.Group("/sessions")
	{
		sessions.GET("/activeSessions", controllers.GetActiveSessions)           // 获取活跃会话列表
		sessions.POST("/logoutSession", controllers.Logout)                      // 退出当前设备
		sessions.POST("/logoutAllSessions", controllers.LogoutAllDevices)        // 退出所有设备
		sessions.DELETE("/logoutSession/:session_id", controllers.RevokeSession) // 吊销指定会话
	}

	r.POST("/admin/addAdmin", controllers.AddAdmin)
	r.DELETE("/admin/removeAdmin", controllers.RemoveAdmin)
	r.GET("/admin/listAdmin", controllers.GetAdminList)

	r.GET("/admin/users", controllers.GetUserList)
	r.PUT("/admin/users/:id", controllers.UpdateUser)
	r.POST("/admin/users/:id/ban", controllers.BanUser)
	r.POST("/admin/users/:id/unban", controllers.UnbanUser)
	r.POST("/admin/users/generateUser", controllers.GenerateUsers)

	r.GET("/user/:id/activity", controllers.GetUserActivity) // 获取用户活跃度

	// 站内信相关路由
	messages := r.Group("/messages")
	{
		messages.GET("/getMessageList", controllers.GetMessageList)        // 获取消息列表
		messages.PUT("/:id/readMessage", controllers.MarkMessageAsRead)    // 标记消息已读
		messages.PUT("/readAll", controllers.MarkAllMessagesAsRead)        // 标记所有消息已读
		messages.POST("/batchRead", controllers.BatchMarkMessagesAsRead)   // 批量标记消息已读
		messages.DELETE("/:id/deleteMessage", controllers.DeleteMessage)   // 删除消息
		messages.GET("/getUnreadCount", controllers.GetUnreadMessageCount) // 获取未读消息数量
	}

	// 题目管理相关路由
	problems := r.Group("/problems")
	{
		problems.POST("", controllers.CreateProblem)                                   // 创建题目                          // 更新题目
		problems.DELETE("/:id", controllers.DeleteProblem)                             // 删除题目
		problems.GET("/:id", controllers.GetProblemDetail)                             // 获取题目详情
		problems.GET("", controllers.GetProblemList)                                   // 获取题目列表
		problems.POST("/switch-difficulty-system", controllers.SwitchDifficultySystem) // 切换难度等级系统
		problems.GET("/difficulty-system", controllers.GetDifficultySystem)            // 获取难度等级系统
	}

	// 管理员专用的题目管理路由
	adminProblems := r.Group("/admin/problems")
	{
		adminProblems.GET("", controllers.AdminGetProblemList)       // 管理员获取题目列表
		adminProblems.GET("/:id", controllers.AdminGetProblemDetail) // 管理员获取题目详情
		adminProblems.PUT("/:id", controllers.AdminUpdateProblem)    // 管理员更新题目
	}

	// 标签管理相关路由
	tags := r.Group("/tags")
	{
		tags.POST("", controllers.CreateTag)            // 创建标签
		tags.PUT("/:id", controllers.UpdateTag)         // 更新标签
		tags.DELETE("/:id", controllers.DeleteTag)      // 删除标签
		tags.GET("/getTagList", controllers.GetTagList) // 获取标签列表

		// 标签分类相关路由
		categories := tags.Group("/categories")
		{
			categories.POST("/createTagCategory", controllers.CreateTagCategory)       // 创建标签分类
			categories.PUT("/:id/updateTagCategory", controllers.UpdateTagCategory)    // 更新标签分类
			categories.DELETE("/:id/deleteTagCategory", controllers.DeleteTagCategory) // 删除标签分类
			categories.GET("/getTagCategoryList", controllers.GetTagCategoryList)      // 获取标签分类列表
			categories.GET("/getTagCategoryTree", controllers.GetTagCategoryTree)      // 获取标签分类树形结构
		}
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

	// 团队相关路由
	teams := r.Group("/teams")
	{
		teams.POST("/createTeam", controllers.CreateTeam)                     // 创建团队
		teams.PUT("/:id/updateTeam", controllers.UpdateTeam)                  // 更新团队信息
		teams.DELETE("/:id/deleteTeam", controllers.DeleteTeam)               // 删除团队
		teams.GET("/:id/getTeamDetail", controllers.GetTeamDetail)            // 获取团队详情
		teams.GET("/getTeamList", controllers.GetTeamList)                    // 获取团队列表
		teams.POST("/:id/createInvitation", controllers.CreateTeamInvitation) // 创建团队邀请
		teams.POST("/join", controllers.JoinTeam)                             // 加入团队
		teams.PUT("/:id/members/role", controllers.UpdateTeamMemberRole)      // 更新成员角色
		teams.DELETE("/:id/members/:user_id", controllers.RemoveTeamMember)   // 移除成员
		teams.GET("/:id/getMembers", controllers.GetTeamMemberList)           // 获取成员列表
		teams.PUT("/:id/changeNickname", controllers.UpdateTeamNickname)      // 更新团队内名称

		// 团队头像相关路由
		teams.POST("/:id/avatar", controllers.UploadTeamAvatar)   // 上传团队头像
		teams.GET("/avatar/:filename", controllers.GetTeamAvatar) // 获取团队头像
		teams.DELETE("/:id/avatar", controllers.RemoveTeamAvatar) // 删除团队头像

		// 团队申请相关路由
		teams.POST("/sendApply", controllers.CreateTeamApplication)       // 申请加入团队
		teams.GET("/getApplications", controllers.GetTeamApplicationList) // 获取申请列表
		teams.POST("/handleApply", controllers.HandleTeamApplication)     // 处理申请

		// 团队作业相关路由
		assignments := teams.Group("/assignments")
		{
			assignments.POST("/createAssignment", controllers.CreateAssignment)           // 创建作业
			assignments.PUT("/:id/updateAssignment", controllers.UpdateAssignment)        // 更新作业
			assignments.GET("/:id/getAssignmentDetail", controllers.GetAssignmentDetail)  // 获取作业详情
			assignments.GET("/getAssignmentList", controllers.GetAssignmentList)          // 获取作业列表
			assignments.GET("/getAvailableProblems", controllers.GetAvailableProblemList) // 获取可用题目列表
			assignments.GET("/getAssignmentProblems", controllers.GetAssignmentProblems)  // 获取作业题目列表
			assignments.GET("/getProblemDetail", controllers.GetAssignmentProblemDetail)  // 获取作业题目详情
			assignments.POST("/submitCode", controllers.SubmitAssignmentCode)             // 提交作业代码
			assignments.GET("/getSubmissions", controllers.GetAssignmentSubmissions)      // 获取作业提交记录
		}

		// 团队私有题目相关路由
		teamProblems := teams.Group("/problems")
		{
			teamProblems.POST("", controllers.CreateTeamProblem)       // 创建团队私有题目
			teamProblems.PUT("/:id", controllers.UpdateTeamProblem)    // 更新团队私有题目
			teamProblems.DELETE("/:id", controllers.DeleteTeamProblem) // 删除团队私有题目
			teamProblems.GET("/:id", controllers.GetTeamProblemDetail) // 获取团队私有题目详情
			teamProblems.GET("", controllers.GetTeamProblemList)       // 获取团队私有题目列表
		}

		// 团队题单相关路由
		problemLists := teams.Group("/problem-lists")
		{
			problemLists.POST("", controllers.CreateProblemList)       // 创建题单
			problemLists.PUT("/:id", controllers.UpdateProblemList)    // 更新题单
			problemLists.GET("/:id", controllers.GetProblemListDetail) // 获取题单详情
			problemLists.GET("", controllers.GetProblemListList)       // 获取题单列表
		}
	}
}
