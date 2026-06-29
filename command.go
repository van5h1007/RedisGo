package main

import "fmt"

type Command any

type SetCommand struct {
	key, val string
}

type GetCommand struct {
	key string
}

type PingCommand struct{}

type DelCommand struct {
	key string
}

type ExistsCommand struct {
	key string
}

type HelloCommand struct {}

type CommandCommand struct {}

type ClientCommand struct {
	value string
}

func parseCommand(v Value) (Command, error) {
	if v.typ != '*' || len(v.array) == 0 {
		return nil, fmt.Errorf("expected a non empty array but got: %+v", v)
	}

	name := v.array[0].bulk

	switch name {
	case "SET", "set":
		if len(v.array) < 3 {
			return nil, fmt.Errorf("SET requires 2 arguements, got %d", len(v.array)-1)
		}

		return SetCommand{
			key: v.array[1].bulk,
			val: v.array[2].bulk,
		}, nil

	case "GET", "get":
		if len(v.array) < 2 {
			return nil, fmt.Errorf("SET requires 1 arguement1, got %d", len(v.array)-1)
		}

		return GetCommand{
			key: v.array[1].bulk,
		}, nil

	case "PING", "ping":
		return PingCommand{}, nil

	case "DELETE", "delete":
		if len(v.array) < 2 {
			return nil, fmt.Errorf("DEL requires 1 arg, got %d", len(v.array)-1)
		}
		return DelCommand{key: v.array[1].bulk}, nil

	case "EXISTS", "exists":
		if len(v.array) < 2 {
			return nil, fmt.Errorf("EXISTS requires 1 arg, got %d", len(v.array)-1)
		}
		return ExistsCommand{key: v.array[1].bulk}, nil
	
	case "HELLO", "hello":
		return HelloCommand{}, nil

	case "COMMAND", "command":
		return CommandCommand{}, nil

	case "CLIENT", "client":
		val := ""
		if len(v.array) >1 {
			val= v.array[1].bulk
		}
		return ClientCommand{value: val}, nil

	default:
		return nil, fmt.Errorf("unknown command: %q", name)
	}
}
