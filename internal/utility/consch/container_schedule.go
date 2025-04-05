package consch

import "fmt"

var noOfRunners = 5
var noOfLanguages = 1

var languages = []string{"python"}
var runners = make(map[string]bool)

var pqChan = make(chan bool)
var pcChan = make(chan bool)

var pythonQueue = ConQueue{}

func init() {
	for range noOfLanguages {
		for j := range noOfRunners {
			runners[fmt.Sprintf("rcc-%s_runner-%s", languages[j], string(rune(j)))] = true
		}
	}

	go PythonSchedule()
}

func AddToPythonQueue(language, code, jid string) {
	pythonQueue.AddCode(language, code, jid)
	pqChan <- true
}

func DoneExec() {
	pcChan <- true
}

func PythonSchedule() {

}
