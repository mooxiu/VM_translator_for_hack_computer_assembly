package code_writer

import (
	"os"
	"vmtranslator/parser"
)

type CodeWriter struct {
	File *os.File
}

func (c *CodeWriter) Translate(p parser.Parsed) []string {
	switch p.CommandType {
	case parser.ArithmeticCommand:
		return translateACommand(p)
	case parser.PushPopCommand:
		return translatePCommand(p)
	case parser.BranchingCommand:
		return translateBCommand(p)
	case parser.FunctionCommand:
		return translateFCommand(p)
	default:
		panic("not support")
	}
}

func (c *CodeWriter) AppendLine(translated []string) {
	for _, str := range translated {
		_, err := c.File.WriteString(str)
		if err != nil {
			panic(err)
		}
		_, err = c.File.WriteString("\n")
		if err != nil {
			panic(err)
		}
	}
}

func (c *CodeWriter) CloseFile() {
	err := c.File.Close()
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) AppendComment(line string) {
	_, err := c.File.WriteString("// " + line)
	if err != nil {
		return
	}
}
