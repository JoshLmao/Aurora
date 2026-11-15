package handler

import (
	"time"

	"teslamate-cyberui/internal/repository"
)

// Handler 处理器集合
type Handler struct {
	repo *repository.Repository
}

// NewHandler 创建处理器
func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// Response 通用响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// parseDateTime 解析日期时间字符串，支持多种格式
// 返回 UTC 时间用于数据库查询
func parseDateTime(dateStr string, isEndDate bool) *time.Time {
	if dateStr == "" {
		return nil
	}

	// Plain date YYYY-MM-DD (the main format).
	// Interpreted in the server's timezone (TZ), then converted to UTC.
	location := localLocation()
	if t, err := time.ParseInLocation("2006-01-02", dateStr, location); err == nil {
		utcTime := t.UTC()
		if isEndDate {
			// End date means the last second of that local day. AddDate keeps
			// this correct on days when DST makes the day 23 or 25 hours long.
			endOfDay := t.AddDate(0, 0, 1).Add(-time.Second).UTC()
			return &endOfDay
		}
		return &utcTime
	}

	// Local time without an offset: 2006-01-02T15:04:05.
	// Interpreted in the server's timezone (TZ), then converted to UTC.
	if t, err := time.ParseInLocation("2006-01-02T15:04:05", dateStr, location); err == nil {
		utcTime := t.UTC()
		return &utcTime
	}

	// 兼容 UTC 格式，直接使用
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return &t
	}
	if t, err := time.Parse("2006-01-02T15:04:05.000Z", dateStr); err == nil {
		return &t
	}
	// 尝试带毫秒的 ISO 格式
	if t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", dateStr); err == nil {
		return &t
	}

	return nil
}

// localLocation returns the server's configured timezone. Go derives
// time.Local from the TZ environment variable (UTC when it is unset).
func localLocation() *time.Location {
	return time.Local
}
