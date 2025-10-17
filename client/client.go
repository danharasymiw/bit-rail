package client

type Client interface {
	SendCommand(string, []string)
}
