package utils

func ToCmdLine(cmd ...string) [][]byte {
	args := make([][]byte, 0, len(cmd))
	for _, v := range cmd {
		args = append(args, []byte(v))
	}
	return args
}

// ToCmdLine2 convert commandName and string-type argument to [][]byte
func ToCmdLine2(commandName string, args ...string) [][]byte {
	result := make([][]byte, len(args)+1)
	result[0] = []byte(commandName)
	for i, s := range args {
		result[i+1] = []byte(s)
	}
	return result
}

// ToCmdLine3 convert commandName and []byte-type argument to CmdLine
func ToCmdLine3(commandName string, args ...[]byte) [][]byte {
	result := make([][]byte, len(args)+1)
	result[0] = []byte(commandName)
	for i, s := range args {
		result[i+1] = s
	}
	return result
}
