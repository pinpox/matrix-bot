package matrixbot

import (
	"fmt"
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

func (gb *MatrixBot) getSenderPower(sender string) int {
	// gomatrix: get "m.room.power_levels" type state event, parse it and find users->userID or use users_default if not found

	// getting the state event: client.StateEvent(room, type, stateKey)
	//TODO
	return 100
}

//RegisterCommand allows to register a command to a handling function
func (bot *MatrixBot) RegisterCommand(pattern string, minpower int, help string, handler func(message string, room string, sender string)) {
	mbch := CommandHandler{
		Pattern:  bot.Name + " " + pattern,
		MinPower: minpower,
		Handler:  handler,
		Help:     help,
	}
	fmt.Println("Registered command: " + mbch.Pattern)
	bot.Handlers = append(bot.Handlers, mbch)
}

func (bot *MatrixBot) handleCommands(message, room, sender string) {

	//Don't do anything if the sender is the bot itself
	//TODO edge-case: bot has the same name as a user but on a different server
	if strings.Contains(sender, bot.matrixUser) {
		return
	}

	for _, v := range bot.Handlers {
		r, _ := regexp.Compile(v.Pattern)
		if r.MatchString(message) {
			if v.MinPower <= bot.getSenderPower(sender) {
				v.Handler(message, room, sender)
			} else {
				bot.SendToRoom(room, "You have not enough power to execute this command ("+v.Pattern+"). Your power: "+string(bot.getSenderPower(sender))+", requeired: "+string(v.MinPower))
			}
		}
	}
}

//SendToRoom sends a message to a specified room
func (gb *MatrixBot) SendToRoom(room, message string) {
	_, err := gb.Client.SendText(room, message)
	if err != nil {
		panic(err)
	}
}

func (bot *MatrixBot) Sync() {

	//Spawn goroutine to keep checking for events
	// go func() {
	// 	for {
	if err := bot.Client.Sync(); err != nil {
		fmt.Println("Sync() returned ", err)
	}
	// Optional: Wait a period of time before trying to sync again.
	// }
	// }()

}

//NewMatrixBot creates a new bot form user credentials
func NewMatrixBot(user, pass string, name string) (*MatrixBot, error) {

	fmt.Println("Logging in")

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
		fmt.Println(ev.Sender + " said: \"" + ev.Content["body"].(string) + "\" in room : " + ev.RoomID)
		bot.handleCommands(ev.Content["body"].(string), ev.RoomID, ev.Sender)

	})

	//Handle member events (kick, invite)
	syncer.OnEventType("m.room.member", func(ev *gomatrix.Event) {
		fmt.Println(ev.Sender + " invited bot to " + ev.RoomID)

		if ev.Content["membership"] == "invite" {

			fmt.Println("Joining Room")

			if resp, err := cli.JoinRoom(ev.RoomID, "", nil); err != nil {
				panic(err)
			} else {
				fmt.Println(resp.RoomID)
			}
		}
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

	bot.SendToRoom(room, helpMsg)
}
