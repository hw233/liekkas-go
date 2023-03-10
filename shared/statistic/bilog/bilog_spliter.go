package bilog

import (
	"bytes"
	"time"
)

type BILogSpliter struct {
	nextFileTime time.Time
}

func getHourTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, t.Hour(), 0, 0, 0, time.Local)
}

func NewBILogSpliter() *BILogSpliter {
	return &BILogSpliter{
		nextFileTime: getHourTime(time.Now()),
	}
}

func (bls *BILogSpliter) OnWrite(writeBytes int) {
}

func (bls *BILogSpliter) TrySplit() bool {
	now := time.Now()
	if now.After(bls.nextFileTime) {
		bls.nextFileTime = bls.nextFileTime.Add(time.Hour)
		return true
	}

	return false
}

func (bls *BILogSpliter) AppendDirAndName(dirBuf, nameBuf *bytes.Buffer) {
	now := time.Now()
	if nameBuf.Len() > 0 {
		nameBuf.WriteByte('_')
	}
	nameBuf.WriteString(now.Format("2006010215"))
}

func (bls *BILogSpliter) MoveLast(dirBuf, nameBuf *bytes.Buffer) error {
	return nil
}

type BILogDailySpliter struct {
	nextFileTime time.Time
}

func getDayTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func NewBILogDailySpliter() *BILogDailySpliter {
	return &BILogDailySpliter{
		nextFileTime: getDayTime(time.Now()),
	}
}

func (bds *BILogDailySpliter) OnWrite(writeBytes int) {
}

func (bds *BILogDailySpliter) TrySplit() bool {
	now := time.Now()
	if now.After(bds.nextFileTime) {
		bds.nextFileTime = bds.nextFileTime.Add(time.Hour * 24)
		return true
	}

	return false
}

func (bds *BILogDailySpliter) AppendDirAndName(dirBuf, nameBuf *bytes.Buffer) {
	now := time.Now()
	if nameBuf.Len() > 0 {
		nameBuf.WriteByte('_')
	}
	nameBuf.WriteString(now.Format("20060102"))
}

func (bds *BILogDailySpliter) MoveLast(dirBuf, nameBuf *bytes.Buffer) error {
	return nil
}
