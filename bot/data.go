package bot

import (
	"github.com/checkTG/db"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const (
	StateStart     = "start"
	StateChoosePic = "choose"
	StateCheckPic  = "check"
	StateWait      = "wait"
)

type UserState struct {
	ID            int64
	State         string
	LastPhotoID   int
	LastMessageID int
}

var users = make(map[int64]*UserState)
var database *db.Database

var Photos = []string{
	"src/1.jpg",
	"src/2.jpg",
}

var Price = []string{
	"150",
	"300",
}

var Index = 0

var ChooseKeyboard = tu.Keyboard(
	tu.KeyboardRow(tu.KeyboardButton("Выбрать аватарку")),
).WithResizeKeyboard()

var InlineKeyboardConfirm = tg.InlineKeyboardMarkup{
	InlineKeyboard: [][]tg.InlineKeyboardButton{
		{
			{
				Text:         "Да",
				CallbackData: "yes",
			},
			{
				Text:         "Нет",
				CallbackData: "no",
			},
		},
	},
}

var InlineKeyboard = tg.InlineKeyboardMarkup{
	InlineKeyboard: [][]tg.InlineKeyboardButton{
		{
			{
				Text:         "Выбрать",
				CallbackData: Price[Index],
			},
			{
				Text:         "Следующая",
				CallbackData: "next",
			},
		},
	},
}
