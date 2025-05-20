package xyzclient

// sendCodeRequestBody defines the structure for the /v1/auth/sendCode API request.
type sendCodeRequestBody struct {
	MobilePhoneNumber string `json:"mobilePhoneNumber"`
	AreaCode          string `json:"areaCode"`
}

// loginOrSignUpWithSMSRequestBody defines the structure for the /v1/auth/loginOrSignUpWithSMS API request.
type loginOrSignUpWithSMSRequestBody struct {
	AreaCode          string `json:"areaCode"`
	VerifyCode        string `json:"verifyCode"`
	MobilePhoneNumber string `json:"mobilePhoneNumber"`
}

// LoginAPIResponse directly maps to the Xiaoyuzhou FM login API success response body.
type LoginAPIResponse struct {
	Data struct {
		User struct {
			UID      string `json:"uid"`
			Nickname string `json:"nickname"`
		} `json:"user"`
	} `json:"data"`
}

// RefreshTokenAPIResponse directly maps to the Xiaoyuzhou FM refresh token API success response body.
type RefreshTokenAPIResponse struct {
	Success           bool   `json:"success"`
	XJikeAccessToken  string `json:"x-jike-access-token"`
	XJikeRefreshToken string `json:"x-jike-refresh-token"`
}

// Picture Structure (used in Avatar and Authorship Image)
type Picture struct {
	PicURL       string `json:"picUrl"`
	LargePicURL  string `json:"largePicUrl"`
	MiddlePicURL string `json:"middlePicUrl"`
	SmallPicURL  string `json:"smallPicUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Format       string `json:"format"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// Avatar Structure (used in UserProfileData and PodcastAuthor)
type Avatar struct {
	Picture          Picture `json:"picture"`
	AvatarDecoration string  `json:"avatarDecoration,omitempty"` // e.g., "PLATFORM_MEMBER_L1"
}

// PodcastAuthor Structure (used in AuthorshipEntry)
type PodcastAuthor struct {
	Type              string      `json:"type"` // e.g., "USER"
	UID               string      `json:"uid"`
	Avatar            Avatar      `json:"avatar"`
	Nickname          string      `json:"nickname"`
	IsNicknameSet     bool        `json:"isNicknameSet"`
	Bio               string      `json:"bio"`
	Gender            string      `json:"gender,omitempty"` // MALE, FEMALE, THIRD
	IsCancelled       bool        `json:"isCancelled"`
	ReadTrackInfo     interface{} `json:"readTrackInfo"` // Empty object {}
	IPLoc             string      `json:"ipLoc"`
	Relation          string      `json:"relation"` // e.g., "STRANGE"
	IsBlockedByViewer bool        `json:"isBlockedByViewer"`
}

// PodcastPermission Structure (used in AuthorshipEntry)
type PodcastPermission struct {
	Name   string `json:"name"`   // e.g., "SHARE", "AI_SUMMARIZE_EPISODE"
	Status string `json:"status"` // e.g., "PERMITTED"
}

// PodcastColor Structure (used in AuthorshipEntry)
type PodcastColor struct {
	Original string `json:"original"`
	Light    string `json:"light"`
	Dark     string `json:"dark"`
}

// AuthorshipEntry Structure (used in UserProfileData)
type AuthorshipEntry struct {
	Type                     string              `json:"type"` // e.g., "PODCAST"
	PID                      string              `json:"pid"`
	Title                    string              `json:"title"`
	Author                   string              `json:"author"`
	Brief                    string              `json:"brief"`
	Description              string              `json:"description"`
	SubscriptionCount        int                 `json:"subscriptionCount"`
	Image                    Picture             `json:"image"` // Reusing Picture struct
	Color                    PodcastColor        `json:"color"`
	HasTopic                 bool                `json:"hasTopic"`
	TopicLabels              []string            `json:"topicLabels"`
	SyncMode                 string              `json:"syncMode"` // e.g., "SELF_HOSTING"
	EpisodeCount             int                 `json:"episodeCount"`
	LatestEpisodePubDate     string              `json:"latestEpisodePubDate"` // ISO Date string
	SubscriptionStatus       string              `json:"subscriptionStatus"`   // e.g., "ON", "OFF"
	SubscriptionPush         bool                `json:"subscriptionPush"`
	SubscriptionPushPriority string              `json:"subscriptionPushPriority"` // e.g., "HIGH"
	SubscriptionStar         bool                `json:"subscriptionStar"`
	Status                   string              `json:"status"` // e.g., "NORMAL"
	Permissions              []PodcastPermission `json:"permissions"`
	PayType                  string              `json:"payType"` // e.g., "FREE"
	PayEpisodeCount          int                 `json:"payEpisodeCount"`
	Podcasters               []PodcastAuthor     `json:"podcasters"`
	ReadTrackInfo            interface{}         `json:"readTrackInfo"` // Empty object {}
	HasPopularEpisodes       bool                `json:"hasPopularEpisodes"`
	Contacts                 []interface{}       `json:"contacts"` // Empty array []
	IsCustomized             bool                `json:"isCustomized"`
	ShowZhuiguangIcon        bool                `json:"showZhuiguangIcon"`
}

// CertificationShow Structure (used in CertificationEntry)
type CertificationShow struct {
	Title string `json:"title"`
}

// CertificationEntry Structure (used in UserProfileData)
type CertificationEntry struct {
	Kind  string              `json:"kind"` // e.g., "PODCASTER"
	Shows []CertificationShow `json:"shows"`
}

// UserProfileData Structure (corresponds to the "data" object in API response)
type UserProfileData struct {
	Type              string               `json:"type"` // e.g., "USER"
	UID               string               `json:"uid"`
	Avatar            Avatar               `json:"avatar"`
	Nickname          string               `json:"nickname"`
	IsNicknameSet     bool                 `json:"isNicknameSet"`
	Bio               string               `json:"bio"`
	Gender            string               `json:"gender"` // e.g., "MALE"
	IsCancelled       bool                 `json:"isCancelled"`
	ReadTrackInfo     interface{}          `json:"readTrackInfo"` // Empty object {}
	IPLoc             string               `json:"ipLoc"`
	Relation          string               `json:"relation"` // e.g., "STRANGE"
	IsBlockedByViewer bool                 `json:"isBlockedByViewer"`
	IsInvited         bool                 `json:"isInvited"`
	Authorship        []AuthorshipEntry    `json:"authorship"`
	Certifications    []CertificationEntry `json:"certifications"`
}

// UserProfileAPIResponse Structure (to wrap the main "data" object)
type UserProfileAPIResponse struct {
	Data UserProfileData `json:"data"`
}

// UserStatsData represents the user statistics data from the API.
type UserStatsData struct {
	FollowerCount      int `json:"followerCount"`
	FollowingCount     int `json:"followingCount"`
	SubscriptionCount  int `json:"subscriptionCount"`
	TotalPlayedSeconds int `json:"totalPlayedSeconds"`
}

// UserStatsAPIResponse wraps the UserStatsData as per the API's structure.
type UserStatsAPIResponse struct {
	Data UserStatsData `json:"data"`
}

// PodcastContact defines a contact method for a podcast.
type PodcastContact struct {
	Name string `json:"name"`
	Note string `json:"note,omitempty"`
	Type string `json:"type"`
	URL  string `json:"url,omitempty"`
}

// PodcastImage defines various image URLs for a podcast.
type PodcastImage struct {
	LargePicURL  string `json:"largePicUrl"`
	MiddlePicURL string `json:"middlePicUrl"`
	PicURL       string `json:"picUrl"`
	SmallPicURL  string `json:"smallPicUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

// PodcasterAvatarPicture defines the structure for a podcaster's avatar picture.
type PodcasterAvatarPicture struct {
	Format       string `json:"format"`
	Height       int    `json:"height"`
	LargePicURL  string `json:"largePicUrl"`
	MiddlePicURL string `json:"middlePicUrl"`
	PicURL       string `json:"picUrl"`
	SmallPicURL  string `json:"smallPicUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Width        int    `json:"width"`
}

// PodcasterAvatar defines the avatar for a podcaster.
type PodcasterAvatar struct {
	Picture PodcasterAvatarPicture `json:"picture"`
}

// Podcaster defines information about a podcaster.
type Podcaster struct {
	Avatar            PodcasterAvatar        `json:"avatar"`
	Bio               string                 `json:"bio"`
	Gender            string                 `json:"gender,omitempty"`
	IPLoc             string                 `json:"ipLoc"`
	IsBlockedByViewer bool                   `json:"isBlockedByViewer"`
	IsCancelled       bool                   `json:"isCancelled"`
	IsNicknameSet     bool                   `json:"isNicknameSet"`
	Nickname          string                 `json:"nickname"`
	ReadTrackInfo     map[string]interface{} `json:"readTrackInfo"` // Assuming it can be an empty object or have arbitrary keys
	Relation          string                 `json:"relation"`
	Type              string                 `json:"type"`
	UID               string                 `json:"uid"`
}

// PodcastDetailData represents the detailed information of a podcast from the API.
type PodcastDetailData struct {
	Author                   string                 `json:"author"`
	Brief                    string                 `json:"brief"`
	Color                    PodcastColor           `json:"color"`
	Contacts                 []PodcastContact       `json:"contacts"`
	Description              string                 `json:"description"`
	EpisodeCount             int                    `json:"episodeCount"`
	HasPopularEpisodes       bool                   `json:"hasPopularEpisodes"`
	Image                    PodcastImage           `json:"image"`
	IsCustomized             bool                   `json:"isCustomized"`
	LatestEpisodePubDate     string                 `json:"latestEpisodePubDate"` // Consider time.Time if parsing is needed
	PayEpisodeCount          int                    `json:"payEpisodeCount"`
	PayType                  string                 `json:"payType"`
	Permissions              []PodcastPermission    `json:"permissions"`
	PID                      string                 `json:"pid"`
	Podcasters               []Podcaster            `json:"podcasters"`
	ReadTrackInfo            map[string]interface{} `json:"readTrackInfo"` // Assuming it can be an empty object or have arbitrary keys
	Status                   string                 `json:"status"`
	SubscriptionCount        int                    `json:"subscriptionCount"`
	SubscriptionPush         bool                   `json:"subscriptionPush"`
	SubscriptionPushPriority string                 `json:"subscriptionPushPriority"`
	SubscriptionStar         bool                   `json:"subscriptionStar"`
	SubscriptionStatus       string                 `json:"subscriptionStatus"`
	SyncMode                 string                 `json:"syncMode"`
	Title                    string                 `json:"title"`
	TopicLabels              []string               `json:"topicLabels"`
	Type                     string                 `json:"type"`
}

// PodcastDetailAPIResponse wraps the PodcastDetailData as per the API's structure.
type PodcastDetailAPIResponse struct {
	Data PodcastDetailData `json:"data"`
}

// LoadMoreKey defines the structure for pagination keys.
type LoadMoreKey struct {
	Direction string `json:"direction,omitempty"`
	PubDate   string `json:"pubDate,omitempty"`
	ID        string `json:"id,omitempty"`
}

// EpisodeListRequest defines the request body for listing podcast episodes.
type EpisodeListRequest struct {
	PID         string       `json:"pid"`
	Order       string       `json:"order,omitempty"` // "asc" or "desc"
	Limit       int          `json:"limit,omitempty"`
	LoadMoreKey *LoadMoreKey `json:"loadMoreKey,omitempty"`
}

// EnclosureInfo defines the structure for episode media enclosure.
type EnclosureInfo struct {
	URL string `json:"url"`
}

// SourceInfo defines the structure for media source.
type SourceInfo struct {
	Mode string `json:"mode"`
	URL  string `json:"url"`
}

// MediaInfo defines the structure for episode media details.
type MediaInfo struct {
	ID       string     `json:"id"`
	Size     int64      `json:"size"` // Changed to int64 for potentially large file sizes
	MimeType string     `json:"mimeType"`
	Source   SourceInfo `json:"source"`
}

// PodcastSummary defines the summary of a podcast, often nested within an episode.
// This is similar to AuthorshipEntry but tailored for the episode context.
type PodcastSummary struct {
	Type                     string              `json:"type"`
	PID                      string              `json:"pid"`
	Title                    string              `json:"title"`
	Author                   string              `json:"author"`
	Brief                    string              `json:"brief"`
	Description              string              `json:"description"`
	SubscriptionCount        int                 `json:"subscriptionCount"`
	Image                    Picture             `json:"image"` // Reusing Picture struct
	Color                    PodcastColor        `json:"color"` // Reusing PodcastColor struct
	HasTopic                 bool                `json:"hasTopic"`
	TopicLabels              []string            `json:"topicLabels"`
	SyncMode                 string              `json:"syncMode"`
	LatestEpisodePubDate     string              `json:"latestEpisodePubDate"`
	SubscriptionStatus       string              `json:"subscriptionStatus"`
	SubscriptionPush         bool                `json:"subscriptionPush"`
	SubscriptionPushPriority string              `json:"subscriptionPushPriority"`
	SubscriptionStar         bool                `json:"subscriptionStar"`
	Status                   string              `json:"status"`
	EpisodeCount             int                 `json:"episodeCount"`
	Permissions              []PodcastPermission `json:"permissions"` // Reusing PodcastPermission
	PayType                  string              `json:"payType"`
	PayEpisodeCount          int                 `json:"payEpisodeCount"`
	IsCustomized             bool                `json:"isCustomized"`
	Podcasters               []PodcastAuthor     `json:"podcasters"` // Reusing PodcastAuthor
	HasPopularEpisodes       bool                `json:"hasPopularEpisodes"`
	Contacts                 []PodcastContact    `json:"contacts"` // Reusing PodcastContact
	PlayTime                 int64               `json:"playTime"`
	ShowZhuiguangIcon        bool                `json:"showZhuiguangIcon"`
}

// WechatShareInfo defines the structure for WeChat sharing details.
type WechatShareInfo struct {
	Style string `json:"style"`
}

// TranscriptInfo defines the structure for episode transcript information.
type TranscriptInfo struct {
	MediaID string `json:"mediaId"`
}

// Episode defines the structure for a single podcast episode.
type Episode struct {
	Type           string              `json:"type"`
	EID            string              `json:"eid"`
	PID            string              `json:"pid"`
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	Shownotes      string              `json:"shownotes"`
	Duration       int                 `json:"duration"`
	Image          Picture             `json:"image"` // Reusing Picture struct
	Enclosure      EnclosureInfo       `json:"enclosure"`
	IsPrivateMedia bool                `json:"isPrivateMedia"`
	MediaKey       string              `json:"mediaKey"`
	Media          MediaInfo           `json:"media"`
	PlayCount      int                 `json:"playCount"`
	ClapCount      int                 `json:"clapCount"`
	CommentCount   int                 `json:"commentCount"`
	FavoriteCount  int                 `json:"favoriteCount"`
	PubDate        string              `json:"pubDate"`
	Status         string              `json:"status"`
	Podcast        PodcastSummary      `json:"podcast"`
	IsPlayed       bool                `json:"isPlayed"`
	IsFinished     bool                `json:"isFinished"`
	IsFavorited    bool                `json:"isFavorited"`
	IsPicked       bool                `json:"isPicked"`
	Permissions    []PodcastPermission `json:"permissions"` // Reusing PodcastPermission
	PayType        string              `json:"payType"`
	WechatShare    WechatShareInfo     `json:"wechatShare"`
	Labels         []interface{}       `json:"labels"`
	Sponsors       []interface{}       `json:"sponsors"`
	IsCustomized   bool                `json:"isCustomized"`
	IPLoc          string              `json:"ipLoc"`
	TopicID        string              `json:"topicId,omitempty"`
	Transcript     TranscriptInfo      `json:"transcript"`
}

// EpisodeListResponseData defines the "data" field in the episode list API response.
type EpisodeListResponseData struct {
	Data        []Episode   `json:"data"`
	Order       string      `json:"order"`
	Total       int         `json:"total"`
	LoadNextKey LoadMoreKey `json:"loadNextKey,omitempty"`
	LoadMoreKey LoadMoreKey `json:"loadMoreKey,omitempty"` // API has both, aliasing to the same struct
}

// EpisodeDetailAPIResponse defines the overall structure for the get episode detail API response.
type EpisodeDetailAPIResponse struct {
	Data Episode `json:"data"`
}

// EpisodeListAPIResponse has been removed as EpisodeListResponseData is the top-level structure.

// --- Search Related Types ---

// SearchAPILoadMoreKey defines the structure for pagination keys in search results.
type SearchAPILoadMoreKey struct {
	LoadMoreKey interface{} `json:"loadMoreKey"` // Can be int or other types based on API
	SearchID    string      `json:"searchId"`
}

// HighlightWord defines the structure for highlighted words in search results.
type HighlightWord struct {
	Words                  []string `json:"words"`
	SingleMaxHighlightTime int      `json:"singleMaxHighlightTime"`
}

// SearchRequest defines the request body for the search API.
type SearchRequest struct {
	Keyword     string                `json:"keyword"`
	Type        string                `json:"type"` // "PODCAST", "EPISODE", "USER"
	PID         string                `json:"pid,omitempty"`
	LoadMoreKey *SearchAPILoadMoreKey `json:"loadMoreKey,omitempty"`
}

// PodcastSearchResultItem represents a podcast item in search results.
// Reusing AuthorshipEntry as it contains most relevant summary fields.
// Note: The example for podcast search result has a 'truncatedDescription' field
// which is not in AuthorshipEntry. If this is critical, AuthorshipEntry might need
// to be extended or a new specific struct created. For now, reusing.
type PodcastSearchResultItem AuthorshipEntry

// EpisodeSearchResultItem represents an episode item in search results.
// Reusing Episode as it contains all relevant fields.
type EpisodeSearchResultItem Episode

// UserSearchResultItem represents a user item in search results.
type UserSearchResultItem struct {
	Type              string                 `json:"type"`
	UID               string                 `json:"uid"`
	Avatar            Picture                `json:"avatar"`
	Nickname          string                 `json:"nickname"`
	IsNicknameSet     bool                   `json:"isNicknameSet"`
	Bio               string                 `json:"bio,omitempty"`
	Gender            string                 `json:"gender,omitempty"`
	IsCancelled       bool                   `json:"isCancelled"`
	ReadTrackInfo     map[string]interface{} `json:"readTrackInfo,omitempty"` // Example shows this can exist
	IPLoc             string                 `json:"ipLoc,omitempty"`
	Relation          string                 `json:"relation,omitempty"`
	IsBlockedByViewer bool                   `json:"isBlockedByViewer,omitempty"`
}

// PodcastSearchResponse defines the API response structure for podcast search.
type PodcastSearchResponse struct {
	Data          []PodcastSearchResultItem `json:"data"`
	HighlightWord *HighlightWord            `json:"highlightWord,omitempty"`
	LoadMoreKey   *SearchAPILoadMoreKey     `json:"loadMoreKey,omitempty"`
}

// EpisodeSearchResponse defines the API response structure for episode search.
type EpisodeSearchResponse struct {
	Data          []EpisodeSearchResultItem `json:"data"`
	HighlightWord *HighlightWord            `json:"highlightWord,omitempty"`
	LoadMoreKey   *SearchAPILoadMoreKey     `json:"loadMoreKey,omitempty"`
}

// UserSearchResponse defines the API response structure for user search.
type UserSearchResponse struct {
	Data          []UserSearchResultItem `json:"data"`
	HighlightWord *HighlightWord         `json:"highlightWord,omitempty"`
	LoadMoreKey   *SearchAPILoadMoreKey  `json:"loadMoreKey,omitempty"`
}
