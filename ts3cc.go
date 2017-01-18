package ts3cc

import (
	"log"

	"github.com/Darfk/ts3"
)

type TS3CC struct {
	client *ts3.Client
}

func NewTS3CC(addr, id, pass string, server int) (*TS3CC, error) {
	client, err := ts3.NewClient(addr)
	if err != nil {
		return nil, err
	}

	_, err = client.Exec(ts3.Login(id, pass))
	if err != nil {
		return nil, err
	}
	_, err = client.Exec(ts3.Use(server))
	if err != nil {
		panic(err)
	}
	return &TS3CC{
		client: client,
	}, nil
}

func (cc *TS3CC) Close() {
	cc.client.Close()
}

func (cc *TS3CC) GetServerInfo() (*Server, error) {

	response, err := cc.client.Exec(ts3.Command{Command: "serverinfo"})
	if err != nil {
		return nil, err
	}
	server, err := NewServer(response.Params[0])
	if err != nil {
		return nil, err
	}
	channelMap := make(ChannelMap)
	if err := cc.client.WalkChannels(func(i int, ch map[string]string) {
		log.Printf("cid: %s", ch["cid"])
		c := ts3.Command{
			Command: "channelinfo",
			Params: map[string][]string{
				"cid": []string{ch["cid"]},
			},
		}
		response, err := cc.client.Exec(c)
		if err != nil {
			panic(err)
		}
		showResponse(&response)
		log.Println()
		channel, err := NewChannel(ch["cid"], response.Params[0])
		if err != nil {
			panic(err)
		}
		channelMap[channel.CID] = channel
	}); err != nil {
		panic(err)
	}

	if err := cc.client.WalkClients(func(i int, cl map[string]string) {
		log.Printf("clID: %s", cl["clid"])
		log.Println(cl)
		if cl["client_type"] == "0" {
			c := ts3.Command{
				Command: "clientinfo",
				Params: map[string][]string{
					"clid": []string{cl["clid"]},
				},
			}
			response, err := cc.client.Exec(c)
			if err != nil {
				panic(err)
			}
			showResponse(&response)
			log.Println()
			client, err := NewClient(cl["clid"], response.Params[0])
			if err != nil {
				panic(err)
			}
			channelMap[client.CID].Clients = append(channelMap[client.CID].Clients, *client)
		}

	}); err != nil {
		panic(err)
	}
	cc.client.Close()
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
