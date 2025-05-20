package tools

import (
	"context"
	"encoding/json"
	"log/slog"

	"xiaoyuzhoufm-mcp/internal/xyzclient"

	"github.com/mark3labs/mcp-go/mcp" // Added for mcp types
)

// GetPodcastDetailsHandler is the MCP handler function for the GetPodcastDetailsTool.
func GetPodcastDetailsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Debug("Executing get_podcast_details tool", "arguments", request.Params.Arguments)

	podcastID, ok := request.Params.Arguments["podcast_id"].(string)
	if !ok || podcastID == "" {
		return mcp.NewToolResultError("输入参数 'podcast_id' 不能为空且必须是字符串类型。"), nil
	}

	podcastDetailsData, err := xyzclient.GetPodcastDetailsByID(podcastID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("调用API获取 PodcastDetails 失败", err), nil
	}

	podcastDetailsJSON, err := json.Marshal(podcastDetailsData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("处理结果失败", err), nil
	}

	slog.Debug("成功获取 PodcastDetails", "podcast_id", podcastID)
	return mcp.NewToolResultText(string(podcastDetailsJSON)), nil
}

// ListPodcastEpisodesHandler is the MCP handler function for the ListPodcastEpisodesTool.
func ListPodcastEpisodesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Debug("Executing list_podcast_episodes tool", "arguments", request.Params.Arguments)

	podcastID, ok := request.Params.Arguments["podcast_id"].(string)
	if !ok || podcastID == "" {
		return mcp.NewToolResultError("错误: 输入参数 'podcast_id' 不能为空且必须是字符串类型。"), nil
	}

	apiRequest := xyzclient.EpisodeListRequest{
		PID:   podcastID,
		Limit: 20,
	}

	if order, ok := request.Params.Arguments["order"].(string); ok && order != "" {
		if order != "asc" && order != "desc" {
			return mcp.NewToolResultError("错误: 输入参数 'order' 必须是 'asc' 或 'desc'。"), nil
		}
		apiRequest.Order = order
	} else {
		apiRequest.Order = "desc" // Default value
	}

	if lmkMap, ok := request.Params.Arguments["load_more_key"].(map[string]interface{}); ok && lmkMap != nil {
		lmk := &xyzclient.LoadMoreKey{}
		if direction, ok := lmkMap["direction"].(string); ok {
			lmk.Direction = direction
		}
		if pubDate, ok := lmkMap["pubDate"].(string); ok {
			lmk.PubDate = pubDate
		}
		if id, ok := lmkMap["id"].(string); ok {
			lmk.ID = id
		}
		// Only assign if at least one field was present, or if the object itself was non-nil
		// The xyzclient.ListPodcastEpisodes will handle nil LoadMoreKey if no fields are set.
		if lmk.Direction != "" || lmk.PubDate != "" || lmk.ID != "" {
			apiRequest.LoadMoreKey = lmk
		}
	}

	slog.Debug("Constructed API request for ListPodcastEpisodes", "apiRequest", apiRequest)

	episodeListData, err := xyzclient.ListPodcastEpisodes(apiRequest)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("调用API获取播客单集列表失败", err), nil
	}

	episodeListJSON, err := json.Marshal(episodeListData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("处理结果失败", err), nil
	}
	slog.Debug("成功获取播客单集列表", "podcast_id", podcastID, "count", len(episodeListData.Data))

	return mcp.NewToolResultText(string(episodeListJSON)), nil
}

// GetEpisodeDetailsHandler is the MCP handler function for the GetEpisodeDetailsTool.
func GetEpisodeDetailsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Debug("Executing get_episode_details tool", "arguments", request.Params.Arguments)

	episodeID, ok := request.Params.Arguments["episode_id"].(string)
	if !ok || episodeID == "" {
		return mcp.NewToolResultError("输入参数 'episode_id' 不能为空且必须是字符串类型。"), nil
	}

	episodeDetailsData, err := xyzclient.GetEpisodeDetailsByID(episodeID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("调用API获取单集详情失败", err), nil
	}

	episodeDetailsJSON, err := json.Marshal(episodeDetailsData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("处理结果失败", err), nil
	}
	slog.Debug("成功获取单集详情", "episode_id", episodeID, "title", episodeDetailsData.Title)

	return mcp.NewToolResultText(string(episodeDetailsJSON)), nil
}
