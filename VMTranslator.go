package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var randomCount = 0

func main() {
	fileName := os.Args[1]
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Name())
	writeFileName := strings.Replace(fileName, ".vm", ".asm", 1)
	toWriteFile, err := os.OpenFile(writeFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		panic(err)
	}
	coder := CodeWriter{File: toWriteFile}
	defer coder.CloseFile()

	bfScanner := bufio.NewScanner(f)
	bfScanner.Split(bufio.ScanLines)
	for bfScanner.Scan() {
		line := strings.TrimSpace(bfScanner.Text())
		if len(line) == 0 {
			continue
		}
		if len(line) > 2 && line[:2] == "//" {
			continue
		}
		p := &Parser{}
		parsed := p.Parse(line)
		coder.AppendComment(line + "\n")
		coder.AppendLine(coder.Translate(parsed))
	}
}

const (
	ArithmeticCommand = "arithmetic"
	PushPopCommand    = "pushpop"
	BranchingCommand  = "branching"
	FunctionCommand   = "function"
)

var (
	arithmeticCommandsKeywords = []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}
	pushPopCommandsKeywords    = []string{"push", "pop"}
	branchingCommandsKeywords  = []string{"label", "goto", "if-goto"}
	functionCommandKeywords    = []string{"Function", "Call", "return"}
)

type Parser struct {
}

type Parsed struct {
	CommandType string
	PCommand    PCommand // push /pop
	ACommand    ACommand
	BCommand    BCommand
	FCommand    FCommand
}

type PCommand struct {
	Action      string // push, pop
	SegmentName string // segment name, local, this, that...
	Num         int
}

type ACommand struct {
	Operation string // add, sub, etc
}

type BCommand struct {
	// TODO
}

type FCommand struct {
	// TODO
}

func (p *Parser) Parse(content string) Parsed {
	content = strings.TrimSpace(content)
	commandType := p.decide(content)
	switch commandType {
	case ArithmeticCommand:
		return parseArithmeticCommand(content)
	case PushPopCommand:
		return parsePushPopCommand(content)
	case BranchingCommand:
		return parseBranchingCommand(content)
	case FunctionCommand:
		return parseFunctionCommand(content)
	default:
		panic("not valid command")
	}
}

func parseFunctionCommand(content string) Parsed {
	// TODO
	return Parsed{}
}

func parseBranchingCommand(content string) Parsed {
	// TODO
	return Parsed{}
}

func parsePushPopCommand(content string) Parsed {
	contents := strings.Split(content, " ")
	num, err := strconv.ParseInt(contents[2], 10, 64)
	if err != nil {
		panic(err)
	}
	return Parsed{
		CommandType: PushPopCommand,
		PCommand: PCommand{
			Action:      contents[0],
			SegmentName: contents[1],
			Num:         int(num),
		},
	}
}

func parseArithmeticCommand(content string) Parsed {
	return Parsed{
		CommandType: ArithmeticCommand,
		ACommand: ACommand{
			Operation: content,
		},
	}
}

func (p *Parser) decide(content string) string {
	contents := strings.Split(content, " ")
	if contains(pushPopCommandsKeywords, contents[0]) {
		return PushPopCommand
	} else if contains(arithmeticCommandsKeywords, contents[0]) {
		return ArithmeticCommand
	} else if contains(branchingCommandsKeywords, contents[0]) {
		return BranchingCommand
	} else if contains(functionCommandKeywords, contents[0]) {
		return FunctionCommand
	}
	panic("not support")
}

func contains(keywords []string, keyword string) bool {
	for _, k := range keywords {
		if k == keyword {
			return true
		}
	}
	return false
}

type CodeWriter struct {
	File *os.File
}

func (c *CodeWriter) Translate(p Parsed) []string {
	switch p.CommandType {
	case ArithmeticCommand:
		return translateACommand(p)
	case PushPopCommand:
		return translatePCommand(p)
	case BranchingCommand:
		return translateBCommand(p)
	case FunctionCommand:
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

/*
*
A Command
*/
func translateACommand(p Parsed) []string {
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
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%d", randomCount), "D;JEQ")
		ans = append(ans, "@SP", "A=M", "M=0", fmt.Sprintf("@FALSE_%d", randomCount), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%d)", randomCount), "@SP", "A=M", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%d)", randomCount))
		ans = append(ans, "@SP", "M=M+1")
		randomCount += 1
	case "gt":
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%d", randomCount), "D;JGT")
		ans = append(ans, "@SP", "A=M", "M=0", fmt.Sprintf("@FALSE_%d", randomCount), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%d)", randomCount), "@SP", "A=M", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%d)", randomCount))
		ans = append(ans, "@SP", "M=M+1")
		randomCount += 1
	case "lt":
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "@SP", "M=M-1", "A=M")
		ans = append(ans, "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%d", randomCount), "D;JLT")
		ans = append(ans, "@SP", "A=M", "M=0", fmt.Sprintf("@FALSE_%d", randomCount), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%d)", randomCount), "@SP", "A=M", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%d)", randomCount))
		ans = append(ans, "@SP", "M=M+1")
		randomCount += 1
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

/*
*
B Command
*/
func translateBCommand(p Parsed) []string {
	// TODO
	return nil
}

/*
*
F Command
*/
func translateFCommand(p Parsed) []string {
	// TODO
	return nil
}

/**
P Command
*/

var (
	addrMap = map[string]string{
		"local":    "LCL",
		"argument": "ARG",
		"this":     "THIS",
		"that":     "THAT",
	}
)

func translatePCommand(p Parsed) []string {
	command := p.PCommand
	ans := make([]string, 0)
	if command.Action == "push" {
		switch command.SegmentName {
		case "local", "argument", "this", "that":
			// addr = local + i
			ans = append(ans, fmt.Sprintf("@%s", addrMap[command.SegmentName]), "D=M")
			ans = append(ans, fmt.Sprintf("@%d", command.Num), "D=D+A")
			// ram[sp] = ram[addr]
			ans = append(ans, "A=D", "D=M")
			ans = append(ans, "@SP", "A=M", "M=D")
			// sp++
			ans = append(ans, "@SP", "M=M+1")
		case "constant":
			// d = const
			ans = append(ans, fmt.Sprintf("@%d", command.Num), "D=A")
			// ram[sp] = d
			ans = append(ans, "@SP", "A=M", "M=D")
			// sp++
			ans = append(ans, "@SP", "M=M+1")
		case "static":
			// D <- @Foo.i
			ans = append(ans, fmt.Sprintf("@Foo.%d", command.Num), "D=M")
			// RAM[SP] <- D
			ans = append(ans, "@SP", "A=M", "M=D")
			// SP++
			ans = append(ans, "@SP", "M=M+1")
		case "pointer":
			if command.Num == 0 {
				// push this
				ans = append(ans, "@THIS")
			} else {
				// push that
				ans = append(ans, "@THAT")
			}
			// D <- this/that
			ans = append(ans, "D=M")
			// RAM[sp] <- D
			ans = append(ans, "@SP", "A=M", "M=D")
			// SP++
			ans = append(ans, "@SP", "M=M+1")
		case "temp":
			// addr = 5 + i
			ans = append(ans, "@R5", "D=A")
			ans = append(ans, fmt.Sprintf("@%d", command.Num), "D=D+A")
			// ram[sp] = ram[addr]
			ans = append(ans, "A=D", "D=M")
			ans = append(ans, "@SP", "A=M", "M=D")
			// sp++
			ans = append(ans, "@SP", "M=M+1")
		}
	} else {
		// "pop"
		switch command.SegmentName {
		case "local", "argument", "this", "that":
			// addr <- LCL + i, store the addr to R13
			ans = append(ans, fmt.Sprintf("@%s", addrMap[command.SegmentName]), "D=M")
			ans = append(ans, fmt.Sprintf("@%d", command.Num), "D=D+A")
			ans = append(ans, "@R13", "M=D")
			// SP--
			ans = append(ans, "@SP", "M=M-1")
			// RAM[addr] <- RAM[SP]
			ans = append(ans, "A=M", "D=M")
			ans = append(ans, "@R13", "A=M", "M=D")
		// no such thing as case "constant":
		case "static":
			// SP--
			ans = append(ans, "@SP", "M=M-1")
			// RAM[addr] <- RAM[SP]
			ans = append(ans, "@SP", "A=M", "D=M")
			ans = append(ans, fmt.Sprintf("@Foo.%d", command.Num), "M=D")
		case "pointer":
			// sp--
			ans = append(ans, "@SP", "M=M-1")
			// D <- RAM[SP]
			ans = append(ans, "@SP", "A=M", "D=M")
			if command.Num == 0 {
				ans = append(ans, "@THIS")
			} else {
				ans = append(ans, "@THAT")
			}
			ans = append(ans, "M=D")
		case "temp":
			// addr <- 5 + i, store the addr to R14
			ans = append(ans, "@R5", "D=A")
			ans = append(ans, fmt.Sprintf("@%d", command.Num), "D=D+A")
			ans = append(ans, "@R14", "M=D")
			// SP--
			ans = append(ans, "@SP", "M=M-1")
			// RAM[addr] <- RAM[SP]
			ans = append(ans, "A=M", "D=M")
			ans = append(ans, "@R14", "A=M", "M=D")
		}
	}
	return ans
}
