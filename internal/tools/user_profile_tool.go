package tools

import (
	"context"
	"encoding/json" // 用于将结果序列化为 JSON 字符串
	"log/slog"

	"xiaoyuzhoufm-mcp/internal/xyzclient"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetUserProfileByIDHandler 是一个工具处理函数，用于获取用户的个人资料。
func GetUserProfileByIDHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Debug("Executing get_user_profile_by_id tool", "arguments", request.Params.Arguments)

	userID, ok := request.Params.Arguments["user_id"].(string)
	if !ok || userID == "" {
		return mcp.NewToolResultError("输入参数 'user_id' 不能为空且必须是字符串类型。"), nil
	}

	profileData, err := xyzclient.GetUserProfileByID(userID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("调用API获取用户 Profile 失败", err), nil
	}

	// 将 profileData (结构体指针) 序列化为 JSON 字符串以放入 TextContent
	profileJSON, err := json.Marshal(profileData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("处理结果失败", err), nil
	}

	slog.Debug("成功获取用户 Profile", "userID", userID)

	return mcp.NewToolResultText(string(profileJSON)), nil
}

// GetUserStatsHandler 是一个工具处理函数，用于获取用户的统计数据。
func GetUserStatsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Debug("Executing get_user_stats tool", "arguments", request.Params.Arguments)

	userID, ok := request.Params.Arguments["user_id"].(string)
	if !ok || userID == "" {
		return mcp.NewToolResultError("输入参数 'user_id' 不能为空且必须是字符串类型。"), nil
	}

	statsData, err := xyzclient.GetUserStats(userID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("调用API获取用户 Stats 失败", err), nil
	}

	statsJSON, err := json.Marshal(statsData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("处理 Stats 结果失败", err), nil
	}

	slog.Debug("成功获取用户 Stats", "userID", userID)

	return mcp.NewToolResultText(string(statsJSON)), nil
}
