package handlers

const (
	headerTemplate = "Search address: %s\nAvailable building addresses and names:"
	lineTemplate   = "%v. %s - %s"
	helpMessage    = `Available commands: 
/addresses - I will return all addresses I know.
/addresses osoi - I will return all addresses that have the prefix "osoi" or "Osoi". You can try any prefix.
/building osoite 1 - I will return information about all buildings with an address of "Osoite 1". You can try any address.
/configure - I will return a menu so that you can manage your preferences.
/help - I will show this message.

If you click the button "Share my location and get the nearest buildings", I will return all known addresses that are close to your location."`
)

var handlersPerCommand = map[string]CommandHandler{
	"start":     {HandlerContainer.start, "Start the bot"},
	"help":      {HandlerContainer.help, "Get help"},
	"settings":  {HandlerContainer.settings, "Configure settings"},
	"addresses": {HandlerContainer.getAllAdresses, "Get all available addresses"},
	"building":  {HandlerContainer.getBuilding, "Get building by address"},
}
