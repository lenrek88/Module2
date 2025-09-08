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

	//_, err1 := c.AddEvent("Встреча", "2025-08-01 09:30", "high")
	//if err1 != nil {
	//	fmt.Println("Ошибка", err1)
	//	return
	//}
	//event2, err2 := c.AddEvent("Уборка", "2025-08-01 12:30", "low")
	//if err2 != nil {
	//	fmt.Println("Ошибка", err2)
	//	return
	//}
	//event3, err3 := c.AddEvent("Поликлиника", "2025-08-02 7:30", "medium")
	//if err3 != nil {
	//	fmt.Println("Ошибка", err3)
	//	return
	//}
	//c.ShowEvents()
	//
	//c.DeleteEvent(event2.ID)
	//err = c.EditEvent(event3.ID, "МФЦ", "2025-08-02 10:00", "high")
	//if err != nil {
	//	fmt.Println("Ошибка:", err)
	//}
	//c.ShowEvents()
	//
	//err = c.SetEventReminder(event3.ID, "мфц", "2025-08-02 09:00")
	//if err != nil {
	//	fmt.Println("Ошибка:", err)
	//}
	//c.ShowEvents()
	//
	//defer func() {
	//	err := c.Save()
	//	if err != nil {
	//		fmt.Println("Ошибка сохранения данных", err)
	//	}
	//}()
	//
	//fmt.Scanln()
}
