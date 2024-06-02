package transports

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"notifications/internal/dto"
)

type TelegramConfig struct {
	Token string `json:"tkn" yaml:"tkn"`
}

type TelegramSender struct {
	cfg *TelegramConfig
	bot *tgbotapi.BotAPI
}

func (ts *TelegramSender) Send(to []string, msg *dto.Message, meta map[string]string) error {
	if ts.bot == nil {
		var initBotErr error
		ts.bot, initBotErr = tgbotapi.NewBotAPI(ts.cfg.Token)
		if initBotErr != nil {
			return initBotErr
		}
	}

	tgMessageFormatted := fmt.Sprintf("<b>%s</b>\n\n%s", msg.Subject, msg.Body)

	for _, chatId := range to {
		i, pErr := strconv.ParseInt(chatId, 10, 64)
		if pErr != nil {
			return pErr
		}
		tgMessage := tgbotapi.NewMessage(i, tgMessageFormatted)
		tgMessage.ParseMode = tgbotapi.ModeHTML
		_, sendErr := ts.bot.Send(tgMessage)
		if sendErr != nil {
			return sendErr
		}
	}

	return nil
}
