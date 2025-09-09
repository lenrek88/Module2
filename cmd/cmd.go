package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/lenrek88/Module2/calendar"
	"github.com/lenrek88/Module2/logger"
	"github.com/lenrek88/Module2/reminder"
)

type Cmd struct {
	calendar *calendar.Calendar
	logger   *Logger
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar: c,
		logger:   &Logger{},
	}
}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionMaxSuggestion(10),
	)
	go func() {
		for msg := range c.calendar.Notification {
			c.PrintWithLog("Reminder! - " + msg + "\n")
		}
	}()
	p.Run()

}

func (c *Cmd) executor(input string) {
	c.logger.AddLogMessage(input)
	if len(input) < 2 {
		c.PrintWithLog("command cannot be empty  \n")
		return
	}
	parts, err := shlex.Split(input)
	if err != nil {
		c.PrintWithLog("error executing command \n")
		return
	}
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "add":
		c.logger.AddLogMessage("add command called")

		if len(parts) < 4 {
			c.PrintWithLog("Формат: add \"название события\" \"дата и время\" \"приоритет\" \n")
			return
		}

		title := parts[1]
		date := parts[2]
		priority := parts[3]

		e, err := c.calendar.AddEvent(title, date, priority)

		if err != nil {
			c.PrintWithLog("failed to add calendar events: " + err.Error() + "\n")
			logger.Error("failed to add calendar events: " + err.Error() + "\n")
		} else {
			c.PrintWithLog("Событие: " + e.Title + " добавлено \n")
			logger.Info("\"Событие:\", e.Title, \"добавлено\"")
		}

	case "list":
		c.logger.AddLogMessage("list command called")

		events := c.calendar.GetEvents()
		for _, e := range events {
			fmt.Println(e.Title, "-", e.StartAt.Format("2006-01-02 15:04"), e.Priority)
		}

	case "remove":
		c.logger.AddLogMessage("remove command called")

		if len(parts) < 2 {
			c.PrintWithLog("Формат: remove \"название события\" \n")
			return
		}

		e := c.calendar.DeleteEvent(parts[1])

		if e != nil {
			c.PrintWithLog("failed to remove calendar events: " + e.Error() + "\n")
			return
		}

		c.PrintWithLog("Событие удалено \n")
		logger.Info("Событие удалено \n")

	case "update":
		c.logger.AddLogMessage("update command called")
		if len(parts) < 5 {
			c.PrintWithLog("Формат: update \"название события\" \"новое название события\" \"дата и время\" \"приоритет\" \n")
			return
		}
		err := c.calendar.EditEvent(parts[1], parts[2], parts[3], parts[4])
		if err != nil {
			c.PrintWithLog("Update Error:" + err.Error())
			return
		}
		c.PrintWithLog("Событие изменено \n")
		logger.Info("Событие изменено \n")

	case "reminder":
		c.logger.AddLogMessage("reminder command called")

		if len(parts) < 4 {
			c.PrintWithLog("Формат: reminder \"название события\" \"сообщение напоминания\" \"дата и время\" \n")
			return
		}
		err := c.calendar.SetEventReminder(parts[1], parts[1]+" |-| "+parts[2], parts[3])
		if errors.Is(err, reminder.ErrEmptyMessage) {
			c.PrintWithLog("Can't set reminder with empty message \n")
			return
		}
		if err != nil {
			c.PrintWithLog("Ошибка создания напоминания:" + err.Error() + "\n")
			return
		}

		c.PrintWithLog("Напоминание создано \n")
		logger.Info("Напоминание создано \n")

	case "cancel_reminder":
		c.logger.AddLogMessage("cancel_reminder command called")
		if len(parts) < 2 {
			c.PrintWithLog("Формат: reminder \"название события\" \n")
			return
		}
		err := c.calendar.CancelEventReminder(parts[1])
		if err != nil {
			c.PrintWithLog("Ошибка отмены напоминания:" + err.Error() + "\n")
			return
		}
		c.PrintWithLog("Напоминание отменено \n")
		logger.Info("Напоминание отменено \n")

	case "help":
		c.logger.AddLogMessage("help command called")
		fmt.Println("	add - Добавить событие \n" +
			"	list - Показать все события \n" +
			"	remove - Удалить событие \n" +
			"	update - Редактировать событие \n" +
			"	reminder - Добавить напоминание \n" +
			"	cancel_reminder - Добавить напоминание \n" +
			"	help - Показать справку \n" +
			"	logger - Показать весь лог \n" +
			"	history - Показать историю ввода-вывода \n" +
			"	exit - Выход из программы")
	case "logger":
		c.logger.AddLogMessage("logger command called")
		file, err := os.Open("app.log")
		if err != nil {
			log.Println("Ошибка открытия файла:", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(os.Stdout, file)
		if err != nil {
			log.Println("Ошибка чтения файла:", err)
		}
	case "history":
		c.logger.AddLogMessage("ioHistory command called")
		c.ShowLog()
	case "exit":
		c.logger.AddLogMessage("exit command called")
		logger.Info("application is exiting")

		err := c.calendar.Save()
		close(c.calendar.Notification)
		if err != nil {
			log.Println("Ошибка сохранения данных", err)
		}
		os.Exit(0)

	default:
		c.PrintWithLog("Неизвестная команда: \n")
		c.PrintWithLog("Введите 'help' для списка команд \n")
	}

}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие"},
		{Text: "list", Description: "Показать все события"},
		{Text: "remove", Description: "Удалить событие"},
		{Text: "update", Description: "Редактировать событие"},
		{Text: "reminder", Description: "Добавить напоминание"},
		{Text: "cancel_reminder", Description: "Отменить напоминание"},
		{Text: "help", Description: "Показать справку"},
		{Text: "logger", Description: "Показать весь лог"},
		{Text: "history", Description: "Показать историю ввода/вывода"},
		{Text: "exit", Description: "Выйти из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
