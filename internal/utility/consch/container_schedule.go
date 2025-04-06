package consch

import (
	"fmt"
	"outputservice/internal/env"
	"outputservice/internal/mq"

	amqp "github.com/rabbitmq/amqp091-go"
)

var noOfRunners = 5
var noOfLanguages = 1

var languages = []string{"python"}
var runners = make(map[string]bool)

var pqChan = make(chan bool)
var pcChan = make(chan bool)

var pyOccupied = false
var pyCount = 0

var pythonQueue = ConQueue{}

var (
	con *amqp.Connection
	ch  *amqp.Channel
)

func init() {
	mq := &mq.MQ{
		User: env.GetMqUser(),
		Pass: env.GetMqPassword(),
		Port: "5672",
	}

	con = mq.ConnectToMq()
	ch = mq.CreateChannel(con)

	for range noOfLanguages {
		for j := range noOfRunners {
			runners[fmt.Sprintf("rcc-%s_runner-%s", languages[j], string(rune(j)))] = true
		}
	}

	go PythonSchedule()
}

func AddToPythonQueue(language, code, jid string) {
	pythonQueue.AddCode(language, code, jid)

	go func() {
		pqChan <- true
	}()
}

func PyDoneExec(number int) {
	runners[fmt.Sprintf("rcc-python_runner-%s", string(rune(number)))] = true

	if pyOccupied {
		pcChan <- true
	}
}

func PythonSchedule() {
	for range pqChan {
		if pyOccupied {
			<-pcChan
		}
		for i := range noOfRunners {
			if runners[fmt.Sprintf("rcc-python_runner-%s", string(rune(i)))] {
				runners[fmt.Sprintf("rcc-python_runner-%s", string(rune(i)))] = false
				pyCount++

				if pyCount == 5 {
					pyOccupied = true
				}

				latest := pythonQueue.LatestCode()
				go execPython(latest, i)

				break
			}
		}
	}
}
