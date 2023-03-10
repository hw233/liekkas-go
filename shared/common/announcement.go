package common

import (
	"context"
	"database/sql"
	"fmt"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/mysql"
)

const (
	AnnouncementModuleForeplay = iota
	AnnouncementModuleGame
)

type Announcement struct {
	Id            int64  `json:"id" db:"id"`
	Type          int32  `json:"type" db:"type"`     //分页类型 0活动
	Module        int32  `json:"module" db:"module"` //所属模块, 0登录前 1登陆后
	Title         string `json:"title" db:"title"`
	Content       string `json:"content" db:"content"`
	Tag           int32  `json:"tag" db:"tag"` //标签 0更新 1修复 2活动 3咨询 4其他
	Image         string `json:"image" db:"image"`
	StartTime     int64  `json:"start_time" db:"start_time"`
	EndTime       int64  `json:"end_time" db:"end_time"`
	ShowStartTime int64  `json:"show_start_time" db:"show_start_time"`
	ShowEndTime   int64  `json:"show_end_time" db:"show_end_time"`
	Priority      int32  `json:"priority" db:"priority"`
}

func NewAnnouncement(id int64, aType, module int32, title, content, image string, tag int32,
	startTime, endTime, showStartTime, showEndTime int64, priority int32) *Announcement {
	return &Announcement{
		Id:            id,
		Type:          aType,
		Module:        module,
		Title:         title,
		Content:       content,
		Tag:           tag,
		Image:         image,
		StartTime:     startTime,
		EndTime:       endTime,
		ShowStartTime: showStartTime,
		ShowEndTime:   showEndTime,
		Priority:      priority,
	}
}

func (a *Announcement) VOAnnouncement() *pb.VOAnnouncement {
	return &pb.VOAnnouncement{
		Id:            a.Id,
		Type:          a.Type,
		Title:         a.Title,
		Content:       a.Content,
		Tag:           a.Tag,
		Image:         a.Image,
		StartTime:     a.StartTime,
		EndTime:       a.EndTime,
		ShowStartTime: a.ShowStartTime,
		ShowEndTime:   a.ShowEndTime,
		Priority:      a.Priority,
	}
}

type Banner struct {
	Id            int64  `json:"id" db:"id"`
	Type          int32  `json:"type" db:"type"`     //分页类型
	Module        int32  `json:"module" db:"module"` //所属模块, 0登录前 1登陆后
	Image         string `json:"image" db:"image"`
	Jump          string `json:"jump" db:"jump"`
	StartTime     int64  `json:"start_time" db:"start_time"`
	EndTime       int64  `json:"end_time" db:"end_time"`
	ShowStartTime int64  `json:"show_start_time" db:"show_start_time"`
	ShowEndTime   int64  `json:"show_end_time" db:"show_end_time"`
	Priority      int32  `json:"priority" db:"priority"`
}

func NewBanner(id int64, aType, module int32, image, jump string, startTime, endTime,
	showStartTime, showEndTime int64, priority int32) *Banner {

	return &Banner{
		Id:            id,
		Type:          aType,
		Module:        module,
		Image:         image,
		StartTime:     startTime,
		EndTime:       endTime,
		ShowStartTime: showStartTime,
		ShowEndTime:   showEndTime,
		Priority:      priority,
	}
}

func (b *Banner) VOBanner() *pb.VOBanner {
	return &pb.VOBanner{
		Id:            b.Id,
		Type:          b.Type,
		Image:         b.Image,
		StartTime:     b.StartTime,
		EndTime:       b.EndTime,
		ShowStartTime: b.ShowStartTime,
		ShowEndTime:   b.ShowEndTime,
		Priority:      b.Priority,
	}
}

type Caution struct {
	Id        int64  `json:"id"`
	Content   string `json:"content"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

func NewCaution(id int64, content string, startTime, endTime int64) *Caution {
	return &Caution{
		Id:        id,
		Content:   content,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

func (c *Caution) VOForeplayCaution() *pb.VOForeplayCaution {
	return &pb.VOForeplayCaution{
		Id:        c.Id,
		Content:   c.Content,
		StartTime: c.StartTime,
		EndTime:   c.EndTime,
	}
}

const (
	AnnouncementDBName = "announcement"
	BannerDBName       = "banner"
	CautionDBName      = "caution"
)

var (
	ErrAnnouncementDBError        = errors.New("announcement db error")
	ErrAnnouncementNotFound       = errors.New("target not found")
	ErrAnnouncementNothingUpdated = errors.New("nothing updated")
)

func DBAddAnnouncement(ctx context.Context, annc *Announcement, db *sql.DB) error {
	columns, placeholders, values, err := mysql.GenInsertParams(annc)
	if err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("insert into `%s`(%s) values(%s)", AnnouncementDBName, columns, placeholders)
	_, err = db.ExecContext(ctx, cmdStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func DBUpdateAnnouncement(ctx context.Context, annc *Announcement, db *sql.DB) error {
	params, values, err := mysql.GenUpdateParams(annc)
	if err != nil {
		return err
	}

	values = append(values, annc.Id)
	cmdStr := fmt.Sprintf("update `%s` set %s where id=?", AnnouncementDBName, params)
	ret, err := db.ExecContext(ctx, cmdStr, values...)
	if err != nil {
		return err
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return ErrAnnouncementDBError
	}
	if affected <= 0 {
		return ErrAnnouncementNothingUpdated
	}

	return nil
}

func DBLoadAnnouncement(ctx context.Context, module int32, endLimit int64, db *sql.DB) ([]*Announcement, error) {
	columnsStr := mysql.GenObjectColumns(&Announcement{})
	cmdStr := fmt.Sprintf("select %s from `%s` where `end_time` > ? and `module` = ?", columnsStr, AnnouncementDBName)
	rows, err := db.QueryContext(ctx, cmdStr, endLimit, module)

	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnLen := len(columns)

	anncs := []*Announcement{}

	for rows.Next() {
		values := make([]interface{}, 0, columnLen)
		for i := 0; i < columnLen; i++ {
			values = append(values, &[]byte{})
		}

		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		annc := &Announcement{}
		err = mysql.SetObjectValue(annc, values)
		if err != nil {
			return nil, err
		}

		anncs = append(anncs, annc)
	}

	return anncs, nil
}

func DBDeleteAnnouncement(ctx context.Context, id int64, db *sql.DB) error {
	cmdStr := fmt.Sprintf("delete from `%s` where id=?", AnnouncementDBName)
	ret, err := db.ExecContext(ctx, cmdStr, id)
	if err != nil {
		return ErrAnnouncementDBError
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return ErrAnnouncementDBError
	}
	if affected <= 0 {
		return ErrAnnouncementNotFound
	}

	return nil
}

func DBAddBanner(ctx context.Context, banner *Banner, db *sql.DB) error {
	columns, placeholders, values, err := mysql.GenInsertParams(banner)
	if err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("insert into `%s`(%s) values(%s)", BannerDBName, columns, placeholders)
	_, err = db.ExecContext(ctx, cmdStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func DBUpdateBanner(ctx context.Context, banner *Banner, db *sql.DB) error {
	params, values, err := mysql.GenUpdateParams(banner)
	if err != nil {
		return err
	}

	values = append(values, banner.Id)
	cmdStr := fmt.Sprintf("update `%s` set %s where id=?", BannerDBName, params)
	ret, err := db.ExecContext(ctx, cmdStr, values...)
	if err != nil {
		return err
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return ErrAnnouncementDBError
	}
	if affected <= 0 {
		return ErrAnnouncementNothingUpdated
	}

	return nil
}

func DBLoadBanner(ctx context.Context, module int32, endLimit int64, db *sql.DB) ([]*Banner, error) {
	columnsStr := mysql.GenObjectColumns(&Banner{})
	cmdStr := fmt.Sprintf("select %s from `%s` where `end_time` > ? and `module` = ?", columnsStr, BannerDBName)
	rows, err := db.QueryContext(ctx, cmdStr, endLimit, module)

	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnLen := len(columns)

	banners := []*Banner{}

	for rows.Next() {
		values := make([]interface{}, 0, columnLen)
		for i := 0; i < columnLen; i++ {
			values = append(values, &[]byte{})
		}

		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		banner := &Banner{}
		err = mysql.SetObjectValue(banner, values)
		if err != nil {
			return nil, err
		}

		banners = append(banners, banner)
	}

	return banners, nil
}

func DBDeleteBanner(ctx context.Context, id int64, db *sql.DB) error {
	cmdStr := fmt.Sprintf("delete from `%s` where id=?", BannerDBName)
	ret, err := db.ExecContext(ctx, cmdStr, id)
	if err != nil {
		return ErrAnnouncementDBError
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return ErrAnnouncementDBError
	}
	if affected <= 0 {
		return ErrAnnouncementNotFound
	}

	return nil
}

func DBAddCaution(ctx context.Context, caution *Caution, db *sql.DB) error {
	columns, placeholders, values, err := mysql.GenInsertParams(caution)
	if err != nil {
		return err
	}

	values = append(values, caution.Id)
	cmdStr := fmt.Sprintf("insert into `%s`(%s) values(%s)", CautionDBName, columns, placeholders)
	_, err = db.ExecContext(ctx, cmdStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func DBUpdateCaution(ctx context.Context, caution *Caution, db *sql.DB) error {
	params, values, err := mysql.GenUpdateParams(caution)
	if err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("update `%s` set %s where id=?", CautionDBName, params)
	ret, err := db.ExecContext(ctx, cmdStr, values...)
	if err != nil {
		return err
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return ErrAnnouncementDBError
	}
	if affected <= 0 {
		return ErrAnnouncementNothingUpdated
	}

	return nil
}

func DBLoadCaution(ctx context.Context, endLimit int64, db *sql.DB) ([]*Caution, error) {
	columnsStr := mysql.GenObjectColumns(&Caution{})
	cmdStr := fmt.Sprintf("select %s from `%s` where `end_time` > ?", columnsStr, CautionDBName)
	rows, err := db.QueryContext(ctx, cmdStr, endLimit)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnLen := len(columns)

	cautions := []*Caution{}

	for rows.Next() {
		values := make([]interface{}, 0, columnLen)
		for i := 0; i < columnLen; i++ {
			values = append(values, &[]byte{})
		}

		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		caution := &Caution{}
		err = mysql.SetObjectValue(caution, values)
		if err != nil {
			return nil, err
		}

		cautions = append(cautions, caution)
	}

	return cautions, nil
}

func DBDeleteCaution(ctx context.Context, id int64, db *sql.DB) error {
	cmdStr := fmt.Sprintf("delete from `%s` where id=?", CautionDBName)
	ret, err := db.ExecContext(ctx, cmdStr, id)
	if err != nil {
		return ErrAnnouncementDBError
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return ErrAnnouncementDBError
	}
	if affected <= 0 {
		return ErrAnnouncementNotFound
	}

	return nil
}
