package mandrill

type Stat interface {
	Operation() string
	Tags() []string
	Value() *string
}

type Statistic struct {
	O string
	T []string
	S *string
}

func (stat *Statistic) Operation() string {
	return stat.O
}

func (stat *Statistic) Tags() []string {
	return stat.T
}

func (stat *Statistic) Value() *string {
	return stat.S
}

type StatConsumer interface {
	Consume(ch chan Stat)
}

type StatNullCollector struct {
}

func (col *StatNullCollector) Consume(ch chan Stat) {
	go func() {
		for len(ch) > 0 {
			<-ch
		}
	}()
}
