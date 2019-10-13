package mandrill

import (
	"reflect"
	"testing"
)

func TestDefaultSystem(t *testing.T) {
	tests := []struct {
		name string
		want PidSystem
	}{
		{
			"DefaultSystem  Test",
			DefaultSystem(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultSystem(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_Register(t *testing.T) {
	type fields struct {
		statConsumer StatConsumer
		registry     map[string]PID
	}
	type args struct {
		name string
		pid  PID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"Test Registry",
			fields{
				&StatNullCollector{},
				map[string]PID{},
			},
			args{
				"test",
				DefaultSystem().SpawnDefault("Test", func(p PID, s PidSystem) bool {
					return true
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				statConsumer: tt.fields.statConsumer,
				registry:     tt.fields.registry,
			}
			s.Register(tt.args.name, tt.args.pid)
		})
	}
}

func TestSystem_Find(t *testing.T) {
	type fields struct {
		statConsumer StatConsumer
		registry     map[string]PID
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   PID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				statConsumer: tt.fields.statConsumer,
				registry:     tt.fields.registry,
			}
			if got := s.Find(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("System.Whois() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_Spawn(t *testing.T) {
	type fields struct {
		statConsumer StatConsumer
		registry     map[string]PID
	}
	type args struct {
		descriptor  string
		mailboxSize int
		concurrency int
		dictionary  map[string]interface{}
		boundFunc   func(PID, PidSystem) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   PID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				statConsumer: tt.fields.statConsumer,
				registry:     tt.fields.registry,
			}
			if got := s.Spawn(tt.args.descriptor, tt.args.mailboxSize, tt.args.concurrency, tt.args.dictionary, tt.args.boundFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("System.Spawn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_SpawnDefault(t *testing.T) {
	type fields struct {
		statConsumer StatConsumer
		registry     map[string]PID
	}
	type args struct {
		descriptor string
		boundFunc  func(PID, PidSystem) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   PID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				statConsumer: tt.fields.statConsumer,
				registry:     tt.fields.registry,
			}
			if got := s.SpawnDefault(tt.args.descriptor, tt.args.boundFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("System.SpawnDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
