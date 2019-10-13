package mandrill

import (
	"errors"
	"reflect"
	"sync"
)

type PID interface {
	Kill()
	Send(message []interface{}) error
	Send1(m interface{}) error
	Send2(m1, m2 interface{}) error
	Send3(m1, m2, m3 interface{}) error
	Read() []interface{}
	Read1(v interface{}) bool
	Read2(v1, v2 interface{}) bool
	Read3(v1, v2, v3 interface{}) bool
	Stats() chan Stat
	ExitChan() chan bool
	PutValue(name string, item interface{})
	GetValue(name string, v interface{}) bool
}

type PIDList []PID

type MandrillPID struct {
	sys             PidSystem
	descriptor      string
	concurrency     int
	mailbox         chan []interface{}
	procStart       chan bool
	boundFunc       func(PID, PidSystem) bool
	exit            bool
	exitChan        chan bool
	stats           chan Stat
	dictionary      map[string]interface{}
	dictionaryMutex sync.Mutex
}

func Spawn(sys PidSystem, descriptor string, mailboxSize int, concurrency int, dictionary map[string]interface{}, boundFunc func(PID, PidSystem) bool) PID {
	return &MandrillPID{sys, descriptor, concurrency, make(chan []interface{}, mailboxSize), make(chan bool, concurrency), boundFunc, false, make(chan bool), make(chan Stat, 100), dictionary, sync.Mutex{}}
}

func SpawnDefault(sys PidSystem, descriptor string, boundFunc func(PID, PidSystem) bool) PID {
	return &MandrillPID{sys, descriptor, 1, make(chan []interface{}, 10000000), make(chan bool, 1), boundFunc, false, make(chan bool), make(chan Stat, 100), map[string]interface{}{}, sync.Mutex{}}
}

func (pid *MandrillPID) GetValue(name string, v interface{}) bool {
	pid.dictionaryMutex.Lock()
	i, ok := pid.dictionary[name]
	pid.dictionaryMutex.Unlock()
	if !ok {
		return false
	}
	val := reflect.ValueOf(v)
	val.Elem().Set(reflect.ValueOf(i))
	return true
}

func (pid *MandrillPID) PutValue(name string, v interface{}) {
	pid.dictionaryMutex.Lock()
	pid.dictionary[name] = v
	pid.dictionaryMutex.Unlock()
}

func (pid *MandrillPID) Send(message []interface{}) error {
	if pid.exit {
		return errors.New("PID has exited")
	}
	pid.mailbox <- message
	go pid.process()
	return nil
}

func (pid *MandrillPID) Send1(m interface{}) error {
	return pid.Send([]interface{}{m})
}

func (pid *MandrillPID) Send2(m1 interface{}, m2 interface{}) error {
	return pid.Send([]interface{}{m1, m2})
}

func (pid *MandrillPID) Send3(m1 interface{}, m2 interface{}, m3 interface{}) error {
	return pid.Send([]interface{}{m1, m2, m3})
}

func (pid *MandrillPID) process() {
	if pid.exit {
		return
	}

	pid.procStart <- true
	if pid.exit {
		return
	}
	pid.stats <- &Statistic{"increment", []string{"processInvoked"}, &pid.descriptor}
	exit := pid.boundFunc(pid, pid.sys)
	<-pid.procStart
	if exit {
		pid.Kill()
	}
}

func (pid *MandrillPID) Read() []interface{} {
	return <-pid.mailbox
}

func (pid *MandrillPID) Read1(v interface{}) bool {
	a := pid.Read()
	if len(a) < 1 {
		return false
	}
	i := a[0]
	if i == nil {
		val := reflect.ValueOf(v)
		val.Elem().Set(reflect.Zero(val.Type()))
	} else {
		val := reflect.ValueOf(v)
		val.Elem().Set(reflect.ValueOf(i))
	}
	return true
}
func (pid *MandrillPID) Read2(v1, v2 interface{}) bool {
	a := pid.Read()
	if len(a) < 2 {
		return false
	}
	i := a[0]
	if i == nil {
		val := reflect.ValueOf(v1)
		val.Elem().Set(reflect.Zero(val.Type()))
	} else {
		val := reflect.ValueOf(v1)
		val.Elem().Set(reflect.ValueOf(i))
	}
	i = a[1]
	if i == nil {
		val := reflect.ValueOf(v2)
		val.Elem().Set(reflect.Zero(val.Type()))
	} else {
		val := reflect.ValueOf(v2)
		val.Elem().Set(reflect.ValueOf(i))
	}
	return true
}
func (pid *MandrillPID) Read3(v1, v2, v3 interface{}) bool {
	a := pid.Read()
	if len(a) < 3 {
		return false
	}
	i := a[0]
	if i == nil {
		val := reflect.ValueOf(v1)
		val.Elem().Set(reflect.Zero(val.Type()))
	} else {
		val := reflect.ValueOf(v1)
		val.Elem().Set(reflect.ValueOf(i))
	}
	i = a[1]
	if i == nil {
		val := reflect.ValueOf(v2)
		val.Elem().Set(reflect.Zero(val.Type()))
	} else {
		val := reflect.ValueOf(v2)
		val.Elem().Set(reflect.ValueOf(i))
	}
	i = a[2]
	if i == nil {
		val := reflect.ValueOf(v3)
		val.Elem().Set(reflect.Zero(val.Type()))
	} else {
		val := reflect.ValueOf(v3)
		val.Elem().Set(reflect.ValueOf(i))
	}
	return true
}

func (pid *MandrillPID) Stats() chan Stat {
	return pid.stats
}

func (pid *MandrillPID) Kill() {
	if pid.exit {
		return
	}
	pid.exit = true
	pid.stats <- &Statistic{"event", []string{"killed"}, &pid.descriptor}
	close(pid.mailbox)
	close(pid.procStart)
	close(pid.stats)
	pid.exitChan <- true
	close(pid.exitChan)
}

func (pid *MandrillPID) ExitChan() chan bool {
	return pid.exitChan
}
