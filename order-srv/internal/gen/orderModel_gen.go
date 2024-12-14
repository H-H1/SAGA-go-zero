// Code generated by goctl. DO NOT EDIT.

package gen

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	orderFieldNames          = builder.RawFieldNames(&Order{})
	orderRows                = strings.Join(orderFieldNames, ",")
	orderRowsExpectAutoSet   = strings.Join(stringx.Remove(orderFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	orderRowsWithPlaceHolder = strings.Join(stringx.Remove(orderFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	orderModel interface {
		Insert(ctx context.Context, data *Order) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Order, error)
		Update(ctx context.Context, data *Order) error
		Delete(ctx context.Context, id int64) error
		FindLastOneByUserIdGoodsId(userId, goodsId int64) (*Order, error)
	}

	defaultOrderModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Order struct {
		Id       int64 `db:"id"`
		UserId   int64 `db:"user_id"`
		GoodsId  int64 `db:"goods_id"`  // 商品id
		Num      int64 `db:"num"`       // 下单数量
		RowState int64 `db:"row_state"` // -1:下单回滚废弃 0:待支付
	}
)

func newOrderModel(conn sqlx.SqlConn) *defaultOrderModel {
	return &defaultOrderModel{
		conn:  conn,
		table: "`order`",
	}
}

func (m *defaultOrderModel) withSession(session sqlx.Session) *defaultOrderModel {
	return &defaultOrderModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`order`",
	}
}

func (m *defaultOrderModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultOrderModel) FindOne(ctx context.Context, id int64) (*Order, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", orderRows, m.table)
	var resp Order
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultOrderModel) Insert(ctx context.Context, data *Order) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, orderRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.GoodsId, data.Num, data.RowState)
	return ret, err
}

func (m *defaultOrderModel) Update(ctx context.Context, data *Order) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, orderRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.GoodsId, data.Num, data.RowState, data.Id)
	return err
}

func (m *defaultOrderModel) tableName() string {
	return m.table
}
func (m *defaultOrderModel) FindLastOneByUserIdGoodsId(userId, goodsId int64) (*Order, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and goods_id =? order by id desc limit 1 ", orderRows, m.table)
	var resp Order
	err := m.conn.QueryRow(&resp, query, userId, goodsId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
