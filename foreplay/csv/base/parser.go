package base

// csv reader module

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

type Parser struct {
	// All line array
	lines [][]string

	// Field name to index map if there is head line
	index2Name map[string]int
}

func NewParser() *Parser {
	return &Parser{
		index2Name: make(map[string]int),
	}
}

// boom check, just remove it if exists it.
func checkBoom(row []byte) []byte {
	if len(row) >= 3 && row[0] == 0xEF && row[1] == 0xBB && row[2] == 0xBF {
		return row[3:]
	} else {
		return row
	}
}

// load csv config, fail return error
// head flag denotes the csv content includes row note.
func (p *Parser) Load(fileName string, head bool) error {
	fs, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fs.Close()

	r := csv.NewReader(fs)
	firstLine := true
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if firstLine { // first line bom check.
			row[0] = string(checkBoom([]byte(row[0])))
			firstLine = false
		}
		p.lines = append(p.lines, row)
	}

	// save index to name map if possible.
	if head && p.GetAllCount() >= 1 {
		size, _ := p.GetCount(0)
		for i := 0; i < size; i++ {
			name := p.lines[0][i]
			if _, ok := p.index2Name[name]; ok {
				return errors.New(fmt.Sprintf("Find field name:%s repeated", name))
			} else {
				p.index2Name[name] = i
			}
		}
	}
	return nil
}

// 获取长度信息。
func (p *Parser) GetAllCount() int {
	return len(p.lines)
}

// 获取index的字段数量
func (p *Parser) GetCount(index uint32) (int, bool) {
	if index >= uint32(p.GetAllCount()) {
		return 0, false
	}
	return len(p.lines[index]), true
}

// 获取对应的字符串。
func (p *Parser) GetField(line uint32, field uint32) (string, bool) {
	fieldSize, _ := p.GetCount(line)
	if line >= uint32(p.GetAllCount()) || field >= uint32(fieldSize) {
		return "", false
	} else {
		return p.lines[line][field], true
	}
}

func (p *Parser) GetFieldByName(line uint32, name string) (string, bool) {
	if index, ok := p.index2Name[name]; ok {
		return p.GetField(line, uint32(index))
	} else {
		return "", false
	}
}

func (p *Parser) GetFieldInt64(line uint32, field uint32) (int64, bool) {
	c, ok := p.GetField(line, field)
	if ok {
		return String2Int64(c)
	} else {
		return 0, false
	}
}

func (p *Parser) GetFieldInt64ByName(line uint32, name string) (int64, bool) {
	c, ok := p.GetFieldByName(line, name)
	if ok {
		return String2Int64(c)
	} else {
		return 0, false
	}
}

// 获取uint64字段信息
func (p *Parser) GetFieldUint64(line uint32, field uint32) (uint64, bool) {
	c, ok := p.GetField(line, field)
	if ok {
		return String2Uint64(c)
	} else {
		return 0, false
	}
}

// uint64 by name
func (p *Parser) GetFieldUint64ByName(line uint32, name string) (uint64, bool) {
	c, ok := p.GetFieldByName(line, name)
	if ok {
		return String2Uint64(c)
	} else {
		return 0, false
	}
}

// 获取所有的信息
func (p *Parser) GetAllContent() [][]string {
	return p.lines
}
