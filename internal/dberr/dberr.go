package dberr

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// 数据库特定错误
var (
	ErrNoRows           = sql.ErrNoRows
	ErrDuplicateKey     = errors.New("记录已存在")
	ErrForeignKey       = errors.New("关联记录不存在")
	ErrCheckViolation   = errors.New("数据校验失败")
	ErrNotNullViolation = errors.New("必填字段缺失")
	ErrInvalidData      = errors.New("数据格式错误")
	ErrConnection       = errors.New("数据库连接失败")

	// Paste 相关错误
	ErrShortLinkDuplicate = errors.New("短链接已存在")
)

// PostgreSQL 错误代码
const (
	PgErrUniqueViolation     = "23505"
	PgErrForeignKeyViolation = "23503"
	PgErrCheckViolation      = "23514"
	PgErrNotNullViolation    = "23502"

	PgErrInvalidTextRepresentation = "22P02"
	PgErrNumericValueOutOfRange    = "22003"

	PgErrConnectionException    = "08000"
	PgErrConnectionDoesNotExist = "08003"
	PgErrConnectionFailure      = "08006"

	PgErrSyntaxError     = "42601"
	PgErrUndefinedColumn = "42703"
	PgErrUndefinedTable  = "42P01"
)

// ParseDBError 解析数据库错误并转换为业务错误
func ParseDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRows
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return parsePgError(pqErr)
	}

	return err
}

func parsePgError(pqErr *pq.Error) error {
	switch pqErr.Code {
	case PgErrUniqueViolation:
		if pqErr.Constraint == "pastes_short_link_key" {
			return ErrShortLinkDuplicate
		}
		return ErrDuplicateKey

	case PgErrForeignKeyViolation:
		return ErrForeignKey

	case PgErrCheckViolation:
		return ErrCheckViolation

	case PgErrNotNullViolation:
		return ErrNotNullViolation

	case PgErrInvalidTextRepresentation, PgErrNumericValueOutOfRange:
		return ErrInvalidData

	case PgErrConnectionException, PgErrConnectionDoesNotExist, PgErrConnectionFailure:
		return ErrConnection

	case PgErrUndefinedColumn, PgErrUndefinedTable, PgErrSyntaxError:
		return pqErr

	default:
		return pqErr
	}
}

// IsDuplicateKeyError 检查是否是唯一键冲突错误
func IsDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == PgErrUniqueViolation
	}
	return false
}

// IsNotFoundError 检查是否是记录不存在错误
func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
