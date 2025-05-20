package tools

import (
	"context"
	"encoding/json"
	"log/slog"

	"xiaoyuzhoufm-mcp/internal/xyzclient"

	"github.com/mark3labs/mcp-go/mcp"
)

// SearchPodcastsHandler is the MCP handler function for the search_podcasts tool.
func SearchPodcastsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Debug("Executing search_podcasts tool", "arguments", request.Params.Arguments)

	keyword, ok := request.Params.Arguments["keyword"].(string)
	if !ok || keyword == "" {
		return mcp.NewToolResultError("参数 'keyword' 不能为空且必须是字符串类型。"), nil
	}

	var loadMoreKey *xyzclient.SearchAPILoadMoreKey
	if lmkMap, ok := request.Params.Arguments["load_more_key"].(map[string]interface{}); ok && lmkMap != nil {
		lmk := &xyzclient.SearchAPILoadMoreKey{}
		if lmkVal, okGet := lmkMap["loadMoreKey"]; okGet { // API uses "loadMoreKey" (interface{}) inside the object
			lmk.LoadMoreKey = lmkVal
		}
		if searchID, okGet := lmkMap["searchId"].(string); okGet {
			lmk.SearchID = searchID
		}
		// Only assign if at least one field was present, or if the object itself was non-nil
		if lmk.LoadMoreKey != nil || lmk.SearchID != "" {
			loadMoreKey = lmk
		}
	}

	searchResult, err := xyzclient.SearchPodcasts(keyword, loadMoreKey)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("调用API搜索播客失败", err), nil
	}

	resultJSON, err := json.Marshal(searchResult)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("处理播客搜索结果失败", err), nil
	}
	slog.Debug("成功搜索播客", "keyword", keyword, "count", len(searchResult.Data))
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// SearchEpisodesHandler is the MCP handler function for the search_episodes tool.
func SearchEpisodesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Debug("Executing search_episodes tool", "arguments", request.Params.Arguments)

	keyword, ok := request.Params.Arguments["keyword"].(string)
	if !ok || keyword == "" {
		return mcp.NewToolResultError("参数 'keyword' 不能为空且必须是字符串类型。"), nil
	}

	pid, _ := request.Params.Arguments["pid"].(string) // pid is optional for this tool

	var loadMoreKey *xyzclient.SearchAPILoadMoreKey
	if lmkMap, ok := request.Params.Arguments["load_more_key"].(map[string]interface{}); ok && lmkMap != nil {
		lmk := &xyzclient.SearchAPILoadMoreKey{}
		if lmkVal, okGet := lmkMap["loadMoreKey"]; okGet {
			lmk.LoadMoreKey = lmkVal
		}
		if searchID, okGet := lmkMap["searchId"].(string); okGet {
			lmk.SearchID = searchID
		}
		if lmk.LoadMoreKey != nil || lmk.SearchID != "" {
			loadMoreKey = lmk
		}
	}

	searchResult, err := xyzclient.SearchEpisodes(keyword, pid, loadMoreKey)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("调用API搜索单集失败", err), nil
	}

	resultJSON, err := json.Marshal(searchResult)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("处理单集搜索结果失败", err), nil
	}
	slog.Debug("成功搜索单集", "keyword", keyword, "pid", pid, "count", len(searchResult.Data))
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// SearchUsersHandler is the MCP handler function for the search_users tool.
func SearchUsersHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Debug("Executing search_users tool", "arguments", request.Params.Arguments)

	keyword, ok := request.Params.Arguments["keyword"].(string)
	if !ok || keyword == "" {
		return mcp.NewToolResultError("参数 'keyword' 不能为空且必须是字符串类型。"), nil
	}

	var loadMoreKey *xyzclient.SearchAPILoadMoreKey
	if lmkMap, ok := request.Params.Arguments["load_more_key"].(map[string]interface{}); ok && lmkMap != nil {
		lmk := &xyzclient.SearchAPILoadMoreKey{}
		if lmkVal, okGet := lmkMap["loadMoreKey"]; okGet {
			lmk.LoadMoreKey = lmkVal
		}
		if searchID, okGet := lmkMap["searchId"].(string); okGet {
			lmk.SearchID = searchID
		}
		if lmk.LoadMoreKey != nil || lmk.SearchID != "" {
			loadMoreKey = lmk
		}
	}

	searchResult, err := xyzclient.SearchUsers(keyword, loadMoreKey)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("调用API搜索用户失败", err), nil
	}

	resultJSON, err := json.Marshal(searchResult)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("处理用户搜索结果失败", err), nil
	}
	slog.Debug("成功搜索用户", "keyword", keyword, "count", len(searchResult.Data))
	return mcp.NewToolResultText(string(resultJSON)), nil
}
