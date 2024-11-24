package bot

import tg "github.com/mymmrac/telego"

var Photos = []string{
	"src/1.jpg",
	"src/2.jpg",
}

var Price = []string{
	"150",
	"300",
}

var Index = 0

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
				CallbackData: string(rune(Index)),
			},
			{
				Text:         "Следующая",
				CallbackData: "next",
			},
		},
	},
}
