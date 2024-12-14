package gen

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StockModel = (*customStockModel)(nil)

type (
	// StockModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStockModel.
	StockModel interface {
		stockModel
	}

	customStockModel struct {
		*defaultStockModel
	}
)

// NewStockModel returns a model for the database table.
func NewStockModel(conn sqlx.SqlConn) StockModel {
	return &customStockModel{
		defaultStockModel: newStockModel(conn),
	}
}
