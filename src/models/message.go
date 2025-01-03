package models

import "time"

// MessageAction 消息可执行的操作
type MessageAction struct {
	Action     string `json:"action"`      // 操作代码
	Name       string `json:"name"`        // 操作名称
	Type       string `json:"type"`        // 操作类型：primary, warning, danger 等
	NeedReason bool   `json:"need_reason"` // 是否需要填写原因
}

// Message 站内信
type Message struct {
	ID            uint64          `json:"id"`
	SenderID      *uint64         `json:"sender_id"` // 为空表示系统消息
	ReceiverID    uint64          `json:"receiver_id"`
	Type          string          `json:"type"`
	Title         string          `json:"title"`
	Content       string          `json:"content"`
	IsRead        bool            `json:"is_read"`
	IsProcessed   bool            `json:"is_processed"`             // 是否已处理（用于申请类消息）
	ApplicationID uint64          `json:"application_id,omitempty"` // 关联的申请ID
	ReadAt        *time.Time      `json:"read_at,omitempty"`        // 阅读时间
	CreatedAt     time.Time       `json:"created_at"`
	Actions       []MessageAction `json:"actions,omitempty" gorm:"-"` // 可执行的操作列表
}

// MessageListRequest 站内信列表请求
type MessageListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Type     string `form:"type"`    // 消息类型，为空表示所有类型
	IsRead   *bool  `form:"is_read"` // 是否已读，为空表示所有状态
}

// MessageListResponse 站内信列表响应
type MessageListResponse struct {
	Messages    []Message `json:"messages"`
	Total       int64     `json:"total"`
	Page        int       `json:"page"`
	PageSize    int       `json:"page_size"`
	UnreadCount int64     `json:"unread_count"` // 总未读数
}

// TeamApplication 团队申请
type TeamApplication struct {
	ID        uint64    `json:"id"`
	TeamID    uint64    `json:"team_id"`
	UserID    uint64    `json:"user_id"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TeamApplicationRequest 团队申请请求
type TeamApplicationRequest struct {
	TeamID  uint64 `json:"team_id" binding:"required"`
	Message string `json:"message"`
}

// TeamApplicationResponse 团队申请响应
type TeamApplicationResponse struct {
	Application TeamApplication `json:"application"`
	Team        Team            `json:"team"`
	User        struct {
		ID       uint64 `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user"`
}

// TeamApplicationListRequest 团队申请列表请求
type TeamApplicationListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	TeamID   uint64 `form:"team_id"`
	Status   string `form:"status"`
}

// TeamApplicationListResponse 团队申请列表响应
type TeamApplicationListResponse struct {
	Applications []TeamApplicationResponse `json:"applications"`
	Total        int64                     `json:"total"`
	Page         int                       `json:"page"`
	PageSize     int                       `json:"page_size"`
}

// TeamApplicationHandleRequest 处理团队申请请求
type TeamApplicationHandleRequest struct {
	ApplicationID uint64 `json:"application_id" binding:"required"`
	Status        string `json:"status" binding:"required,oneof=approved rejected"`
	Message       string `json:"message"` // 审批意见
}

// BatchReadRequest 批量标记已读请求
type BatchReadRequest struct {
	MessageIDs []uint64 `json:"message_ids" binding:"required,min=1"`
}

// UnreadCountResponse 未读消息数量响应
type UnreadCountResponse struct {
	Total   int64                    `json:"total"`    // 总未读数
	ByTypes map[string]TypeCountInfo `json:"by_types"` // 按类型统计的未读数
}

// TypeCountInfo 消息类型统计信息
type TypeCountInfo struct {
	Count       int64  `json:"count"`       // 未读数量
	Description string `json:"description"` // 类型描述
}

// 预定义消息类型
const (
	MessageTypeSystem          = "system"           // 系统消息
	MessageTypeTeamApplication = "team_application" // 团队申请
	MessageTypeTeamInvitation  = "team_invitation"  // 团队邀请
	MessageTypeTeamNotice      = "team_notice"      // 团队通知
)

// GetMessageTypeDescription 获取消息类型描述
func GetMessageTypeDescription(msgType string) string {
	switch msgType {
	case MessageTypeSystem:
		return "系统消息"
	case MessageTypeTeamApplication:
		return "团队申请"
	case MessageTypeTeamInvitation:
		return "团队邀请"
	case MessageTypeTeamNotice:
		return "团队通知"
	default:
		return "其他消息"
	}
}

// 预定义消息操作
const (
	ActionMarkRead = "mark_read" // 标记已读
	ActionApprove  = "approve"   // 通过
	ActionReject   = "reject"    // 拒绝
	ActionDelete   = "delete"    // 删除
)

// GetMessageActions 获取消息可执行的操作列表
func GetMessageActions(msg *Message, isAdmin bool) []MessageAction {
	actions := make([]MessageAction, 0)

	// 根据消息类型添加特定操作
	switch msg.Type {
	case MessageTypeTeamApplication:
		// 只有管理员可以处理未处理的申请
		if isAdmin && !msg.IsProcessed {
			actions = append(actions, MessageAction{
				Action:     ActionApprove,
				Name:       "通过",
				Type:       "primary",
				NeedReason: true,
			}, MessageAction{
				Action:     ActionReject,
				Name:       "拒绝",
				Type:       "danger",
				NeedReason: true,
			})
		}
		return actions
	default:
		// 其他类型的消息
		// 未读消息可以标记已读
		if !msg.IsRead {
			actions = append(actions, MessageAction{
				Action: ActionMarkRead,
				Name:   "标记已读",
				Type:   "default",
			})
		}
		// 所有非申请类消息都可以删除
		actions = append(actions, MessageAction{
			Action: ActionDelete,
			Name:   "删除",
			Type:   "danger",
		})
	}

	return actions
}
