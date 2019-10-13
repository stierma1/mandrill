package mandrill

import (
	"reflect"
	"testing"
)

func TestEnd2End1(t *testing.T) {
	name := "End 2 end 1"
	t.Run(name, func(t *testing.T) {

		system := DefaultSystem()
		p := system.SpawnDefault("test1", func(pid PID, system PidSystem) bool {
			var message string
			pid.Read1(&message)

			if !reflect.DeepEqual(message, "Hello") {
				t.Errorf("DefaultSystem() = %v, want %v", message, "Hello")
			}

			return true
		})
		p.Send1("Hello")

		<-p.ExitChan()

	})
	name = "End 2 end 2 - multiple messages processed sequentially"
	t.Run(name, func(t *testing.T) {
		var count int
		system := DefaultSystem()
		p := system.SpawnDefault("test1", func(pid PID, system PidSystem) bool {
			var message string
			pid.Read1(&message)

			count++
			if count > 1 {
				if !reflect.DeepEqual(message, "World") {
					t.Errorf("DefaultSystem() = %v, want %v", message, "World")
				}
				return true
			}
			if !reflect.DeepEqual(message, "Hello") {
				t.Errorf("DefaultSystem() = %v, want %v", message, "Hello")
			}
			return false
		})
		p.Send1("Hello")
		p.Send1("World")

		<-p.ExitChan()

	})
	name = "End 2 end 3 - pid passed to pid"
	t.Run(name, func(t *testing.T) {

		system := DefaultSystem()
		p := system.SpawnDefault("test1", func(pid PID, system PidSystem) bool {
			var message string
			var pid2 PID
			pid.Read2(&message, &pid2)

			if message == "Hello" {
				pid2.Send2(message, pid)
				return false
			}

			if !reflect.DeepEqual(message, "Done") {
				t.Errorf("DefaultSystem() = %v, want %v", message, "Done")
			}

			return true
		})
		p2 := system.SpawnDefault("test2", func(pid PID, system PidSystem) bool {
			var message string
			var pid2 PID
			pid.Read2(&message, &pid2)

			if !reflect.DeepEqual(message, "Hello") {
				t.Errorf("DefaultSystem() = %v, want %v", message, "Hello")
			}

			pid2.Send2("Done", pid)
			return true
		})
		p.Send2("Hello", p2)

		<-p.ExitChan()

	})

}
