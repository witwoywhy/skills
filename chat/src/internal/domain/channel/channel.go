package channel

import channeltype "agent-chat/internal/enum/channel-type"

type Channel struct {
	Name        string           `json:"name"`
	Members     []ChannelMember  `json:"members"`
	Type        channeltype.Type `json:"type"`
	Description string           `json:"description"`
}

type ChannelMember struct {
	Fleet  string `json:"fleet"`
	Member string `json:"member"`
}
