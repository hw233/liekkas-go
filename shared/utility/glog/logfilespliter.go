package glog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
)

type FileSizeSpliter struct {
	splitNo     int
	maxFileSize int
	curFileSize int
}

func NewFileSizeSpliter(maxSize int) *FileSizeSpliter {
	return &FileSizeSpliter{
		splitNo:     0,
		maxFileSize: maxSize,
		curFileSize: 0,
	}
}

func (fss *FileSizeSpliter) SetSplitNo(num int) {
	fss.splitNo = num
}

func (fss *FileSizeSpliter) OnWrite(writeBytes int) {
	fss.curFileSize = fss.curFileSize + writeBytes
}

func (fss *FileSizeSpliter) Restart() {
	fss.splitNo = 0
	fss.curFileSize = 0
}

func (fss *FileSizeSpliter) TrySplit() bool {
	if fss.maxFileSize <= 0 {
		return false
	}

	if fss.curFileSize > fss.maxFileSize {
		fss.splitNo = fss.splitNo + 1
		fss.curFileSize = 0
		return true
	}

	return false
}

func (fss *FileSizeSpliter) AppendDirAndName(dirBuf, nameBuf *bytes.Buffer) {
	if fss.splitNo > 0 {
		if nameBuf.Len() > 0 {
			nameBuf.WriteByte('-')
		}
		nameBuf.WriteString(strconv.Itoa(fss.splitNo))
	}
}

type LogFileSpliter interface {
	OnWrite(writeBytes int)
	TrySplit() bool
	AppendDirAndName(dirBuf, nameBuf *bytes.Buffer)
	MoveLast(dirBuf, nameBuf *bytes.Buffer) error
}

type NowTime func() time.Time

func defaultNow() time.Time {
	return time.Now()
}

type LogFileHourSpliter struct {
	nextFileTime time.Time
	nowTime      NowTime
	sizeSpliter  *FileSizeSpliter
}

func getHourTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, t.Hour(), 0, 0, 0, time.Local)
}

func NewLogFileHourSpliter(timeFunc NowTime, maxFileSize int) *LogFileHourSpliter {
	if timeFunc == nil {
		timeFunc = defaultNow
	}

	return &LogFileHourSpliter{
		nextFileTime: getHourTime(timeFunc()),
		nowTime:      timeFunc,
		sizeSpliter:  NewFileSizeSpliter(maxFileSize),
	}
}

func (lfhs *LogFileHourSpliter) OnWrite(writeBytes int) {
	lfhs.sizeSpliter.OnWrite(writeBytes)
}

func (lfhs *LogFileHourSpliter) TrySplit() bool {
	now := lfhs.nowTime()
	if now.After(lfhs.nextFileTime) {
		lfhs.nextFileTime = lfhs.nextFileTime.Add(time.Hour)
		lfhs.sizeSpliter.Restart()
		return true
	}

	return lfhs.sizeSpliter.TrySplit()
}

func (lfhs *LogFileHourSpliter) AppendDirAndName(dirBuf, nameBuf *bytes.Buffer) {
	now := lfhs.nowTime()
	lfhs.appendDir(dirBuf, &now)
	lfhs.appendName(nameBuf, &now)
	lfhs.sizeSpliter.AppendDirAndName(dirBuf, nameBuf)
}

func (lfhs *LogFileHourSpliter) appendDir(dirBuf *bytes.Buffer, t *time.Time) {
	if dirBuf.Len() > 0 {
		dirBuf.WriteByte('/')
	}
	dirBuf.WriteString(t.Format("2006-01-02"))
}

func (lfhs *LogFileHourSpliter) appendName(nameBuf *bytes.Buffer, t *time.Time) {
	if nameBuf.Len() > 0 {
		nameBuf.WriteByte('-')
	}
	nameBuf.WriteString(strconv.Itoa(t.Hour()))
}

func (lfhs *LogFileHourSpliter) MoveLast(dirBuf, nameBuf *bytes.Buffer) error {
	lfhs.restart()

	now := lfhs.nowTime()
	lfhs.appendDir(dirBuf, &now)

	fileInfos, err := ioutil.ReadDir(dirBuf.String())
	if err != nil {
		return err
	}

	regStr := fmt.Sprintf(`.*?%d(?:-([\d]+))?.*`, now.Hour())
	reg := regexp.MustCompile(regStr)

	maxSplitNo := 0
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		filename := fileInfo.Name()
		matches := reg.FindStringSubmatch(filename)
		if len(matches) < 2 {
			continue
		}

		splitStr := matches[1]
		if splitStr == "" {
			continue
		}

		splitNo, err := strconv.ParseInt(splitStr, 10, 64)
		if err != nil {
			continue
		}

		if splitNo > int64(maxSplitNo) {
			maxSplitNo = int(splitNo)
		}
	}

	lfhs.sizeSpliter.SetSplitNo(maxSplitNo)

	return nil
}

func (lfhs *LogFileHourSpliter) restart() {
	lfhs.nextFileTime = getHourTime(lfhs.nowTime()).Add(time.Hour)
	lfhs.sizeSpliter.Restart()
}

type LogFileDateSpliter struct {
	nextFileTime time.Time
	nowTime      NowTime
	sizeSpliter  *FileSizeSpliter
}

func getDateTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func NewLogFileDateSpliter(timeFunc NowTime, maxFileSize int) *LogFileDateSpliter {
	if timeFunc == nil {
		timeFunc = defaultNow
	}

	return &LogFileDateSpliter{
		nextFileTime: getDateTime(timeFunc()),
		nowTime:      timeFunc,
		sizeSpliter:  NewFileSizeSpliter(maxFileSize),
	}
}

func (lfds *LogFileDateSpliter) OnWrite(writeBytes int) {
	lfds.sizeSpliter.OnWrite(writeBytes)
}

func (lfds *LogFileDateSpliter) TrySplit() bool {
	now := lfds.nowTime()
	if now.After(lfds.nextFileTime) {
		lfds.nextFileTime = lfds.nextFileTime.AddDate(0, 0, 1)
		lfds.sizeSpliter.Restart()
		return true
	}

	return lfds.sizeSpliter.TrySplit()
}

func (lfds *LogFileDateSpliter) AppendDirAndName(dirBuf, nameBuf *bytes.Buffer) {
	now := lfds.nowTime()
	lfds.appendDir(dirBuf, &now)
	lfds.appendName(nameBuf, &now)
	lfds.sizeSpliter.AppendDirAndName(dirBuf, nameBuf)
}

func (lfds *LogFileDateSpliter) appendDir(dirBuf *bytes.Buffer, t *time.Time) {
}

func (lfds *LogFileDateSpliter) appendName(nameBuf *bytes.Buffer, t *time.Time) {
	if nameBuf.Len() > 0 {
		nameBuf.WriteByte('-')
	}
	nameBuf.WriteString(t.Format("2006-01-02"))
}

func (lfds *LogFileDateSpliter) MoveLast(dirBuf, nameBuf *bytes.Buffer) error {
	lfds.restart()

	now := lfds.nowTime()
	lfds.appendDir(dirBuf, &now)

	fileInfos, err := ioutil.ReadDir(dirBuf.String())
	if err != nil {
		return err
	}

	regStr := fmt.Sprintf(`.*?%s(?:-([\d]+))?.*`, now.Format("2006-01-02"))
	reg := regexp.MustCompile(regStr)

	maxSplitNo := 0
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		filename := fileInfo.Name()
		matches := reg.FindStringSubmatch(filename)
		if len(matches) < 2 {
			continue
		}

		splitStr := matches[1]
		if splitStr == "" {
			continue
		}

		splitNo, err := strconv.ParseInt(splitStr, 10, 64)
		if err != nil {
			continue
		}

		if splitNo > int64(maxSplitNo) {
			maxSplitNo = int(splitNo)
		}
	}

	lfds.sizeSpliter.SetSplitNo(maxSplitNo)

	return nil
}

func (lfds *LogFileDateSpliter) restart() {
	lfds.nextFileTime = getDateTime(lfds.nowTime()).AddDate(0, 0, 1)
	lfds.sizeSpliter.Restart()
}
