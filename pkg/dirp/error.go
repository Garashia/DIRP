package dirp

import (
	"errors"
	"fmt"
)

// ErrorCode は C++ の error_code に寄せた「機械判定用コード」。
type ErrorCode int

const (
	ErrUnexpectedToken ErrorCode = iota + 1
	ErrMissingDelimiter
	ErrInvalidTemplate
	ErrInvalidArgument
	ErrEmptyName
)

func (c ErrorCode) Category() string {
	switch c {
	case ErrUnexpectedToken, ErrMissingDelimiter, ErrEmptyName:
		return "parse"
	case ErrInvalidTemplate, ErrInvalidArgument:
		return "template"
	default:
		return "unknown"
	}
}

func (c ErrorCode) String() string {
	switch c {
	case ErrUnexpectedToken:
		return "unexpected_token"
	case ErrMissingDelimiter:
		return "missing_delimiter"
	case ErrInvalidTemplate:
		return "invalid_template"
	case ErrInvalidArgument:
		return "invalid_argument"
	case ErrEmptyName:
		return "empty_name"
	default:
		return "unknown_error"
	}
}

type Error struct {
	Code  ErrorCode
	Pos   int // 入力文字列のバイト位置（一次情報）
	Line  int // 表示用（WithLineCol で補完）
	Col   int // 表示用（WithLineCol で補完）
	Msg   string
	Cause error
}

// Error は人間向けの簡易文字列表現。
// CLI 出力は FormatError を使うことで統一する。
func (e *Error) Error() string {
	if e.Line > 0 && e.Col > 0 {
		return fmt.Sprintf("%s:%s at %d:%d: %s", e.Code.Category(), e.Code.String(), e.Line, e.Col, e.Msg)
	}
	if e.Pos >= 0 {
		return fmt.Sprintf("%s:%s at pos %d: %s", e.Code.Category(), e.Code.String(), e.Pos, e.Msg)
	}
	return fmt.Sprintf("%s:%s: %s", e.Code.Category(), e.Code.String(), e.Msg)
}

// Unwrap で errors.As / errors.Is に対応させる。
func (e *Error) Unwrap() error {
	return e.Cause
}

// NewError は現在地点で新規エラーを作る。
func NewError(code ErrorCode, pos int, format string, args ...any) *Error {
	return &Error{
		Code: code,
		Pos:  pos,
		Msg:  fmt.Sprintf(format, args...),
	}
}

// WrapError は下位エラーを保持したまま文脈を付加する。
func WrapError(code ErrorCode, pos int, cause error, format string, args ...any) *Error {
	return &Error{
		Code:  code,
		Pos:   pos,
		Msg:   fmt.Sprintf(format, args...),
		Cause: cause,
	}
}

// AsError は一般 error から dirp.Error を安全に取り出す。
func AsError(err error) (*Error, bool) {
	var de *Error
	if errors.As(err, &de) {
		return de, true
	}
	return nil, false
}

// WithLineCol は Pos だけ持つエラーへ line/col を付与する。
func WithLineCol(err error, input string) error {
	de, ok := AsError(err)
	if !ok {
		return err
	}
	line, col := posToLineCol(input, de.Pos)
	out := *de
	out.Line = line
	out.Col = col
	return &out
}

// FormatError は path:line:col を含む統一フォーマットを返す。
func FormatError(path string, err error) string {
	de, ok := AsError(err)
	if !ok {
		return err.Error()
	}
	msg := de.Msg
	if de.Cause != nil {
		msg = msg + ": " + de.Cause.Error()
	}
	if de.Line > 0 && de.Col > 0 {
		return fmt.Sprintf("%s:%d:%d: %s error: %s", path, de.Line, de.Col, de.Code.Category(), msg)
	}
	return fmt.Sprintf("%s: %s error: %s", path, de.Code.Category(), msg)
}

// posToLineCol はバイト位置を line/col に変換する。
// 入力は UTF-8 を想定し、改行のみ特別扱いする。
func posToLineCol(input string, pos int) (int, int) {
	if pos < 0 {
		return 1, 1
	}
	if pos > len(input) {
		pos = len(input)
	}
	line, col := 1, 1
	for i := 0; i < pos; i++ {
		if input[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}
