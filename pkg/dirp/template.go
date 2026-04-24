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

// expandRange は #(start,end,step) を厳密に展開する。
func expandRange(body string, bodyPos int) ([]string, error) {
	parts := strings.Split(body, ",")
	if len(parts) != 3 {
		return nil, NewError(ErrInvalidTemplate, bodyPos, "invalid range template: expected 3 args")
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, NewError(ErrInvalidArgument, bodyPos, "invalid range start %q", parts[0])
	}
	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, NewError(ErrInvalidArgument, bodyPos+len(parts[0])+1, "invalid range end %q", parts[1])
	}
	step, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return nil, NewError(ErrInvalidArgument, bodyPos+len(parts[0])+len(parts[1])+2, "invalid range step %q", parts[2])
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
		return nil, NewError(ErrInvalidTemplate, bodyPos, "range expansion produced no values")
	}
	return values, nil
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
