package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ElvisChiang/tgpokemongobot/process"
	"github.com/mrd0ll4r/tbotapi"
	"github.com/mrd0ll4r/tbotapi/examples/boilerplate"
)

// DEBUG show verbose debug msg
const DEBUG = true

var pokedexFile = "./data/pokedex.csv"

var gameData []process.GameData

func main() {
	ok := false
	gameData, ok = process.LoadGameData(pokedexFile)
	if !ok {
		fmt.Printf("Game data loading fail\n")
		return
	}
	if DEBUG {
		for _, data := range gameData {
			fmt.Printf("#%d: %s %v %v [%s/%s] %s\n",
				data.Number,
				data.Name,
				data.Evolve,
				data.Nickname,
				data.Type1, data.Type2,
				data.AvatarFile)
		}
	}
	startBot()
	// Command(nil, nil, "/pm")
	// Command(nil, nil, "/pm 1")
	// Command(nil, nil, "/pm 031")
	// Command(nil, nil, "/pm 種子")
	// Command(nil, nil, "/pm23")
	// Command(nil, nil, "123")
	// Command(nil, nil, "456")
	// Command(nil, nil, "種子")
	// Command(nil, nil, "/help")
}

func sendText(api *tbotapi.TelegramBotAPI, chat *tbotapi.Chat, text string) (ok bool) {
	if api == nil || chat == nil {
		fmt.Printf("Error sending text: %s, err = no api or chat\n", text)
		return false
	}
	outMsg, err := api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(*chat), text).Send()
	if err != nil {
		fmt.Printf("Error sending text: %s, err = %s\n", text, err)
		return false
	}
	fmt.Printf("->%d, To:\t%s, %s\n", outMsg.Message.ID, outMsg.Message.Chat, text)
	return true
}

func sendPokemonPic(api *tbotapi.TelegramBotAPI, chat *tbotapi.Chat, pokemon process.GameData) (ok bool) {
	ok = false
	// send a photo
	file, err := os.Open(pokemon.AvatarFile)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		ok = false
		return
	}
	defer file.Close()
	caption := fmt.Sprintf("# %d: %s\n暱稱: %v\n屬性: %s %s\n",
		pokemon.Number, pokemon.Name, pokemon.Nickname, pokemon.Type1, pokemon.Type2)
	for _, evolve := range pokemon.Evolve {
		caption += fmt.Sprintf("進化: /pm%s\n", evolve)
	}
	caption += fmt.Sprintf("")
	fmt.Println(caption)
	if api == nil || chat == nil {
		fmt.Println("Error sending pic, err = no api or chat")
		return false
	}
	photo := api.NewOutgoingPhoto(tbotapi.NewRecipientFromChat(*chat), "pokemon.png", file)
	photo.SetCaption(caption)
	outMsg, err := photo.Send()
	if err != nil {
		fmt.Printf("Error sending photo: %s\n", err)
		return
	}
	fmt.Printf("->%d, To:\t%s, (Photo)\n", outMsg.Message.ID, outMsg.Message.Chat)
	ok = true
	return
}

func sendSticker(api *tbotapi.TelegramBotAPI, chat *tbotapi.Chat, id string) (ok bool) {
	outMsg, err := api.NewOutgoingStickerResend(tbotapi.NewRecipientFromChat(*chat), id).Send()
	if err != nil {
		fmt.Printf("Error sending sticker: %s, err = %s\n", id, err)
		return false
	}
	fmt.Printf("->%d, To:\t%s, sticker %s\n", outMsg.Message.ID, outMsg.Message.Chat, id)
	return true
}

// Command parse tg command line
func Command(api *tbotapi.TelegramBotAPI, chat *tbotapi.Chat, msg string) (ok bool) {

	fmt.Println("- raw msg = " + msg)

	ok = false
	result := strings.Fields(msg)
	if len(result) == 0 {
		return
	}
	command := result[0]
	lowerCmd := strings.ToLower(command)
	lowerCmd = strings.Replace(lowerCmd, "@pgoplusbot", "", -1)
	// fmt.Println("lowerCmd: " + lowerCmd)
	if lowerCmd == "/help" || lowerCmd == "/start" {
		text := "/pm 名字, 暱稱, 編號"
		sendText(api, chat, text)
		ok = true
		return
	}

	msg = strings.TrimPrefix(msg, command)
	msg = strings.TrimSpace(msg)

	if len(msg) == 0 && strings.HasPrefix(lowerCmd, "/pm") {
		msg = strings.Replace(lowerCmd, "/pm", "", -1)
		if len(msg) == 0 {
			return
		}
	} else if len(msg) == 0 {
		msg = lowerCmd
	}
	fmt.Printf("msg = %s, len = %d\n", msg, len(msg))

	pokemon, found := process.FindPokemon(gameData, msg)
	if !found {
		text := "醒醒吧，你沒有" + msg
		sendText(api, chat, text)
		fmt.Println(text)
		return
	}
	ok = true

	sendPokemonPic(api, chat, pokemon)
	return ok
}

func startBot() {
	updateFunc := func(update tbotapi.Update, api *tbotapi.TelegramBotAPI) {
		switch update.Type() {
		case tbotapi.MessageUpdate:
			msg := update.Message
			typ := msg.Type()
			text := "(nil)"
			if typ == tbotapi.StickerMessage {
				sticker := update.Message.Sticker
				fmt.Printf("\tSticker id: %s size: %d\n",
					sticker.FileBase.ID, sticker.FileBase.Size)
			}
			if typ == tbotapi.TextMessage {
				text = *msg.Text
			}
			fmt.Printf("<-%d, From:\t%s, Text: %s \n", msg.ID, msg.Chat, text)
			if typ != tbotapi.TextMessage {
				//ignore non-text messages for now
				fmt.Println("Ignoring non-text message")
				return
			}
			if msg.IsReply() {
				fmt.Println("Ignoring replied message")
				return
			}

			// Note: Bots cannot receive from channels, at least no text messages. So we don't have to distinguish anything here

			// display the incoming message
			// msg.Chat implements fmt.Stringer, so it'll display nicely
			// we know it's a text message, so we can safely use the Message.Text pointer
			// fmt.Printf("<-%d, From:\t%s, Text: %s \n", msg.ID, msg.Chat, *msg.Text)

			Command(api, &msg.Chat, *msg.Text)
		case tbotapi.InlineQueryUpdate:
			fmt.Println("Ignoring received inline query: ", update.InlineQuery.Query)
		case tbotapi.ChosenInlineResultUpdate:
			fmt.Println("Ignoring chosen inline query result (ID): ", update.ChosenInlineResult.ID)
		default:
			fmt.Printf("Ignoring unknown Update type.")
		}
	}

	// run the bot, this will block
	boilerplate.RunBot(apiToken, updateFunc, "pgoplusbot", "Reply Pokemon Go information")
}
