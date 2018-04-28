package MoneyBot

import (
	"github.com/shvimas/teleBot"
	"strconv"
	"strings"
	"errors"
)

type MoneyBot struct {
	Handler  teleBot.RequestHandler
	Provider Provider
	Chans    LogChans
}

func defaultBot() MoneyBot {
	provider := FSProvider{dataPath} // FIXME: use DB
	handler := teleBot.NewProxyRequestHandler(privateToken, proxyAddr, username, password)
	chans := *NewLogChans()
	return MoneyBot{Handler: handler, Provider: provider, Chans: chans}
}

var DefaultBot = defaultBot()
var helper = botHelper{DefaultBot}

const (
	addUsage     = "Usage: add <category> <amount>"
	removeUsage  = "Usage: remove <category>"
	lookUsage    = "Usage: look"
	historyUsage = "Usage: history"
	setUsage     = "Usage: set <container name>"
	resetUsage   = "Usage: reset [name to save]"
	eraseUsage   = "Usage: erase"
	deleteUsage  = "Usage: delete [name1[, name2, ...]]"
)

func (bot MoneyBot) help(words []string) string {
	if len(words) > 1 {
		if len(words) > 2 {
			return "Usage: help [command]"
		}
		switch strings.ToLower(words[1]) {
		case "add":
			return addUsage
		case "remove":
			return removeUsage
		case "look", "l":
			return lookUsage
		case "history", "h":
			return historyUsage
		case "set":
			return setUsage
		case "reset":
			return resetUsage
		case "erase", "empty":
			return eraseUsage
		case "help":
			return bot.help([]string{})
		default:
			return "Unknown command: " + strings.ToLower(words[1])
		}
	}
	// telegram was asked to let people use monospace fonts since 2014. Morons...
	return `Here are my commands:
		add           add new spending
		remove     remove category
		look          get current container. Alias: l
		history      history of containers. Alias: h
		set            choose current container(from ones in history)
		reset         reset current container(also saves current to history)
		erase        erase current container(does not save it to history)
		empty       alias for erase
		delete       delete specified or all containers from history
		help          get this help

	All commands are case-insensitive.
	You can type help <cmd> to get usage of <cmd>`
}

func (bot MoneyBot) add(words []string, msg teleBot.Message) string {
	if len(words) != 3 {
		return addUsage
	}
	container := helper.GetContainer(msg)
	defer helper.PushContainer(container, msg)
	category := words[1]
	amount, err := strconv.ParseFloat(words[2], 64)
	helper.ParseErr(err, "while parsing "+words[2]+" to float")
	err = container.Add(category, amount)
	helper.ParseErr(err, "while adding "+words[2]+" to "+category+"in"+container.String())
	return "Added " + words[2] + " to " + category
}

func (bot MoneyBot) remove(words []string, msg teleBot.Message) string {
	if len(words) != 2 {
		return removeUsage
	}
	container := helper.GetContainer(msg)
	defer helper.PushContainer(container, msg)
	category := words[1]
	err := container.Delete(category)
	helper.ParseErr(err, "while deleting"+category+"from "+container.String()+": ")
	return "Removed " + category
}

func (bot MoneyBot) look(words []string, msg teleBot.Message) string {
	if len(words) != 1 {
		return lookUsage
	}
	container := helper.GetContainer(msg)
	return container.ToString()
}

func (bot MoneyBot) history(words []string, msg teleBot.Message) string {
	if len(words) != 1 {
		return historyUsage
	}
	history := helper.GetHistory(msg)
	return history.ToString()
}

func (bot MoneyBot) set(words []string, msg teleBot.Message) string {
	if len(words) != 2 {
		return setUsage
	}
	container := helper.GetContainer(msg)
	defer helper.PushContainer(container, msg)
	if container.IsFilled() {
		return "Can't replace a non-empty container.\n" +
			"Please empty the current container with reset(saves current)/erase or empty(do not save current)"
	}
	history := helper.GetHistory(msg)
	defer helper.PushHistory(history, msg)
	name := words[1]
	newContainer, ok := history.Get(name)
	if !ok {
		return "Do not have a container called " + name + " in history"
	}
	container = &newContainer
	return "Set " + name + " as current container"
}

func (bot MoneyBot) reset(words []string, msg teleBot.Message) string {
	if len(words) < 1 || len(words) > 2 {
		return resetUsage
	}
	name := ""
	if len(words) == 2 {
		name = words[1]
	}
	container := helper.GetContainer(msg)
	defer helper.PushContainer(container, msg)
	history := helper.GetHistory(msg)
	defer helper.PushHistory(history, msg)
	name = history.Add(*container, name)
	err := container.Erase()
	helper.ParseErr(err, "while erasing container "+container.String())
	return "Removed all from current container (saved in history as " + name + ")"
}

func (bot MoneyBot) erase(words []string, msg teleBot.Message) string {
	if len(words) != 1 {
		return eraseUsage
	}
	container := helper.GetContainer(msg)
	defer helper.PushContainer(container, msg)
	err := container.Erase()
	helper.ParseErr(err, "while erasing container "+container.String())
	return "Erased current container"
}

func (bot MoneyBot) delete(words []string, msg teleBot.Message) string {
	if len(words) < 2 {
		return deleteUsage
	}
	history := helper.GetHistory(msg)
	defer helper.PushHistory(history, msg)
	for _, name := range words[1:] {
		ok := history.Delete(name)
		if !ok {
			return "Failed to delete " + name + " from history"
		}
	}
	return "Deleted"
}

func (bot MoneyBot) myId(words []string, msg teleBot.Message) string {
	return strconv.Itoa(msg.From.Id)
}

func (bot MoneyBot) parseMessage(msg teleBot.Message) string {
	bot.Chans.LogChan <- "Parsing " + msg.String()
	err := bot.Provider.RegisterUser(msg.From)
	helper.ParseErr(err, "while registering user "+msg.From.String())
	words := strings.Fields(msg.Text)
	if len(words) > 0 {
		switch strings.ToLower(words[0]) {
		case "help":
			return bot.help(words)
		case "add":
			return bot.add(words, msg)
		case "look", "l":
			return bot.look(words, msg)
		case "history", "h":
			return bot.history(words, msg)
		case "remove":
			return bot.remove(words, msg)
		case "reset":
			return bot.reset(words, msg)
		case "set":
			return bot.set(words, msg)
		case "erase", "empty":
			return bot.erase(words, msg)
		case "delete":
			return bot.delete(words, msg)
		case "myid":
			return bot.myId(words, msg)
		default:
			return "Can't understand your request. Type 'help' to get list of available commands"
		}
	} else {
		bot.Chans.LogChan <- "Got empty message somehow: " + msg.String()
		helper.ParseErr(errors.New("Got empty message somehow: "+msg.String()), "")
		return "Failed to parse request: " + msg.Text
	}
}

func (bot MoneyBot) parseUpdates(updates []teleBot.Update) {
	// need to parse updates in the coming order; no goroutines here
	for _, upd := range updates {
		msg := upd.Message
		text := bot.parseMessage(msg)
		params := map[string][]string{
			"chat_id": {strconv.Itoa(msg.Chat.Id)},
			"text":    {text},
		}
		response := teleBot.ResponseUpdate{}
		err := bot.Handler.SendMessage(params, &response)
		if err != nil {
			helper.ParseErr(err, "Error while parsing "+upd.String()+": ")
		} else {
			bot.Chans.LogChan <- "Parsed " + msg.String()
			bot.Chans.LogChan <- "Telegram response: " + response.String()
		}
	}
}

func (bot MoneyBot) requestUpdates(resp *teleBot.GetUpdatesResponse, offset int) {
	params := map[string][]string{
		"offset":  {strconv.Itoa(offset)},
		"timeout": {"10"},
	}
	err := bot.Handler.GetUpdates(params, resp)
	if err != nil {
		helper.ParseErr(err, "while requesting updates: ")
		return
	}
	if len(resp.Res) > 0 {
		bot.Chans.LogChan <- "Got updates: " + resp.String()
	}
}

func (bot MoneyBot) Run() {
	bot.Provider.Connect()
	defer bot.Provider.Close()
	bot.Chans.Start()
	defer bot.Chans.Close()

	resp := teleBot.GetUpdatesResponse{}
	offset := 0
	for ; ; {
		bot.requestUpdates(&resp, offset)
		results := resp.Res
		if len(results) != 0 {
			offset = teleBot.MaxId(results) + 1
			for _, updates := range teleBot.GroupByChatId(&resp) {
				bot.parseUpdates(updates)
			}
		}
	}
}
