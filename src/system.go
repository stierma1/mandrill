package mandrill

type PidSystem interface {
	Register(string, PID)
	Find(string) PID
	Spawn(string, int, int, map[string]interface{}, func(PID, PidSystem) bool) PID
	SpawnDefault(descriptor string, boundFunc func(PID, PidSystem) bool) PID
}

type System struct {
	statConsumer StatConsumer
	registry     map[string]PID
}

var globalSystems = make(map[string]PidSystem)

func DefaultSystem() PidSystem {
	return &System{&StatNullCollector{}, map[string]PID{}}
}

func RegisterGlobalSystem(name string, p PidSystem) {
	globalSystems[name] = p
}

func FindSystem(name string) PidSystem {
	p, _ := globalSystems[name]
	return p
}

func (s *System) Register(name string, pid PID) {
	s.registry[name] = pid
}

func (s *System) Find(name string) PID {
	return s.registry[name]
}

func (s *System) Spawn(descriptor string, mailboxSize int, concurrency int, dictionary map[string]interface{}, boundFunc func(PID, PidSystem) bool) PID {
	pid := Spawn(s, descriptor, mailboxSize, concurrency, dictionary, boundFunc)
	s.statConsumer.Consume(pid.Stats())
	return pid
}

func (s *System) SpawnDefault(descriptor string, boundFunc func(PID, PidSystem) bool) PID {
	pid := SpawnDefault(s, descriptor, boundFunc)
	s.statConsumer.Consume(pid.Stats())
	return pid
}
