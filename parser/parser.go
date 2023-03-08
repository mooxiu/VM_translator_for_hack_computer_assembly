package parser

import (
	"strconv"
	"strings"
)

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
