package ts3cc

import (
	"log"

	"github.com/Darfk/ts3"
)

type TS3CC struct {
	addr     string
	id       string
	pass     string
	serverID int
}

func NewTS3CC(addr, id, pass string, server int) (*TS3CC, error) {
	return &TS3CC{
		addr,
		id,
		pass,
		server,
	}, nil
}

func (cc *TS3CC) GetServerInfo() (*Server, error) {
	client, err := ts3.NewClient(cc.addr)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	_, err = client.Exec(ts3.Login(cc.id, cc.pass))
	if err != nil {
		return nil, err
	}
	_, err = client.Exec(ts3.Use(cc.serverID))
	if err != nil {
		panic(err)
	}

	response, err := client.Exec(ts3.Command{Command: "serverinfo"})
	if err != nil {
		return nil, err
	}
	server, err := NewServer(response.Params[0])
	if err != nil {
		return nil, err
	}
	channelMap := make(ChannelMap)
	if err := client.WalkChannels(func(i int, ch map[string]string) {
		c := ts3.Command{
			Command: "channelinfo",
			Params: map[string][]string{
				"cid": []string{ch["cid"]},
			},
		}
		response, err := client.Exec(c)
		if err != nil {
			panic(err)
		}
		channel, err := NewChannel(ch["cid"], response.Params[0])
		if err != nil {
			panic(err)
		}
		channelMap[channel.CID] = channel
	}); err != nil {
		panic(err)
	}

	if err := client.WalkClients(func(i int, cl map[string]string) {
		if cl["client_type"] == "0" {
			c := ts3.Command{
				Command: "clientinfo",
				Params: map[string][]string{
					"clid": []string{cl["clid"]},
				},
			}
			response, err := client.Exec(c)
			if err != nil {
				panic(err)
			}
			client, err := NewClient(cl["clid"], response.Params[0])
			if err != nil {
				panic(err)
			}
			channelMap[client.CID].Clients = append(channelMap[client.CID].Clients, *client)
		}

	}); err != nil {
		panic(err)
	}
	server.Channels = channelMap.MakeSlice()
	server.SortChannels()
	return server, nil
}

func showResponse(r *ts3.Response) {
	for i, e := range r.Params {
		for k, v := range e {
			log.Printf("%d: %s: %s", i, k, v)
		}
	}
}
