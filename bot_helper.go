package MoneyBot

import "teleBot"

type botHelper struct {
	bot MoneyBot
}

func (bh botHelper) GetContainer(msg teleBot.Message) *container {
	container, err := bh.bot.Provider.GetContainer(msg.From.Id)
	if err != nil {
		bh.bot.Chans.LogChan <- "Got " + err.Error() + " while parsing " + msg.String()
		bh.bot.Chans.ErrChan <- ErrChanStruct{err, "while getting container for " + msg.From.String() + ": "}
	}
	return container
}

func (bh botHelper) PushContainer(container *container, msg teleBot.Message) {
	err := bh.bot.Provider.PushContainer(msg.From.Id, container)
	if err != nil {
		bh.bot.Chans.LogChan <- "Got " + err.Error() + " while parsing " + msg.String()
		bh.bot.Chans.ErrChan <- ErrChanStruct{err, "while pushing " + container.String()}
	}
}

func (bh botHelper) GetHistory(msg teleBot.Message) *history {
	history, err := bh.bot.Provider.GetHistory(msg.From.Id)
	if err != nil {
		bh.bot.Chans.LogChan <- "Got " + err.Error() + " while parsing " + msg.String()
		bh.bot.Chans.ErrChan <- ErrChanStruct{err, "while getting history for " + msg.From.String() + ": "}
	}
	return history
}

func (bh botHelper) PushHistory(history *history, msg teleBot.Message) {
	err := bh.bot.Provider.PushHistory(msg.From.Id, history)
	if err != nil {
		bh.bot.Chans.LogChan <- "Got " + err.Error() + " while parsing " + msg.String()
		bh.bot.Chans.ErrChan <- ErrChanStruct{err, "while pushing " + history.String()}
	}
}

func (bh botHelper) ParseErr(err error, msg string) {
	bh.bot.Chans.ParseErr(err, msg)
}
