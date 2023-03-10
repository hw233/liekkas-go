package mail

import (
	"context"
	"fmt"
	"mailserver/manager"
	"shared/common"
	"shared/utility/mysql"
	"shared/utility/servertime"
)

const (
	groupMailTableName    = "server_group_mails"
	personalMailTableName = "server_personal_mails"
)

func DBAddGroupMail(ctx context.Context, mail *common.ServerGroupMail) error {
	columns, placeholders, values, err := mysql.GenInsertParams(mail)
	if err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("insert into `%s`(%s) values(%s)", groupMailTableName, columns, placeholders)
	_, err = manager.DB.ExecContext(ctx, cmdStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func DBLoadGroupMail(ctx context.Context) ([]*common.ServerGroupMail, error) {
	now := servertime.Now().Unix()

	columnsStr := mysql.GenObjectColumns(&common.ServerGroupMail{})
	cmdStr := fmt.Sprintf("select %s from `%s` where `end_time` > ?", columnsStr, groupMailTableName)
	rows, err := manager.DB.QueryContext(ctx, cmdStr, now)

	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnLen := len(columns)

	var mails []*common.ServerGroupMail

	for rows.Next() {
		values := make([]interface{}, 0, columnLen)
		for i := 0; i < columnLen; i++ {
			values = append(values, &[]byte{})
		}

		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		mail := &common.ServerGroupMail{}
		err = mysql.SetObjectValue(mail, values)
		if err != nil {
			return nil, err
		}

		mails = append(mails, mail)
	}

	return mails, nil
}

func DBAddPersonalMail(ctx context.Context, mail *common.ServerPersonalMail) error {
	columns, placeholders, values, err := mysql.GenInsertParams(mail)
	if err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("insert into `%s`(%s) values(%s)", personalMailTableName, columns, placeholders)
	_, err = manager.DB.ExecContext(ctx, cmdStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func DBLoadPersonalMail(ctx context.Context) ([]*common.ServerPersonalMail, error) {
	now := servertime.Now().Unix()
	columnsStr := mysql.GenObjectColumns(&common.ServerPersonalMail{})
	cmdStr := fmt.Sprintf("select %s from `%s` where `end_time` > ?", columnsStr, personalMailTableName)
	rows, err := manager.DB.QueryContext(ctx, cmdStr, now)

	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnLen := len(columns)

	mails := []*common.ServerPersonalMail{}

	for rows.Next() {
		values := make([]interface{}, 0, columnLen)
		for i := 0; i < columnLen; i++ {
			values = append(values, &[]byte{})
		}

		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		mail := &common.ServerPersonalMail{}
		err = mysql.SetObjectValue(mail, values)
		if err != nil {
			return nil, err
		}

		mails = append(mails, mail)
	}

	return mails, nil
}
