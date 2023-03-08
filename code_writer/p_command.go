package code_writer

import (
	"fmt"
	"vmtranslator/parser"
)

var (
	addrMap = map[string]string{
		"local":    "LCL",
		"argument": "ARG",
		"this":     "THIS",
		"that":     "THAT",
	}
)

func translatePCommand(p parser.Parsed) []string {
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
