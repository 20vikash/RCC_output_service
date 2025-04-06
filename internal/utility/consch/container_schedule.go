package consch

import (
	"fmt"
	"log"
	"outputservice/internal/env"
	"outputservice/internal/mq"
	"sync"

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

var mu sync.Mutex

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

	for i := range noOfLanguages {
		for j := range noOfRunners {
			runners[fmt.Sprintf("rcc-%s_runner-%s", languages[i], string(rune(j+1)))] = true
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
	log.Println("Done")
	runners[fmt.Sprintf("rcc-python_runner-%s", string(rune(number)))] = true

	mu.Lock()
	pyCount--
	mu.Unlock()

	if pyOccupied {
		pyOccupied = false
		pcChan <- true
	}
}

func PythonSchedule() {
	for range pqChan {
		log.Println(pyCount)
		if pyOccupied {
			log.Println("2 asdasdasda")
			<-pcChan
		}
		for i := range noOfRunners {
			if runners[fmt.Sprintf("rcc-python_runner-%s", string(rune(i+1)))] {
				runners[fmt.Sprintf("rcc-python_runner-%s", string(rune(i+1)))] = false

				mu.Lock()
				pyCount++
				mu.Unlock()

				if pyCount == 5 {
					pyOccupied = true
				}

				latest := pythonQueue.AckCode()
				log.Println(latest.Language)
				go execPython(latest, i+1)

				break
			}
		}
	}
}
