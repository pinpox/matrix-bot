package matrixbot

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"

	"github.com/matrix-org/gomatrix"
)

//MatrixBot struct to hold the bot and it's methods
type MatrixBot struct {
	//Map a repository to matrix rooms
	Client     *gomatrix.Client
	matrixPass string
	matrixUser string
	Handlers   []CommandHandler
	Name       string
}

//CommandHandler struct to hold a pattern/command asocciated with the
//handling funciton and the needed minimum power of the user in the room
type CommandHandler struct {

	//The pattern or command to handle
	Pattern string

	//The minimal power requeired to execute this command
	MinPower int

	//The function to handle this command
	Handler func(message, room, sender string)

	//Help to be displayed for this command
	Help string
}

func (bot *MatrixBot) getUserPower(room, user string) int {

	powerLevels := struct {
		Users   map[string]int `json:"users"`
		Default int            `json:"users_default"`
	}{}

	if err := bot.Client.StateEvent(room, "m.room.power_levels", "", &powerLevels); err != nil {
		log.Fatal(err)
	}

	//Return the users power or the default user power, if not found
	if power, ok := powerLevels.Users[user]; ok {
		log.Debugf("Found %s found in %v, his power is %v", user, powerLevels.Users, power)
		return power
	}
	log.Debugf("User %s not found in %v", user, powerLevels.Users)
	return powerLevels.Default
}

//RegisterCommand allows to register a command to a handling function
func (bot *MatrixBot) RegisterCommand(pattern string, minpower int, help string, handler func(message string, room string, sender string)) {
	mbch := CommandHandler{
		Pattern:  bot.Name + " " + pattern,
		MinPower: minpower,
		Handler:  handler,
		Help:     help,
	}
	log.Debugf("Registered command: %s [%v]", mbch.Pattern, mbch.MinPower)
	bot.Handlers = append(bot.Handlers, mbch)
}

func (bot *MatrixBot) handleCommands(message, room, sender string) {

	//Don't do anything if the sender is the bot itself
	//TODO edge-case: bot has the same name as a user but on a different server
	if strings.Contains(sender, bot.matrixUser) {
		return
	}

	userPower := bot.getUserPower(room, sender)

	for _, v := range bot.Handlers {
		r, _ := regexp.Compile(v.Pattern)
		if r.MatchString(message) {
			if v.MinPower <= bot.getUserPower(room, sender) {
				v.Handler(message, room, sender)
			} else {
				bot.SendTextToRoom(room, "You have not enough power to execute this command (!"+v.Pattern+").\nYour power: "+strconv.Itoa(userPower)+"\nRequired: "+strconv.Itoa(v.MinPower))
			}
		}
	}
}

//SendHTMLToRoom sends a formattet message to a specified room
func (bot *MatrixBot) SendHTMLToRoom(room, message, messageAlt string) {
	bot.Client.SendMessageEvent(room, "m.room.message", gomatrix.HTMLMessage{
		Body:          messageAlt,
		MsgType:       "m.notice",
		Format:        "org.matrix.custom.html",
		FormattedBody: message,
	})
}

//SendTextToRoom sends a plain-text message to a specified room
func (bot *MatrixBot) SendTextToRoom(room, message string) {
	_, err := bot.Client.SendNotice(room, message)
	if err != nil {
		log.Fatal(err)
	}
}

//Sync syncs the matrix events
func (bot *MatrixBot) Sync() {
	if err := bot.Client.Sync(); err != nil {
		log.Warningf("Sync() returned %s", err)
	}
}

//NewMatrixBot creates a new bot form user credentials
func NewMatrixBot(user, pass string, name string) (*MatrixBot, error) {

	log.Infof("Logging in as %s", user)

	cli, _ := gomatrix.NewClient("http://matrix.org", "", "")

	resp, err := cli.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     user,
		Password: pass,
	})

	if err != nil {
		return nil, err
	}

	cli.SetCredentials(resp.UserID, resp.AccessToken)

	bot := &MatrixBot{
		matrixPass: pass,
		matrixUser: user,
		Client:     cli,
		Name:       name,
	}

	//Setup Syncer and to handle events
	syncer := cli.Syncer.(*gomatrix.DefaultSyncer)

	//Handle messages send to the channel
	syncer.OnEventType("m.room.message", func(ev *gomatrix.Event) {
		log.Debugf("%s said \"%s\" in room %s", ev.Sender, ev.Content["body"], ev.RoomID)
		bot.handleCommands(ev.Content["body"].(string), ev.RoomID, ev.Sender)

	})

	//Handle member events (kick, invite)
	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		log.Debugf("%s invited bot to %s", ev.Sender, ev.RoomID)
		if ev.Content["membership"] == "invite" {
			log.Debugf("Joining Room %s", ev.RoomID)
			if resp, err := cli.JoinRoom(ev.RoomID, "", nil); err != nil {
				log.Fatal(err)
			} else {
				log.Debugf("Joined room %s", resp.RoomID)
			}
		}
	})

	syncer.OnEventType("m.room.power_levels", func(ev *gomatrix.Event) {
		log.Debug("got powerlevel event")
		log.Debug(ev.Body())

	})

	bot.RegisterCommand("help", 0, "Display this help", bot.handleCommandHelp)
	return bot, nil
}

func (bot *MatrixBot) handleCommandHelp(message, room, sender string) {
	//TODO make this a markdown table?

	helpMsg := `The following commands are avaitible for this bot:

Command			Power required		Explanation
----------------------------------------------------------------`

	for _, v := range bot.Handlers {
		helpMsg = helpMsg + "\n!" + v.Pattern + "\t\t\t[" + strconv.Itoa(v.MinPower) + "]\t\t\t\t\t" + v.Help
	}

	bot.SendTextToRoom(room, helpMsg)
}
