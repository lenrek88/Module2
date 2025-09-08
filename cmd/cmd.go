package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/lenrek88/app/calendar"
)

var mu sync.Mutex

type Cmd struct {
	calendar *calendar.Calendar
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar: c,
	}
}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
	)
	go func() {
		//mu.Lock()
		//defer mu.Unlock()
		for msg := range c.calendar.Notification {
			fmt.Println("Reminder!", msg)
			file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal("Ошибка открытия файла лога:", err)

			}
			multiWriter := io.MultiWriter(os.Stdout, file)
			log.SetOutput(multiWriter)
			log.Println(msg)
		}
	}()
	p.Run()

}

func (c *Cmd) executor(input string) {
	parts, err := shlex.Split(input)
	if err != nil {
		log.Println("ошибка разделения строки", err)
		return
	}
	cmd := strings.ToLower(parts[0])

	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Ошибка открытия файла лога:", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	switch cmd {
	case "add":
		if len(parts) < 4 {
			log.Println("Формат: add \"название события\" \"дата и время\" \"приоритет\"")
			return
		}

		title := parts[1]
		date := parts[2]
		priority := parts[3]

		e, err := c.calendar.AddEvent(title, date, priority)
		if err != nil {
			log.Println("Ошибка добавления:", err)
		} else {
			log.Println("Событие:", e.Title, "добавлено")
		}

	case "list":
		events := c.calendar.GetEvents()
		for _, e := range events {
			log.Println(e.ID, e.Title, "-", e.StartAt.Format("2006-01-02 15:04"), e.Priority)
		}

	case "remove":
		if len(parts) < 2 {
			log.Println("Формат: remove \"ключ\"")
			return
		}
		c.calendar.DeleteEvent(parts[1])

	case "update":
		if len(parts) < 5 {
			log.Println("Формат: update \"ключ события\" \"название события\" \"дата и время\" \"приоритет\"")
			return
		}
		err := c.calendar.EditEvent(parts[1], parts[2], parts[3], parts[4])
		if err != nil {
			log.Println("Ошибка редактирования:", err)
			return
		}

	case "reminder":
		if len(parts) < 4 {
			log.Println("Формат: reminder \"ключ события\" \"сообщение напоминания\" \"дата и время\"")
			return
		}
		err := c.calendar.SetEventReminder(parts[1], parts[2], parts[3])
		if err != nil {
			log.Println("Ошибка создания напоминания:", err)
		}
	case "help":
		log.Println("add - Добавить событие \n" +
			"list - Показать все события \n" +
			"remove - Удалить событие \n" +
			"update - Редактировать событие \n" +
			"reminder - Добавить напоминание \n" +
			"help - Показать справку \n" +
			"log - Показать весь лог \n" +
			"exit - Выход из программы")
	case "log":
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
	case "exit":
		err := c.calendar.Save()
		close(c.calendar.Notification)
		if err != nil {
			log.Println("Ошибка сохранения данных", err)
		}
		os.Exit(0)

	default:
		log.Println("Неизвестная команда:")
		log.Println("Введите 'help' для списка команд")
	}

}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие"},
		{Text: "list", Description: "Показать все события"},
		{Text: "remove", Description: "Удалить событие"},
		{Text: "update", Description: "Редактировать событие"},
		{Text: "reminder", Description: "Добавить напоминание"},
		{Text: "help", Description: "Показать справку"},
		{Text: "log", Description: "Показать весь лог"},
		{Text: "exit", Description: "Выйти из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
