# matrix-bot
BYOB (Build your own bot) - Build a matrix bot that acts on !commands

## Usage
Here is a minimal example on how to build a custom bot that replies to a message "!ping" with "pong".
After starting it, you can invite it to any matrix room and it will join.


```go
package main

import 	"github.com/binaryplease/matrix-bot"

// PingPongBot is a custom bot that will reply to !ping with "pong"
type PingPongBot struct {
	*matrixbot.MatrixBot
}

func main() {

	pass := "supersecretpass"
	user := "myawesomebot"

	bot, err := matrixbot.NewMatrixBot(user, pass)

	if err != nil {
		panic(err)
	}

	mypingPongBot := PingPongBot{bot}
  
        // Register a command like this
	bot.RegisterCommand("!ping", 0, mypingPongBot.handlePing)

	for {
		//Loop forever. If you don't have anything that keeps running, the bot will exit.
	}
}

// Handles the !ping message
func (mybot *PingPongBot) handlePing(message, room, sender string) {
	mybot.SendToRoom(room, "pong!")
}

```
