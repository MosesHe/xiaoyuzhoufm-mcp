package server

import (
	"log/slog"

	"xiaoyuzhoufm-mcp/internal/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RunStdioServer initializes and runs a basic MCP server over stdio.
func RunStdioServer() {
	s := server.NewMCPServer(
		"XiaoyuzhouFM Stdio Server", // Server name
		"0.0.1",                     // Server version
		server.WithLogging(),        // Optional: enable basic logging
	)

	getUserProfileByIDTool := mcp.NewTool("get_user_profile_by_id",
		mcp.WithDescription("根据用户 UID 获取指定用户的公开信息。"),
		mcp.WithString("user_id",
			mcp.Description("要查询的用户的唯一标识符 (UID)。"),
			mcp.Required(),
		),
	)
	s.AddTool(getUserProfileByIDTool, tools.GetUserProfileByIDHandler)

	getUserStatsTool := mcp.NewTool("get_user_stats",
		mcp.WithDescription("获取指定用户的统计数据（如关注数、粉丝数、订阅播客数、收听时长）。"),
		mcp.WithString("user_id",
			mcp.Description("要查询的用户的唯一标识符 (UID)。"),
			mcp.Required(),
		),
	)
	s.AddTool(getUserStatsTool, tools.GetUserStatsHandler)

	podcastDetailsTool := mcp.NewTool("get_podcast_details",
		mcp.WithDescription("获取指定播客的详细信息。"),
		mcp.WithString("podcast_id",
			mcp.Description("播客的唯一标识符 (PID)。"),
			mcp.Required(),
		),
	)
	s.AddTool(podcastDetailsTool, tools.GetPodcastDetailsHandler)

	listPodcastEpisodesTool := mcp.NewTool("list_podcast_episodes",
		mcp.WithDescription("获取指定播客的单集列表。"),
		mcp.WithString("podcast_id",
			mcp.Description("播客的唯一标识符 (PID)。"),
			mcp.Required(),
		),
		mcp.WithString("order",
			mcp.Description("排序方式。"),
			mcp.Enum("asc", "desc"),
		),
		mcp.WithObject("load_more_key",
			mcp.Description("用于分页查询的键。对应 API 中的 loadMoreKey 对象。"), // Description for the object itself
			// The load_more_key object itself is optional.
			mcp.Properties(map[string]interface{}{ // Corrected: Pass a map to mcp.Properties
				"direction": map[string]interface{}{
					"type":        "string",
					"description": "分页方向 (例如 'NEXT')。",
				},
				"pubDate": map[string]interface{}{
					"type":        "string",
					"description": "分页锚点的发布日期 (ISO 8601 格式的字符串)。",
				},
				"id": map[string]interface{}{
					"type":        "string",
					"description": "分页锚点的单集ID。",
				},
			}),
		),
	)
	s.AddTool(listPodcastEpisodesTool, tools.ListPodcastEpisodesHandler)

	getEpisodeDetailsTool := mcp.NewTool("get_episode_details",
		mcp.WithDescription("获取指定单集的详细信息。"),
		mcp.WithString("episode_id",
			mcp.Description("要查询的单集的唯一标识符 (EID)。"),
			mcp.Required(),
		),
	)
	s.AddTool(getEpisodeDetailsTool, tools.GetEpisodeDetailsHandler)

	// Search Podcasts Tool
	searchPodcastsTool := mcp.NewTool("search_podcasts",
		mcp.WithDescription("根据关键词搜索播客。"),
		mcp.WithString("keyword",
			mcp.Description("搜索关键词。"),
			mcp.Required(),
		),
		mcp.WithObject("load_more_key",
			mcp.Description("用于分页查询的键，存在于先前搜索请求的响应中。"),
			// load_more_key object itself is optional
			mcp.Properties(map[string]interface{}{
				// API example shows loadMoreKey as an integer for podcast/episode/user search results
				"loadMoreKey": map[string]interface{}{
					"type":        "integer",
					"description": "分页加载键（通常是数字）。",
				},
				"searchId": map[string]interface{}{
					"type":        "string",
					"description": "搜索会话ID。",
				},
			}),
		),
	)
	s.AddTool(searchPodcastsTool, tools.SearchPodcastsHandler)

	// Search Episodes Tool
	searchEpisodesTool := mcp.NewTool("search_episodes",
		mcp.WithDescription("根据关键词搜索单集。可选择在特定播客 (pid) 内搜索。"),
		mcp.WithString("keyword",
			mcp.Description("搜索关键词。"),
			mcp.Required(),
		),
		mcp.WithString("pid", // Optional podcast ID to search within
			mcp.Description("可选参数，如果需要在特定播客内搜索单集，请提供播客ID。"),
			// This parameter is optional, so no mcp.Required()
		),
		mcp.WithObject("load_more_key",
			mcp.Description("用于分页查询的键，存在于先前搜索请求的响应中。"),
			mcp.Properties(map[string]interface{}{
				// API example shows loadMoreKey as an integer for podcast/episode/user search results
				"loadMoreKey": map[string]interface{}{
					"type":        "integer",
					"description": "分页加载键（通常是数字）。",
				},
				"searchId": map[string]interface{}{
					"type":        "string",
					"description": "搜索会话ID。",
				},
			}),
		),
	)
	s.AddTool(searchEpisodesTool, tools.SearchEpisodesHandler)

	// Search Users Tool
	searchUsersTool := mcp.NewTool("search_users",
		mcp.WithDescription("根据关键词搜索用户。"),
		mcp.WithString("keyword",
			mcp.Description("搜索关键词。"),
			mcp.Required(),
		),
		mcp.WithObject("load_more_key",
			mcp.Description("用于分页查询的键，存在于先前搜索请求的响应中。"),
			mcp.Properties(map[string]interface{}{
				// API example shows loadMoreKey as an integer for podcast/episode/user search results
				"loadMoreKey": map[string]interface{}{
					"type":        "integer",
					"description": "分页加载键（通常是数字）。",
				},
				"searchId": map[string]interface{}{
					"type":        "string",
					"description": "搜索会话ID。",
				},
			}),
		),
	)
	s.AddTool(searchUsersTool, tools.SearchUsersHandler)

	slog.Debug("MCP Stdio Server starting with 'hello', 'get_user_profile_by_id', 'get_user_stats', 'get_podcast_details', 'list_podcast_episodes', 'get_episode_details', 'search_podcasts', 'search_episodes', and 'search_users' tools...")

	if err := server.ServeStdio(s); err != nil {
		slog.Error("MCP Stdio Server failed", "error", err)
	}

	slog.Debug("MCP Stdio Server stopped.")
}
