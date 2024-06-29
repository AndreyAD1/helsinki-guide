package handlers

const (
	headerTemplate = "Search address: %s\nAvailable building addresses and names:"
	buttonTemplate = "%s - %s"
	helpMessage    = `If you send me a message, I will provide all addresses I know that are similar to your message.
If you click the button "Share my location and get the nearest buildings", I will provide all known addresses that are close to your location.
I am aware of buildings located in these Helsinki neighbourhoods: Munkkiniemi, Munkkivuori, Laajasalo, Lauttasaari, and Pohjois-Haaga.

Available commands:
/start - I will send a greeting message.
/addresses - I will return all addresses I know.
/settings - I will return a menu so that you can manage your preferences.
/help - I will show this message.`
	BUILDING_BUTTON    = "building"
	NEXT_BUTTON        = "next"
	LANGUAGE_BUTTON    = "language"
	MAX_MESSAGE_LENGTH = 50
)

var handlersPerCommand = map[string]CommandHandler{
	"start":     {HandlerContainer.start, "Start the bot"},
	"help":      {HandlerContainer.help, "Get help"},
	"settings":  {HandlerContainer.settings, "Configure settings"},
	"addresses": {HandlerContainer.getAllAdresses, "Get all available addresses"},
}
var languageCodes = map[string]string{
	"fi": "Finnish",
	"en": "English",
	"ru": "Russian",
}
