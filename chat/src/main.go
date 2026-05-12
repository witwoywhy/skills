package main

import (
	"agent-chat/infrastructure"
	addmembertochannel "agent-chat/internal/service/add-member-to-channel"
	"agent-chat/internal/service/install"
	listchannel "agent-chat/internal/service/list-channel"
	readmessage "agent-chat/internal/service/read-message"
	registerchannel "agent-chat/internal/service/register-channel"
	removememberfromchannel "agent-chat/internal/service/remove-member-from-channel"
	sendmessage "agent-chat/internal/service/send-message"
	"agent-chat/serve"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Install bool `short:"i" long:"install" description:"initiate resource for agent-chat"`

	// channel
	ListChannel             bool                            `long:"list-channel" description:"list channel"`
	RegisterChannel         registerchannel.Request         `group:"register channel"`
	AddMemberToChannel      addmembertochannel.Request      `group:"add member to channel"`
	RemoveMemberFromChannel removememberfromchannel.Request `group:"remove member from channel"`

	// message
	SendMessage sendmessage.Request `group:"send message"`
	ReadMessage readmessage.Request `group:"read message"`
}

func init() {
	infrastructure.InitValidate()
	infrastructure.InitPath()
}

func main() {
	var options Options

	p := flags.NewParser(&options, flags.Default)
	_, err := p.Parse()
	switch err := err.(type) {
	case nil:
		break
	case *flags.Error:
		if err.Type == flags.ErrHelp {
			os.Exit(0)
		}

		panic(err)
	default:
		panic(err)
	}

	services := serve.New()

	if options.Install {
		_, err := services.Install.Execute(&install.Request{})
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	if options.RegisterChannel.Do {
		_, err := services.RegisterChannel.Execute(&options.RegisterChannel)
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	if options.AddMemberToChannel.Do {
		_, err := services.AddMemberToChannel.Execute(&options.AddMemberToChannel)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if options.RemoveMemberFromChannel.Do {
		_, err := services.RemoveMemberFromChannel.Execute(&options.RemoveMemberFromChannel)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if options.ListChannel {
		_, err := services.ListChannel.Execute(&listchannel.Request{})
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if options.SendMessage.Do {
		_, err := services.SendMessage.Execute(&options.SendMessage)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if options.ReadMessage.Do {
		_, err := services.ReadMessage.Execute(&options.ReadMessage)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

}
