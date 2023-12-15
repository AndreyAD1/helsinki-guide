package handlers

const (
	headerTemplate = "Search address: %s\nAvailable building addresses and names:"
	lineTemplate   = "%v. %s - %s"
)

var handlersPerCommand = map[string]CommandHandler{
	"start":     {HandlerContainer.start, "Start the bot"},
	"help":      {HandlerContainer.help, "Get help"},
	"settings":  {HandlerContainer.settings, "Configure settings"},
	"addresses": {HandlerContainer.getAllAdresses, "Get all available addresses"},
	"building":  {HandlerContainer.getBuilding, "Get building by address"},
}
