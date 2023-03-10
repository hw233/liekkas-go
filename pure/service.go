package main

import (
	"bufio"
	"os"
	"strings"
	"unicode/utf8"

	"pure/pb"

	"github.com/huichen/sego"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
)

type server struct {
	replaceWord string
	dirtyWords  map[string]bool
	segmenter   sego.Segmenter
}

func (s *server) init(c *cli.Context) {
	s.replaceWord = c.String("replace-word")
	s.dirtyWords = make(map[string]bool)

	dictionary := c.String("dictionary")
	dirty := c.String("dirty")

	// 载入字典
	log.Info("Loading Dictionary...")
	s.segmenter.LoadDictionary(dictionary)
	log.Info("Dictionary Loaded")

	// 读取脏词库
	log.Info("Loading Dirty Words...")
	f, err := os.Open(dirty)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// 逐行扫描
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		words := strings.Split(strings.ToUpper(strings.TrimSpace(scanner.Text())), " ") // 均处理为大写
		if words[0] != "" {
			s.dirtyWords[words[0]] = true
		}
	}
	log.Info("Dirty Words Loaded")
}

func (s *server) Filter(ctx context.Context, in *pb.FilterRequest) (*pb.FilterReply, error) {
	bin := []byte(in.Text)
	segments := s.segmenter.Segment(bin)
	pureText := make([]byte, 0, len(bin))
	for _, seg := range segments {
		word := bin[seg.Start():seg.End()]
		if s.dirtyWords[strings.ToUpper(string(word))] {
			pureText = append(pureText, []byte(strings.Repeat(s.replaceWord, utf8.RuneCount(word)))...)
		} else {
			pureText = append(pureText, word...)
		}
	}
	return &pb.FilterReply{Text: string(pureText)}, nil
}
