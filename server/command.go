package server

import (
	"hedis/codec"
	"hedis/core"
)

type CommandContext struct {
	session *Session
	name    *core.String
	args    []*core.String
	command Command
}

type Command func(s *Session, args []*core.String) codec.Message

func CommandPing(s *Session, args []*core.String) codec.Message {
	var msg codec.Message
	if len(args) == 1 {
		msg = codec.NewSimpleStr(args[0])
	} else {
		msg = codec.NewSimpleString("PONG")
	}

	return msg
}

func CommandNotFound(s *Session, args []*core.String) codec.Message {
	msg := codec.NewErrorString("Command not supported")

	return msg
}

func CommandParseFailed(s *Session, args []*core.String) codec.Message {
	msg := codec.NewErrorString("Command parse failed")

	return msg
}

type AllCommands struct {
	cmMap *core.Hash
}

func (t *AllCommands) Get(name *core.String) Command {
	find, i := t.cmMap.Get(name)

	if !find {
		return CommandNotFound
	}

	cmd, ok := i.(Command)
	if !ok {
		return CommandParseFailed
	}

	return cmd
}

func (t *AllCommands) add(name string, cmd Command) {
	t.cmMap.Put(core.NewStringStr(name), cmd)
}

func newAllCommands() *AllCommands {
	commands := &AllCommands{}
	commands.cmMap = core.NewHashSize(1024)
	commands.add("ping", CommandPing)

	return commands
}

var commands = newAllCommands()
