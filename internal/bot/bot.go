package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiimazurok/ecoflow-status-bot/internal/config"
	"github.com/serhiimazurok/ecoflow-status-bot/pkg/logger"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func New(cfg *config.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.ApiToken)
	if err != nil {
		return nil, err
	}

	tg := &Bot{}
	tg.bot = bot

	tg.bot.Debug = cfg.Telegram.Debug

	logger.Infof("Authorized on account %s", tg.bot.Self.UserName)

	return tg, nil
}

func (tg *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tg.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "devices":
			var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Delta Pro", "Delta Pro"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Delta", "Delta"),
				),
			)
			msg.Text = "List of your devices."
			msg.ReplyMarkup = numericKeyboard
		case "help":
			msg.Text = "I understand /devices."
		default:
			msg.Text = "I don't know that command."
		}

		_, err := tg.bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (tg *Bot) Stop() {
	tg.bot.StopReceivingUpdates()
}
