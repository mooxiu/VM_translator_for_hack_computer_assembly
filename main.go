package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"vmtranslator/code_writer"
	"vmtranslator/parser"
)

func main() {
	fileName := os.Args[1]
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Name())
	writeFileName := strings.ReplaceAll(fileName, ".vm", ".asm")
	toWriteFile, err := os.OpenFile(writeFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		panic(err)
	}
	coder := code_writer.CodeWriter{File: toWriteFile}
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
		p := &parser.Parser{}
		parsed := p.Parse(line)
		coder.AppendComment(line + "\n")
		coder.AppendLine(coder.Translate(parsed))
	}
}
