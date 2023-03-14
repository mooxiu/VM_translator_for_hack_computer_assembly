package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	randomCount              = 0
	functionReturnPointCount = 0
	originals                = make(map[string][]string, 0)
)

func main() {
	fileName := os.Args[1]
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	/*
		create the file, and get the writer
	*/
	stat, _ := f.Stat()
	var writeFileName string
	if stat.IsDir() {
		writeFileName = f.Name() + "/" + filepath.Base(f.Name()) + ".asm"
	} else {
		writeFileName = filepath.Dir(f.Name()) + "/" + strings.Replace(filepath.Base(f.Name()), ".vm", "", 1) + ".asm"
	}
	toWriteFile, err := os.OpenFile(writeFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	coder := &CodeWriter{File: toWriteFile}
	defer coder.CloseFile()
	/*
		keep each of the file into a map
	*/
	p := &Parser{}
	fillMap(f)
	if stat.IsDir() {
		writeBootstrap(coder, p)
	}
	writeRes(p, coder, originals)
}

func writeBootstrap(coder *CodeWriter, p *Parser) {
	coder.File.WriteString("@256\n")
	coder.File.WriteString("D=A\n")
	coder.File.WriteString("@SP\n")
	coder.File.WriteString("M=D\n")
	coder.AppendLine(coder.Translate(p.Parse("call Sys.init 0"), "bootstrap"))
	coder.File.WriteString("////////////////End of BootStrap")
}

func writeRes(p *Parser, coder *CodeWriter, originals map[string][]string) {
	for fileName, lines := range originals {
		for _, line := range lines {
			coder.AppendComment(line + "\n")
			coder.AppendLine(coder.Translate(p.Parse(line), fileName))
		}
	}
}

func fillMap(f *os.File) {
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	if !stat.IsDir() {
		bfScanner := bufio.NewScanner(f)
		bfScanner.Split(bufio.ScanLines)

		for bfScanner.Scan() {
			line := trimLine(bfScanner.Text())
			if len(line) == 0 {
				continue
			}
			if len(line) > 2 && line[:2] == "//" {
				continue
			}
			originals[stat.Name()] = append(originals[stat.Name()], line)
		}
	} else {
		fileInfos, err := ioutil.ReadDir(f.Name())
		if err != nil {
			panic(err)
		}
		for _, subF := range fileInfos {
			if !subF.IsDir() && !strings.HasSuffix(subF.Name(), ".vm") {
				continue
			}
			subFile, _ := os.Open(f.Name() + "/" + subF.Name())
			defer subFile.Close()
			fillMap(subFile)
		}
	}
}

func trimLine(line string) string {
	idx := strings.Index(line, "//")
	if idx == 0 {
		return ""
	}
	if idx > 0 {
		return trimLine(line[:idx])
	}
	ans := strings.TrimSpace(line)
	return ans
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
	functionCommandKeywords    = []string{"function", "call", "return"}
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
	Action    string // label, goto, if-goto
	LabelName string
}

type FCommand struct {
	Action       string // return, call, function
	FunctionName string
	Count        int
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
	contents := strings.Split(content, " ")
	if len(contents) == 1 {
		return Parsed{
			CommandType: FunctionCommand,
			FCommand: FCommand{
				Action: contents[0],
			},
		}
	} else {
		c, _ := strconv.ParseInt(contents[2], 10, 64)
		return Parsed{
			CommandType: FunctionCommand,
			FCommand: FCommand{
				Action:       contents[0],
				FunctionName: contents[1],
				Count:        int(c),
			},
		}
	}
}

func parseBranchingCommand(content string) Parsed {
	contents := strings.Split(content, " ")
	return Parsed{
		CommandType: BranchingCommand,
		BCommand: BCommand{
			Action:    contents[0],
			LabelName: contents[1],
		},
	}
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

func (c *CodeWriter) Translate(p Parsed, fileName string) []string {
	switch p.CommandType {
	case ArithmeticCommand:
		return translateACommand(p)
	case PushPopCommand:
		return translatePCommand(p, fileName)
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
		ans = append(ans, "A=A-1", "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%d", randomCount), "D;JEQ")
		ans = append(ans, "@SP", "A=M-1", "M=0", fmt.Sprintf("@FALSE_%d", randomCount), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%d)", randomCount), "@SP", "A=M-1", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%d)", randomCount))
	case "gt":
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "A=A-1", "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%d", randomCount), "D;JGT")
		ans = append(ans, "@SP", "A=M-1", "M=0", fmt.Sprintf("@FALSE_%d", randomCount), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%d)", randomCount), "@SP", "A=M-1", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%d)", randomCount))
	case "lt":
		ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
		ans = append(ans, "A=A-1", "D=M-D")
		ans = append(ans, fmt.Sprintf("@TRUE_%d", randomCount), "D;JLT")
		ans = append(ans, "@SP", "A=M-1", "M=0", fmt.Sprintf("@FALSE_%d", randomCount), "0;JMP")
		ans = append(ans, fmt.Sprintf("(TRUE_%d)", randomCount), "@SP", "A=M-1", "M=-1")
		ans = append(ans, fmt.Sprintf("(FALSE_%d)", randomCount))
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
	randomCount += 1
	return ans
}

/*
*
B Command
*/
func translateBCommand(p Parsed) []string {
	ans := make([]string, 0)
	c := p.BCommand
	switch c.Action {
	case "label":
		ans = append(ans, fmt.Sprintf("(%s)", c.LabelName))
	case "goto":
		ans = append(ans, fmt.Sprintf("@%s", c.LabelName), "0;JMP")
	case "if-goto":
		ans = append(ans, "@SP", "AM=M-1", "D=M")
		ans = append(ans, fmt.Sprintf("@%s", c.LabelName), "D;JNE")
	default:
		panic("not support")
	}
	return ans
}

/*
*
F Command
*/
func translateFCommand(p Parsed) []string {
	ans := make([]string, 0)
	c := p.FCommand
	switch c.Action {
	case "return":
		return translateFCommandReturnType(ans, c)
	case "function":
		return translateFCommandFunctionType(ans, c)
	case "call":
		return translateFCommandCallType(ans, c)
	default:
		panic("not support")
	}
	return nil
}

func translateFCommandReturnType(ans []string, c FCommand) []string {
	thisFrame := fmt.Sprintf("@endFrame")
	// endFrame = LCL
	ans = append(ans, "@LCL", "D=M", thisFrame, "M=D")
	// retAddr = *(endFrame - 5)
	ans = append(ans, "@5", "D=D-A", "A=D", "D=M", "@RETURN", "M=D")
	// arg -> pop()
	ans = append(ans, "@SP", "M=M-1", "A=M", "D=M")
	ans = append(ans, "@ARG", "A=M", "M=D")
	// sp = arg + 1
	ans = append(ans, "@ARG", "D=M+1", "@SP", "M=D")
	// that, this, arg, lcl restore
	ans = append(ans, "@1", "D=A", thisFrame, "D=M-D", "A=D", "D=M", "@THAT", "M=D")
	ans = append(ans, "@2", "D=A", thisFrame, "D=M-D", "A=D", "D=M", "@THIS", "M=D")
	ans = append(ans, "@3", "D=A", thisFrame, "D=M-D", "A=D", "D=M", "@ARG", "M=D")
	ans = append(ans, "@4", "D=A", thisFrame, "D=M-D", "A=D", "D=M", "@LCL", "M=D")
	// go to retAddr
	ans = append(ans, "@RETURN", "A=M", "0;JMP")
	return ans
}

func translateFCommandFunctionType(ans []string, c FCommand) []string {
	// inject entry point
	ans = append(ans, fmt.Sprintf("(%s)", c.FunctionName))
	// initialize local segment
	for i := 0; i < c.Count; i++ {
		ans = append(ans, "@SP", "A=M", "M=0")
		ans = append(ans, "@SP", "M=M+1")
	}
	return ans
}

func translateFCommandCallType(ans []string, c FCommand) []string {
	functionReturnPointCount++
	// save the return address, when the execution ends, which line we want to go to
	ans = append(ans, fmt.Sprintf("@returnPoint%d", functionReturnPointCount), "D=A")
	ans = append(ans, "@SP", "A=M", "M=D", "@SP", "M=M+1")
	// save, LCL, ARG, THIS, THAT
	ans = append(ans, "// Call: LCL save")
	ans = append(ans, "@LCL", "D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1")

	ans = append(ans, "// Call: ARG save")
	ans = append(ans, "@ARG", "D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1")

	ans = append(ans, "// Call: THIS save")
	ans = append(ans, "@THIS", "D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1")

	ans = append(ans, "// Call: THAT save")
	ans = append(ans, "@THAT", "D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1")

	// reassign ARG
	ans = append(ans, "@SP", "D=M", "@5", "D=D-A", fmt.Sprintf("@%d", c.Count), "D=D-A")
	ans = append(ans, "@ARG", "M=D")
	// reassign LCL
	ans = append(ans, "@SP", "D=M", "@LCL", "M=D")
	// go to function name
	ans = append(ans, fmt.Sprintf("@%s", c.FunctionName), "0;JMP")
	// insert label to
	ans = append(ans, fmt.Sprintf("(returnPoint%d)", functionReturnPointCount))
	return ans
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

func translatePCommand(p Parsed, fileName string) []string {
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
			ans = append(ans, fmt.Sprintf("@%s.%d", fileName, command.Num), "D=M")
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
			ans = append(ans, fmt.Sprintf("@%s.%d", fileName, command.Num), "M=D")
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
