package serve

import (
	getchannelbyname "agent-chat/internal/repository/get-channel-by-name"
	getchannels "agent-chat/internal/repository/get-channels"
	initiatefolder "agent-chat/internal/repository/initiate-folder"
	readmessageby "agent-chat/internal/repository/read-message-by"
	savechannels "agent-chat/internal/repository/save-channels"
	savehistory "agent-chat/internal/repository/save-history"
	sr "agent-chat/internal/repository/send-message"
	addmembertochannel "agent-chat/internal/service/add-member-to-channel"
	"agent-chat/internal/service/install"
	listchannel "agent-chat/internal/service/list-channel"
	readmessage "agent-chat/internal/service/read-message"
	registerchannel "agent-chat/internal/service/register-channel"
	removememberfromchannel "agent-chat/internal/service/remove-member-from-channel"
	sv "agent-chat/internal/service/send-message"
)

type Serve struct {
	Install install.Service

	// channel
	RegisterChannel         registerchannel.Service
	AddMemberToChannel      addmembertochannel.Service
	RemoveMemberFromChannel removememberfromchannel.Service
	ListChannel             listchannel.Service

	// message
	SendMessage sv.Service
	ReadMessage readmessage.Service
}

func New() *Serve {
	return &Serve{
		Install: install.New(initiatefolder.NewAdaptorFile()),

		// channel
		RegisterChannel: registerchannel.New(
			getchannels.NewAdaptorFile(),
			savechannels.NewAdaptorFile(),
		),
		AddMemberToChannel: addmembertochannel.New(
			getchannels.NewAdaptorFile(),
			savechannels.NewAdaptorFile(),
		),
		RemoveMemberFromChannel: removememberfromchannel.New(
			getchannels.NewAdaptorFile(),
			savechannels.NewAdaptorFile(),
		),
		ListChannel: listchannel.New(
			getchannels.NewAdaptorFile(),
		),

		// message
		SendMessage: sv.New(
			getchannelbyname.NewAdaptorFile(),
			savehistory.NewAdaptorFile(),
			sr.NewAdaptorTmux(),
		),
		ReadMessage: readmessage.New(
			getchannelbyname.NewAdaptorFile(),
			readmessageby.NewAdaptorFile(),
		),
	}
}
