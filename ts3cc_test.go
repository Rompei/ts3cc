package ts3cc

import (
	"testing"
)

func TestTS3CC(t *testing.T) {
	cl, err := NewTS3CC("localhost:10011", "testuser", "I2PfW1D2", 1)
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()

	server, err := cl.GetServerInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", server)
	showChannel(t, server.Channels, "")
}

func showChannel(t *testing.T, channels Channels, before string) {
	before += "\t"
	for _, channel := range channels {
		t.Logf("%s%s", before, channel.ChannelName)
		for _, client := range channel.Clients {
			t.Logf("%s %s", before, client.Nickname)
		}
		showChannel(t, channel.ChildChannels, before)
	}
}
