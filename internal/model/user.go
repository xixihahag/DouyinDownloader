package model

type UserPost struct {
	StatusCode int64       `json:"status_code"`
	AwemeList  []AwemeList `json:"aweme_list"`
	MaxCursor  int64       `json:"max_cursor"`
	MinCursor  int64       `json:"min_cursor"`
	HasMore    bool        `json:"has_more"`
	Extra      Extra       `json:"extra"`
}

type AwemeList struct {
	LongVideo    interface{} `json:"long_video"`
	Images       interface{} `json:"images"`
	Desc         string      `json:"desc"`
	Author       Author      `json:"author"`
	AwemeType    int64       `json:"aweme_type"`
	ImageInfos   []ImageInfo `json:"image_infos"`
	LabelTopText interface{} `json:"label_top_text"`
	Promotions   interface{} `json:"promotions"`
	AwemeID      string      `json:"aweme_id"`
	ChaList      interface{} `json:"cha_list"`
	VideoLabels  interface{} `json:"video_labels"`
	CommentList  interface{} `json:"comment_list"`
	VideoText    interface{} `json:"video_text"`
	Video        VideoPost   `json:"video"`
	Statistics   Statistics  `json:"statistics"`
	TextExtra    []TextExtra `json:"text_extra"`
	Geofencing   interface{} `json:"geofencing"`
}

type Author struct {
	MixInfo                interface{}   `json:"mix_info"`
	Uid                    string        `json:"uid"`
	Signature              Signature     `json:"signature"`
	FollowingCount         int64         `json:"following_count"`
	PlatformSyncInfo       interface{}   `json:"platform_sync_info"`
	IsAdFake               bool          `json:"is_ad_fake"`
	SECUid                 SECUid        `json:"sec_uid"`
	CardEntries            interface{}   `json:"card_entries"`
	TotalFavorited         string        `json:"total_favorited"`
	Region                 Region        `json:"region"`
	VideoIcon              AvatarLarger  `json:"video_icon"`
	WithCommerceEntry      bool          `json:"with_commerce_entry"`
	Nickname               Nickname      `json:"nickname"`
	FollowStatus           int64         `json:"follow_status"`
	FollowerCount          int64         `json:"follower_count"`
	CustomVerify           string        `json:"custom_verify"`
	Rate                   int64         `json:"rate"`
	UserBadRate            int64         `json:"user_bad_rate"`
	AwemeCount             int64         `json:"aweme_count"`
	VerificationType       int64         `json:"verification_type"`
	EnterpriseVerifyReason string        `json:"enterprise_verify_reason"`
	PolicyVersion          interface{}   `json:"policy_version"`
	WithFusionShopEntry    bool          `json:"with_fusion_shop_entry"`
	UserCanceled           bool          `json:"user_canceled"`
	AvatarThumb            AvatarLarger  `json:"avatar_thumb"`
	AvatarMedium           AvatarLarger  `json:"avatar_medium"`
	FavoritingCount        int64         `json:"favoriting_count"`
	FollowersDetail        interface{}   `json:"followers_detail"`
	StoryOpen              bool          `json:"story_open"`
	Secret                 int64         `json:"secret"`
	HasOrders              bool          `json:"has_orders"`
	Geofencing             interface{}   `json:"geofencing"`
	IsGovMediaVip          bool          `json:"is_gov_media_vip"`
	TypeLabel              []interface{} `json:"type_label"`
	ShortID                string        `json:"short_id"`
	AvatarLarger           AvatarLarger  `json:"avatar_larger"`
	UniqueID               string        `json:"unique_id"`
	WithShopEntry          bool          `json:"with_shop_entry"`
}

type AvatarLarger struct {
	URLList []string `json:"url_list"`
	URI     string   `json:"uri"`
}

type ImageInfo struct {
}

type Nickname string

type Region string

type SECUid string

type Signature string

type Ratio string

type VideoPost struct {
	Ratio         Ratio         `json:"ratio"`
	DownloadAddr  *AvatarLarger `json:"download_addr,omitempty"`
	HasWatermark  bool          `json:"has_watermark"`
	PlayAddrLowbr *AvatarLarger `json:"play_addr_lowbr,omitempty"`
	Height        int64         `json:"height"`
	Width         int64         `json:"width"`
	DynamicCover  *AvatarLarger `json:"dynamic_cover,omitempty"`
	OriginCover   AvatarLarger  `json:"origin_cover"`
	BitRate       interface{}   `json:"bit_rate"`
	Duration      int64         `json:"duration"`
	PlayAddr      AvatarLarger  `json:"play_addr"`
	Cover         AvatarLarger  `json:"cover"`
	Vid           string        `json:"vid"`
}

type TextExtra struct {
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Type        int64  `json:"type"`
	HashtagName string `json:"hashtag_name"`
	HashtagID   int64  `json:"hashtag_id"`
}

type Extra struct {
	Now   int64  `json:"now"`
	Logid string `json:"logid"`
}

type Statistics struct {
	PlayCount    int64  `json:"play_count"`
	ShareCount   int64  `json:"share_count"`
	ForwardCount int64  `json:"forward_count"`
	AwemeID      string `json:"aweme_id"`
	CommentCount int64  `json:"comment_count"`
	DiggCount    int64  `json:"digg_count"`
}
