//go:build wireinject

package structproviders

import "github.com/google/wire"

func InitializeEvent(phrase string, code int) (Event, error) {
	// `*` 表示通配符，都赋值
	//wire.Build(NewEvent, NewGreeter, wire.Struct(new(Message), "*"))
	// 指定字段，只赋值指定字段
	wire.Build(NewEvent, NewGreeter, wire.Struct(new(Message), "Content"))
	return Event{}, nil
}
