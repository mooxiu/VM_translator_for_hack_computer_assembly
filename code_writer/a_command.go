package code_writer

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"translator/parser"
)

func translateACommand(p parser.Parsed) []string {
	command := p.ACommand
	ans := make([]string, 0)
	switch command.Operation {
	case "add":
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "M=M+D")
		ans = append(ans, "@SP", "M=M+1")
	case "sub":
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "M=M-D")
		ans = append(ans, "@SP", "M=M+1")
	case "neg":
		ans = append(ans, "@SP", "A=M-1")
		ans = append(ans, "M=-M")
	case "eq":
		uuid := uuid2.New().String()
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%s", uuid), "D;JEQ")
		ans = append(ans, "@SP", "A=M", "M=0", fmt.Sprintf("@FALSE_%s", uuid), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%s)", uuid), "@SP", "A=M", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%s)", uuid))
		ans = append(ans, "@SP", "M=M+1")
	case "gt":
		uuid := uuid2.New().String()
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%s", uuid), "D;JGT")
		ans = append(ans, "@SP", "A=M", "M=0", fmt.Sprintf("@FALSE_%s", uuid), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%s)", uuid), "@SP", "A=M", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%s)", uuid))
		ans = append(ans, "@SP", "M=M+1")
	case "lt":
		uuid := uuid2.New().String()
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%s", uuid), "D;JLT")
		ans = append(ans, "@SP", "A=M", "M=0", fmt.Sprintf("@FALSE_%s", uuid), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%s)", uuid), "@SP", "A=M", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%s)", uuid))
		ans = append(ans, "@SP", "M=M+1")
	case "and":
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "M=M&D")
		ans = append(ans, "@SP", "M=M+1")
	case "or":
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "M=M|D")
		ans = append(ans, "@SP", "M=M+1")
	case "not":
		ans = append(ans, "@SP", "A=M-1", "M=!M")
	}
	return ans
}
