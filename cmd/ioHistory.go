package cmd

import (
	"fmt"
	"sync"
)

type Logger struct {
	logs  []string
	mutex sync.Mutex
}

func (l *Logger) AddLogMessage(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	l.logs = append(l.logs, message)
}

func (l *Logger) GetLog() []string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.logs
}

func (c *Cmd) PrintWithLog(message string) {
	c.logger.AddLogMessage(message)
	fmt.Print(message)

}

func (c *Cmd) ShowLog() {
	logs := c.logger.GetLog()
	if len(logs) == 0 {
		fmt.Println("LogIsEmpty")
		return
	}

	for _, e := range logs {
		fmt.Println(e)
	}
}
