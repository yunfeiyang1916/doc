package database

import "strings"

var cmdTable = make(map[string]*command)

type command struct {
	name     string
	executor ExecFunc
	// arity means allowed number of cmdArgs, arity < 0 means len(args) >= -arity.
	// for example: the arity of `get` is 2, `mget` is -2
	arity int
}

func registerCommand(name string, executor ExecFunc, arity int) *command {
	name = strings.ToLower(name)
	cmd := &command{
		name:     name,
		executor: executor,
		arity:    arity,
	}
	cmdTable[name] = cmd
	return cmd
}
