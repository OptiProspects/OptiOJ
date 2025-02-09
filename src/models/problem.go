package models

import "time"

// DifficultySystem 难度等级系统
type DifficultySystem string

const (
	DifficultySystemNormal DifficultySystem = "normal" // 普通难度等级系统
	DifficultySystemOI     DifficultySystem = "oi"     // OI 难度等级系统
)

// 普通难度等级
const (
	DifficultyNormalUnrated = "unrated" // 暂无评级
	DifficultyNormalEasy    = "easy"    // 简单
	DifficultyNormalMedium  = "medium"  // 中等
	DifficultyNormalHard    = "hard"    // 困难
)

// OI 难度等级
const (
	DifficultyOIUnrated   = "unrated"    // 暂无评级
	DifficultyOIBeginner  = "beginner"   // 入门/蒟蒻
	DifficultyOIBasic     = "basic"      // 普及-
	DifficultyOIBasicPlus = "basicplus"  // 普及/提高-
	DifficultyOIAdv       = "advanced"   // 普及+/提高
	DifficultyOIAdvPlus   = "advplus"    // 提高+/省选-
	DifficultyOIProv      = "provincial" // 省选/NOI-
	DifficultyOINOI       = "noi"        // NOI/NOI+/CTSC/神犇
)

// DifficultyConfig 难度等级配置
type DifficultyConfig struct {
	System       DifficultySystem
	Difficulties map[string]string // key: 难度代码, value: 难度显示名称
}

var (
	// DefaultDifficultyConfigs 默认难度等级配置
	DefaultDifficultyConfigs = map[DifficultySystem]DifficultyConfig{
		DifficultySystemNormal: {
			System: DifficultySystemNormal,
			Difficulties: map[string]string{
				DifficultyNormalUnrated: "暂无评级",
				DifficultyNormalEasy:    "简单",
				DifficultyNormalMedium:  "中等",
				DifficultyNormalHard:    "困难",
			},
		},
		DifficultySystemOI: {
			System: DifficultySystemOI,
			Difficulties: map[string]string{
				DifficultyOIUnrated:   "暂无评级",
				DifficultyOIBeginner:  "入门/蒟蒻",
				DifficultyOIBasic:     "普及-",
				DifficultyOIBasicPlus: "普及/提高-",
				DifficultyOIAdv:       "普及+/提高",
				DifficultyOIAdvPlus:   "提高+/省选-",
				DifficultyOIProv:      "省选/NOI-",
				DifficultyOINOI:       "NOI/NOI+/CTSC/神犇",
			},
		},
	}

	// CustomDifficultyConfigs 自定义难度等级配置
	CustomDifficultyConfigs = make(map[DifficultySystem]DifficultyConfig)
)

// GetDifficultyConfig 获取难度等级配置
func GetDifficultyConfig(system DifficultySystem) (DifficultyConfig, bool) {
	// 优先使用自定义配置
	if config, ok := CustomDifficultyConfigs[system]; ok {
		return config, true
	}
	// 否则使用默认配置
	config, ok := DefaultDifficultyConfigs[system]
	return config, ok
}

// IsValidDifficulty 验证难度等级是否有效
func IsValidDifficulty(system DifficultySystem, difficulty string) bool {
	config, ok := GetDifficultyConfig(system)
	if !ok {
		return false
	}
	_, exists := config.Difficulties[difficulty]
	return exists
}

// GetDifficultyDisplay 获取难度等级的显示名称
func GetDifficultyDisplay(system DifficultySystem, difficulty string) string {
	config, ok := GetDifficultyConfig(system)
	if !ok {
		return difficulty
	}
	if display, exists := config.Difficulties[difficulty]; exists {
		return display
	}
	return difficulty
}

// Problem 题目模型
type Problem struct {
	ID                uint64           `json:"id"`
	Title             string           `json:"title"`
	Description       string           `json:"description"`
	InputDescription  string           `json:"input_description"`
	OutputDescription string           `json:"output_description"`
	SampleCases       string           `json:"sample_cases" gorm:"column:samples"`
	Hint              string           `json:"hint"`
	Source            string           `json:"source"`
	DifficultySystem  DifficultySystem `json:"difficulty_system"` // 难度等级系统
	Difficulty        string           `json:"difficulty"`        // 难度等级
	TimeLimit         int              `json:"time_limit"`
	MemoryLimit       int              `json:"memory_limit"`
	IsPublic          bool             `json:"is_public"`
	CreatedBy         uint64           `json:"created_by"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// ProblemCategory 题目分类
type ProblemCategory struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint64   `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// TagCategory 标签分类
type TagCategory struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint64   `json:"parent_id"` // 父分类ID，为空表示一级分类
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TagCategoryDetail 标签分类详情（包含子分类）
type TagCategoryDetail struct {
	TagCategory
	Children []TagCategoryDetail `json:"children,omitempty"` // 子分类列表
}

// ProblemTag 题目标签
type ProblemTag struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	Color      string    `json:"color"`
	CategoryID *uint64   `json:"category_id"` // 所属分类ID
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TagWithCategory 带有分类信息的标签
type TagWithCategory struct {
	ProblemTag
	Category     *TagCategory `json:"category,omitempty"`      // 所属分类
	CategoryPath []string     `json:"category_path,omitempty"` // 分类路径
}

// TestCase 测试用例
type TestCase struct {
	ID         uint64    `json:"id"`
	ProblemID  uint64    `json:"problem_id"`
	InputFile  string    `json:"input_file"`  // 输入文件路径
	OutputFile string    `json:"output_file"` // 输出文件路径
	CreatedAt  time.Time `json:"created_at"`
}

// TestCaseWithLocalID 带有局部ID的测试用例信息
type TestCaseWithLocalID struct {
	TestCase
	LocalID int `json:"local_id"` // 测试用例在当前题目中的序号
}

// TestCaseContentResponse 测试用例内容响应
type TestCaseContentResponse struct {
	ID        uint64 `json:"id"`         // 测试用例ID
	LocalID   int    `json:"local_id"`   // 局部ID
	ProblemID uint64 `json:"problem_id"` // 题目ID
	Input     string `json:"input"`      // 输入内容
	Output    string `json:"output"`     // 输出内容
}

// ProblemDetail 题目详细信息（包含分类和标签）
type ProblemDetail struct {
	Problem
	Categories []ProblemCategory `json:"categories"`
	Tags       []ProblemTag      `json:"tags"`
	UserStatus *string           `json:"user_status"` // 用户状态：null-未提交, accepted-已通过, attempted-尝试过
}

// CreateProblemRequest 创建题目请求
type CreateProblemRequest struct {
	Title             string           `json:"title" binding:"required"`
	Description       string           `json:"description" binding:"required"`
	InputDescription  string           `json:"input_description"`
	OutputDescription string           `json:"output_description"`
	Samples           string           `json:"samples"`
	Hint              string           `json:"hint"`
	Source            string           `json:"source"`
	DifficultySystem  DifficultySystem `json:"difficulty_system" binding:"required,oneof=normal oi"`
	Difficulty        string           `json:"difficulty" binding:"required"`
	TimeLimit         int              `json:"time_limit" binding:"required,min=100,max=10000"`
	MemoryLimit       int              `json:"memory_limit" binding:"required,min=16,max=1024"`
	IsPublic          bool             `json:"is_public"`
	CategoryIDs       []uint64         `json:"category_ids"`
	TagIDs            []uint64         `json:"tag_ids"`
}

// UpdateProblemRequest 更新题目请求
type UpdateProblemRequest struct {
	Title             *string           `json:"title"`
	Description       *string           `json:"description"`
	InputDescription  *string           `json:"input_description"`
	OutputDescription *string           `json:"output_description"`
	Samples           *string           `json:"samples"`
	Hint              *string           `json:"hint"`
	Source            *string           `json:"source"`
	DifficultySystem  *DifficultySystem `json:"difficulty_system"`
	Difficulty        *string           `json:"difficulty"`
	TimeLimit         *int              `json:"time_limit" binding:"omitempty,min=100,max=10000"`
	MemoryLimit       *int              `json:"memory_limit" binding:"omitempty,min=16,max=1024"`
	IsPublic          *bool             `json:"is_public"`
	CategoryIDs       []uint64          `json:"category_ids"`
	TagIDs            []uint64          `json:"tag_ids"`
}

// ProblemListRequest 获取题目列表的请求参数
type ProblemListRequest struct {
	Page       int      `form:"page" binding:"required,min=1"`
	PageSize   int      `form:"page_size" binding:"required,min=1,max=100"`
	Title      string   `form:"title"`
	Difficulty string   `form:"difficulty"`
	Tags       []uint64 `form:"tags"`
	Categories []uint64 `form:"categories"`
	IsPublic   *bool    `form:"is_public"`
}

// ProblemListItem 题目列表项
type ProblemListItem struct {
	ID              uint64            `json:"id"`
	Title           string            `json:"title"`
	Difficulty      string            `json:"difficulty"`
	Tags            []ProblemTag      `json:"tags"`
	Categories      []ProblemCategory `json:"categories"`
	AcceptCount     int64             `json:"accept_count"`     // 通过次数
	SubmissionCount int64             `json:"submission_count"` // 提交总数
	AcceptRate      float64           `json:"accept_rate"`      // 通过率
	UserStatus      *string           `json:"user_status"`      // 用户状态：null-未提交, accepted-已通过, attempted-尝试过
}

// ProblemListResponse 题目列表响应
type ProblemListResponse struct {
	Problems    []ProblemListItem `json:"problems"`
	TotalCount  int64             `json:"total_count"`
	PageSize    int               `json:"page_size"`
	CurrentPage int               `json:"current_page"`
}

// TestCaseUploadRequest 测试用例上传请求
type TestCaseUploadRequest struct {
	ProblemID uint64 `form:"problem_id" binding:"required"`
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name       string  `json:"name" binding:"required,max=30"`
	Color      string  `json:"color" binding:"required,len=7"` // 十六进制颜色值，如 #FF0000
	CategoryID *uint64 `json:"category_id"`                    // 所属分类ID
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name       *string `json:"name" binding:"omitempty,max=30"`
	Color      *string `json:"color" binding:"omitempty,len=7"` // 十六进制颜色值，如 #FF0000
	CategoryID *uint64 `json:"category_id"`                     // 所属分类ID
}

// TagListRequest 标签列表请求
type TagListRequest struct {
	Page       int     `form:"page" binding:"required,min=1"`
	PageSize   int     `form:"page_size" binding:"required,min=1,max=100"`
	Name       string  `form:"name"`        // 标签名称模糊搜索
	CategoryID *uint64 `form:"category_id"` // 分类ID过滤
}

// TagListResponse 标签列表响应
type TagListResponse struct {
	Tags       []TagWithCategory `json:"tags"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	Categories []TagCategory     `json:"categories"` // 一级分类列表
}

// CreateTagCategoryRequest 创建标签分类请求
type CreateTagCategoryRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	Description string  `json:"description" binding:"max=200"`
	ParentID    *uint64 `json:"parent_id"` // 父分类ID，为空表示一级分类
}

// UpdateTagCategoryRequest 更新标签分类请求
type UpdateTagCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
	ParentID    *uint64 `json:"parent_id"` // 父分类ID，为空表示一级分类
}

// GetTagCategoryListRequest 获取标签分类列表请求
type GetTagCategoryListRequest struct {
	ParentID *uint64 `form:"parent_id"` // 父分类ID，为空表示获取一级分类
}

// GetTagCategoryListResponse 获取标签分类列表响应
type GetTagCategoryListResponse struct {
	Categories []TagCategoryDetail `json:"categories"`
}

// SwitchDifficultySystemRequest 切换难度等级系统请求
type SwitchDifficultySystemRequest struct {
	DifficultySystem DifficultySystem `json:"difficulty_system" binding:"required,oneof=normal oi"`
}

// DifficultyMappingRule 难度等级映射规则
type DifficultyMappingRule struct {
	FromSystem DifficultySystem  `json:"from_system"`
	ToSystem   DifficultySystem  `json:"to_system"`
	Mappings   map[string]string `json:"mappings"` // key: 原难度, value: 目标难度
}

var (
	// DefaultDifficultyMappings 默认难度等级映射规则
	DefaultDifficultyMappings = []DifficultyMappingRule{
		{
			FromSystem: DifficultySystemNormal,
			ToSystem:   DifficultySystemOI,
			Mappings: map[string]string{
				DifficultyNormalUnrated: DifficultyOIUnrated,
				DifficultyNormalEasy:    DifficultyOIBeginner,
				DifficultyNormalMedium:  DifficultyOIBasic,
				DifficultyNormalHard:    DifficultyOIAdv,
			},
		},
		{
			FromSystem: DifficultySystemOI,
			ToSystem:   DifficultySystemNormal,
			Mappings: map[string]string{
				DifficultyOIUnrated:   DifficultyNormalUnrated,
				DifficultyOIBeginner:  DifficultyNormalEasy,
				DifficultyOIBasic:     DifficultyNormalEasy,
				DifficultyOIBasicPlus: DifficultyNormalMedium,
				DifficultyOIAdv:       DifficultyNormalMedium,
				DifficultyOIAdvPlus:   DifficultyNormalHard,
				DifficultyOIProv:      DifficultyNormalHard,
				DifficultyOINOI:       DifficultyNormalHard,
			},
		},
	}
)

// GetDifficultyMapping 获取难度等级映射
func GetDifficultyMapping(fromSystem, toSystem DifficultySystem, difficulty string) string {
	// 如果系统相同，直接返回原难度
	if fromSystem == toSystem {
		return difficulty
	}

	// 查找映射规则
	for _, rule := range DefaultDifficultyMappings {
		if rule.FromSystem == fromSystem && rule.ToSystem == toSystem {
			if mapped, ok := rule.Mappings[difficulty]; ok {
				return mapped
			}
			break
		}
	}

	// 如果没有找到映射规则，返回暂无评级
	if toSystem == DifficultySystemNormal {
		return DifficultyNormalUnrated
	}
	return DifficultyOIUnrated
}

// DifficultyLevel 难度等级
type DifficultyLevel struct {
	Code    string `json:"code"`    // 难度代码
	Display string `json:"display"` // 显示名称
}

// DifficultySystemInfo 难度等级系统信息
type DifficultySystemInfo struct {
	System       DifficultySystem  `json:"system"`       // 系统代码
	Name         string            `json:"name"`         // 系统名称
	Difficulties []DifficultyLevel `json:"difficulties"` // 有序的难度等级列表
}

var (
	// OrderedDifficultySystems 有序的难度等级系统配置
	OrderedDifficultySystems = []DifficultySystemInfo{
		{
			System: DifficultySystemNormal,
			Name:   "普通难度",
			Difficulties: []DifficultyLevel{
				{Code: DifficultyNormalUnrated, Display: "暂无评级"},
				{Code: DifficultyNormalEasy, Display: "简单"},
				{Code: DifficultyNormalMedium, Display: "中等"},
				{Code: DifficultyNormalHard, Display: "困难"},
			},
		},
		{
			System: DifficultySystemOI,
			Name:   "OI难度",
			Difficulties: []DifficultyLevel{
				{Code: DifficultyOIUnrated, Display: "暂无评级"},
				{Code: DifficultyOIBeginner, Display: "入门/蒟蒻"},
				{Code: DifficultyOIBasic, Display: "普及-"},
				{Code: DifficultyOIBasicPlus, Display: "普及/提高-"},
				{Code: DifficultyOIAdv, Display: "普及+/提高"},
				{Code: DifficultyOIAdvPlus, Display: "提高+/省选-"},
				{Code: DifficultyOIProv, Display: "省选/NOI-"},
				{Code: DifficultyOINOI, Display: "NOI/NOI+/CTSC/神犇"},
			},
		},
	}
)

// GetDifficultySystemResponse 获取难度等级系统响应
type GetDifficultySystemResponse struct {
	CurrentSystem DifficultySystem       `json:"current_system"` // 当前使用的难度等级系统
	Systems       []DifficultySystemInfo `json:"systems"`        // 所有可用的难度等级系统配置（有序）
}
