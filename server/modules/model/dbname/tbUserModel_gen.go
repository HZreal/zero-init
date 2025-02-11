// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.2

package dbname

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tbUserFieldNames          = builder.RawFieldNames(&TbUser{})
	tbUserRows                = strings.Join(tbUserFieldNames, ",")
	tbUserRowsExpectAutoSet   = strings.Join(stringx.Remove(tbUserFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tbUserRowsWithPlaceHolder = strings.Join(stringx.Remove(tbUserFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	tbUserModel interface {
		Insert(ctx context.Context, data *TbUser) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*TbUser, error)
		FindOneByMobile(ctx context.Context, mobile string) (*TbUser, error)
		Update(ctx context.Context, data *TbUser) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultTbUserModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TbUser struct {
		Id         uint64    `db:"id"`
		Name       string    `db:"name"`     // 用户姓名
		Gender     uint64    `db:"gender"`   // 用户性别
		Mobile     string    `db:"mobile"`   // 用户电话
		Password   string    `db:"password"` // 用户密码
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func newTbUserModel(conn sqlx.SqlConn) *defaultTbUserModel {
	return &defaultTbUserModel{
		conn:  conn,
		table: "`tb_user`",
	}
}

func (m *defaultTbUserModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTbUserModel) FindOne(ctx context.Context, id uint64) (*TbUser, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tbUserRows, m.table)
	var resp TbUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTbUserModel) FindOneByMobile(ctx context.Context, mobile string) (*TbUser, error) {
	var resp TbUser
	query := fmt.Sprintf("select %s from %s where `mobile` = ? limit 1", tbUserRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, mobile)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTbUserModel) Insert(ctx context.Context, data *TbUser) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, tbUserRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.Gender, data.Mobile, data.Password)
	return ret, err
}

func (m *defaultTbUserModel) Update(ctx context.Context, newData *TbUser) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tbUserRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Name, newData.Gender, newData.Mobile, newData.Password, newData.Id)
	return err
}

func (m *defaultTbUserModel) tableName() string {
	return m.table
}
