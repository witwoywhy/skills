package infrastructure

import "os"

var (
	Path = PathInfo{
		Root: "",

		Channel:             "channel",
		ChannelsJson:        "channels.json",
		ChannelDir:          "",
		ChannelFileJsonPath: "",

		History:    "history",
		HistoryDir: "",
	}
)

type PathInfo struct {
	Root                string
	Channel             string
	ChannelsJson        string
	ChannelDir          string
	ChannelFileJsonPath string
	History             string
	HistoryDir          string
}

func InitPath() {
	root, _ := os.UserHomeDir()
	Path.Root = root + "/.config/agent-chat"

	Path.ChannelDir = Path.Root + "/" + Path.Channel
	Path.ChannelFileJsonPath = Path.ChannelDir + "/" + Path.ChannelsJson

	Path.HistoryDir = Path.Root + "/" + Path.History
}
