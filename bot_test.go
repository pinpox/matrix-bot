package matrixbot

import (
	"reflect"
	"testing"

	"github.com/matrix-org/gomatrix"
)

func TestMatrixBot_getUserPower(t *testing.T) {
	type fields struct {
		Client     *gomatrix.Client
		matrixPass string
		matrixUser string
		Handlers   []CommandHandler
		Name       string
	}
	type args struct {
		room string
		user string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &MatrixBot{
				Client:     tt.fields.Client,
				matrixPass: tt.fields.matrixPass,
				matrixUser: tt.fields.matrixUser,
				Handlers:   tt.fields.Handlers,
				Name:       tt.fields.Name,
			}
			if got := bot.getUserPower(tt.args.room, tt.args.user); got != tt.want {
				t.Errorf("MatrixBot.getUserPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrixBot_RegisterCommand(t *testing.T) {
	type fields struct {
		Client     *gomatrix.Client
		matrixPass string
		matrixUser string
		Handlers   []CommandHandler
		Name       string
	}
	type args struct {
		pattern  string
		minpower int
		help     string
		handler  func(message string, room string, sender string)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &MatrixBot{
				Client:     tt.fields.Client,
				matrixPass: tt.fields.matrixPass,
				matrixUser: tt.fields.matrixUser,
				Handlers:   tt.fields.Handlers,
				Name:       tt.fields.Name,
			}
			bot.RegisterCommand(tt.args.pattern, tt.args.minpower, tt.args.help, tt.args.handler)
		})
	}
}

func TestMatrixBot_handleCommands(t *testing.T) {
	type fields struct {
		Client     *gomatrix.Client
		matrixPass string
		matrixUser string
		Handlers   []CommandHandler
		Name       string
	}
	type args struct {
		message string
		room    string
		sender  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &MatrixBot{
				Client:     tt.fields.Client,
				matrixPass: tt.fields.matrixPass,
				matrixUser: tt.fields.matrixUser,
				Handlers:   tt.fields.Handlers,
				Name:       tt.fields.Name,
			}
			bot.handleCommands(tt.args.message, tt.args.room, tt.args.sender)
		})
	}
}

func TestMatrixBot_SendToRoom(t *testing.T) {
	type fields struct {
		Client     *gomatrix.Client
		matrixPass string
		matrixUser string
		Handlers   []CommandHandler
		Name       string
	}
	type args struct {
		room    string
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &MatrixBot{
				Client:     tt.fields.Client,
				matrixPass: tt.fields.matrixPass,
				matrixUser: tt.fields.matrixUser,
				Handlers:   tt.fields.Handlers,
				Name:       tt.fields.Name,
			}
			bot.SendTextToRoom(tt.args.room, tt.args.message)
		})
	}
}

func TestMatrixBot_Sync(t *testing.T) {
	type fields struct {
		Client     *gomatrix.Client
		matrixPass string
		matrixUser string
		Handlers   []CommandHandler
		Name       string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &MatrixBot{
				Client:     tt.fields.Client,
				matrixPass: tt.fields.matrixPass,
				matrixUser: tt.fields.matrixUser,
				Handlers:   tt.fields.Handlers,
				Name:       tt.fields.Name,
			}
			bot.Sync()
		})
	}
}

func TestNewMatrixBot(t *testing.T) {
	type args struct {
		user string
		pass string
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *MatrixBot
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMatrixBot(tt.args.user, tt.args.pass, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMatrixBot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMatrixBot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrixBot_handleCommandHelp(t *testing.T) {
	type fields struct {
		Client     *gomatrix.Client
		matrixPass string
		matrixUser string
		Handlers   []CommandHandler
		Name       string
	}
	type args struct {
		message string
		room    string
		sender  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &MatrixBot{
				Client:     tt.fields.Client,
				matrixPass: tt.fields.matrixPass,
				matrixUser: tt.fields.matrixUser,
				Handlers:   tt.fields.Handlers,
				Name:       tt.fields.Name,
			}
			bot.handleCommandHelp(tt.args.message, tt.args.room, tt.args.sender)
		})
	}
}
