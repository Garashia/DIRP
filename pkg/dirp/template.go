package dirp

import (
	"strconv"
	"strings"
)

// expandName は entity 名に含まれる template を展開する入口。
// basePos は元入力上の開始位置で、詳細エラー位置へ変換するときに使う。
func expandName(pattern string, basePos int) ([]string, error) {
	if pattern == "" {
		return nil, NewError(ErrEmptyName, basePos, "empty entity name")
	}
	return expandPattern(pattern, basePos)
}

// expandPattern は文字列中の最初の template を1つ展開し、
// 残り suffix を再帰展開して直積を作る。
func expandPattern(s string, basePos int) ([]string, error) {
	hashIdx := strings.Index(s, "#(")
	atIdx := strings.Index(s, "@(")

	if hashIdx == -1 && atIdx == -1 {
		return []string{s}, nil
	}

	// "#(" と "@(" のうち、より左側にある template を先に処理する。
	start := hashIdx
	isHash := true
	if start == -1 || (atIdx != -1 && atIdx < hashIdx) {
		start = atIdx
		isHash = false
	}

	end := strings.Index(s[start:], ")")
	if end == -1 {
		return nil, NewError(ErrMissingDelimiter, basePos+start, "missing closing ')' in template")
	}
	end = start + end

	// prefix + item + tail で最終候補を作る。
	prefix := s[:start]
	body := s[start+2 : end]
	suffix := s[end+1:]

	var items []string
	var err error
	if isHash {
		items, err = expandRange(body, basePos+start+2)
	} else {
		items, err = expandList(body, basePos+start+2)
	}
	if err != nil {
		return nil, err
	}

	tails, err := expandPattern(suffix, basePos+end+1)
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(items)*len(tails))
	for _, item := range items {
		for _, tail := range tails {
			out = append(out, prefix+item+tail)
		}
	}
	return out, nil
}

// expandRange は #(end), #(start,end), #(start,end,step) を展開する。
// - 1引数: start=1, stepは end に向かう向きで自動決定
// - 2引数: stepは end に向かう向きで自動決定
// - 3引数: stepを明示指定
func expandRange(body string, bodyPos int) ([]string, error) {
	parts := strings.Split(body, ",")
	if len(parts) < 1 || len(parts) > 3 {
		return nil, NewError(ErrInvalidTemplate, bodyPos, "invalid range template: expected 1 to 3 args")
	}

	argPos := make([]int, len(parts))
	offset := 0
	for i, part := range parts {
		argPos[i] = bodyPos + offset
		offset += len(part) + 1
	}

	toInt := func(idx int, label string) (int, error) {
		v, convErr := strconv.Atoi(strings.TrimSpace(parts[idx]))
		if convErr != nil {
			return 0, NewError(ErrInvalidArgument, argPos[idx], "invalid range %s %q", label, parts[idx])
		}
		return v, nil
	}

	var (
		start int
		end   int
		step  int
		err   error
	)

	switch len(parts) {
	case 1:
		start = 1
		end, err = toInt(0, "end")
		if err != nil {
			return nil, err
		}
		step = inferStep(start, end)
	case 2:
		start, err = toInt(0, "start")
		if err != nil {
			return nil, err
		}
		end, err = toInt(1, "end")
		if err != nil {
			return nil, err
		}
		step = inferStep(start, end)
	case 3:
		start, err = toInt(0, "start")
		if err != nil {
			return nil, err
		}
		end, err = toInt(1, "end")
		if err != nil {
			return nil, err
		}
		step, err = toInt(2, "step")
		if err != nil {
			return nil, err
		}
	default:
		return nil, NewError(ErrInvalidTemplate, bodyPos, "invalid range template: expected 1 to 3 args")
	}

	if step == 0 {
		return nil, NewError(ErrInvalidArgument, bodyPos, "range step must not be 0")
	}

	values := []string{}
	// step の符号で増加/減少の向きを決定。
	if step > 0 {
		for i := start; i <= end; i += step {
			values = append(values, strconv.Itoa(i))
		}
	} else {
		for i := start; i >= end; i += step {
			values = append(values, strconv.Itoa(i))
		}
	}

	if len(values) == 0 {
		return nil, NewError(ErrInvalidTemplate, bodyPos, "range step direction is invalid")
	}
	return values, nil
}

func inferStep(start int, end int) int {
	if end >= start {
		return 1
	}
	return -1
}

// expandList は @(a,b,c) を配列へ展開する。
func expandList(body string, bodyPos int) ([]string, error) {
	raw := strings.Split(body, ",")
	if len(raw) == 0 {
		return nil, NewError(ErrInvalidTemplate, bodyPos, "invalid list template")
	}

	out := make([]string, 0, len(raw))
	offset := 0 // body 内での相対位置（空要素エラー報告用）
	for _, item := range raw {
		v := strings.TrimSpace(item)
		if v == "" {
			return nil, NewError(ErrInvalidArgument, bodyPos+offset, "list template contains empty item")
		}
		out = append(out, v)
		offset += len(item) + 1
	}
	return out, nil
}
