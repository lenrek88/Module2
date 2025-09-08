package main

import (
	"fmt"

	"github.com/lenrek88/Module2/calendar"
	"github.com/lenrek88/Module2/cmd"
	"github.com/lenrek88/Module2/logger"
	"github.com/lenrek88/Module2/storage"
)

func main() {

	err := logger.Init()
	if err != nil {
		fmt.Println("Ошибка инициализации логгера", err)
	}
	defer logger.CloseFile()
	logger.Info("application is running")

	//s := storage.NewJsonStorage("calendar.json")
	zs := storage.NewZipStorage("calendar.zip")
	c := calendar.NewCalendar(zs)

	err = c.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки данных", err)
	}

	cli := cmd.NewCmd(c)
	cli.Run()

}
