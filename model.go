package ts3cc

import (
	"errors"
	"sort"
	"strconv"
)

type Server struct {
	ServerName string   `json:"ServerName"`
	Channels   Channels `json:"channels"`
}

func NewServer(param map[string]string) (*Server, error) {
	serverName, ok := param["virtualserver_name"]
	if !ok {
		return nil, errors.New("virtualserver_name is not found")
	}
	return &Server{
		ServerName: serverName,
	}, nil
}

func (s *Server) SortChannels() {
	s.Channels.Sort()
}

type Channel struct {
	CID           string   `json:"cid"`
	PID           string   `json:"pid"`
	ChannelName   string   `json:"channelName"`
	ChannelOrder  int      `json:"channelOrder"`
	IsLocked      bool     `json:"isLocked"`
	Clients       []Client `json:"clients"`
	ChildChannels Channels `json:"childChannels"`
}

func NewChannel(cid string, param map[string]string) (*Channel, error) {
	pid, ok := param["pid"]
	if !ok {
		return nil, errors.New("pid is not found")
	}
	channelName, ok := param["channel_name"]
	if !ok {
		return nil, errors.New("channel_name is not found")
	}
	channelOrderStr, ok := param["channel_order"]
	if !ok {
		return nil, errors.New("channel_order is not found")
	}
	channelOrder, err := strconv.Atoi(channelOrderStr)
	if err != nil {
		return nil, err
	}
	pass, ok := param["channel_flag_password"]
	if !ok {
		return nil, errors.New("password flag is not found")
	}
	var isLocked bool
	if pass == "1" {
		isLocked = true
	}
	return &Channel{
		CID:          cid,
		PID:          pid,
		ChannelName:  channelName,
		ChannelOrder: channelOrder,
		IsLocked:     isLocked,
	}, nil
}

func (c *Channel) Sort() {
	sort.Sort(c.ChildChannels)
	for _, v := range c.ChildChannels {
		v.Sort()
	}
}

type Channels []*Channel

func (chs Channels) Len() int {
	return len(chs)
}

func (chs Channels) Swap(i, j int) {
	chs[i], chs[j] = chs[j], chs[i]
}

func (chs Channels) Less(i, j int) bool {
	return chs[i].ChannelOrder < chs[j].ChannelOrder
}

func (chs Channels) Sort() {
	sort.Sort(chs)
	for _, v := range chs {
		v.Sort()
	}
}

type ChannelMap map[string]*Channel

func (cm ChannelMap) MakeSlice() Channels {
	for _, v := range cm {
		if v.PID != "0" {
			if parent, ok := cm[v.PID]; ok {
				parent.ChildChannels = append(parent.ChildChannels, v)
			}
		}
	}
	var firstChannels Channels
	for _, v := range cm {
		if v.PID == "0" {
			firstChannels = append(firstChannels, v)
		}
	}
	return firstChannels
}

type Client struct {
	CID            string `json:"cid"`
	ClID           string `json:"clid"`
	Nickname       string `json:"nickname"`
	IsAway         bool   `json:"isAway"`
	IsMicMuted     bool   `json:"isMicMuted"`
	IsSpeakerMuted bool   `json:"isSpeakerMuted"`
}

func NewClient(clid string, param map[string]string) (*Client, error) {
	cid, ok := param["cid"]
	if !ok {
		return nil, errors.New("cid is not found")
	}
	nickname, ok := param["client_nickname"]
	if !ok {
		return nil, errors.New("client_nickname is not found")
	}
	var isMicMuted, isSpeakerMuted, isAway bool
	inputMuted, ok := param["client_input_muted"]
	if !ok {
		return nil, errors.New("client_input_muted is not found")
	}
	if inputMuted == "1" {
		isMicMuted = true
	}
	outputMuted, ok := param["client_output_muted"]
	if !ok {
		return nil, errors.New("client_output_muted is not found")
	}
	if outputMuted == "1" {
		isSpeakerMuted = true
	}
	away, ok := param["client_away"]
	if !ok {
		return nil, errors.New("client_away is not found")
	}
	if away == "1" {
		isAway = true
	}

	return &Client{
		cid,
		clid,
		nickname,
		isAway,
		isMicMuted,
		isSpeakerMuted,
	}, nil
}
