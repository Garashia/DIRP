package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Node struct {
	Name     string
	Children []*Node
}

type ParseError struct {
	Pos int
	Msg string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at char %d: %s", e.Pos, e.Msg)
}

type InputSource struct {
	Text string
	File string
}

type Parser struct {
	src []rune
	pos int
}

func NewParser(input string) *Parser {
	return &Parser{src: []rune(input)}
}

func (p *Parser) ParseProgram() ([]*Node, error) {
	p.skipEntityLeadingWhitespace()
	nodes, err := p.parseEntities(false)
	if err != nil {
		return nil, err
	}
	p.skipEntityLeadingWhitespace()
	if !p.eof() {
		return nil, p.errf("unexpected trailing token %q", string(p.peek()))
	}
	return nodes, nil
}

func (p *Parser) parseEntities(stopAtRightBrace bool) ([]*Node, error) {
	var out []*Node
	for {
		p.skipEntityLeadingWhitespace()
		if p.eof() {
			if stopAtRightBrace {
				return nil, p.errf("missing closing '}'")
			}
			break
		}
		if stopAtRightBrace && p.peek() == '}' {
			break
		}

		entityNodes, err := p.parseEntity()
		if err != nil {
			return nil, err
		}
		out = append(out, entityNodes...)

		p.skipEntityLeadingWhitespace()
		if p.eof() {
			break
		}
		if stopAtRightBrace && p.peek() == '}' {
			break
		}

		// Lenient mode: allow implicit sibling separation right after a closed child block.
		// Example: "... }next_entity { ... }" is treated like "... }|next_entity { ... }".
		if p.allowImplicitSeparatorAfterRightBrace() {
			continue
		}

		if !p.consumeSeparator() {
			return nil, p.errf("expected separator ',' '|' or newline")
		}

		p.skipEntityLeadingWhitespace()
		if stopAtRightBrace && p.peekIfExists('}') {
			return nil, p.errf("empty entity before '}'")
		}
		if p.eof() {
			return nil, p.errf("dangling separator at end")
		}
		if p.peekIfExists(',') || p.peekIfExists('|') || p.peekIfExists('\n') {
			return nil, p.errf("empty entity between separators")
		}
	}

	if len(out) == 0 {
		return nil, p.errf("empty entity list")
	}
	return out, nil
}

func (p *Parser) parseEntity() ([]*Node, error) {
	names, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipEntityLeadingWhitespace()
	var children []*Node
	if p.peekIfExists('{') {
		p.pos++
		children, err = p.parseEntities(true)
		if err != nil {
			return nil, err
		}
		p.skipEntityLeadingWhitespace()
		if !p.peekIfExists('}') {
			return nil, p.errf("missing closing '}'")
		}
		p.pos++
	}

	nodes := make([]*Node, 0, len(names))
	for _, name := range names {
		nodes = append(nodes, &Node{Name: name, Children: children})
	}
	return nodes, nil
}

func (p *Parser) parseExpression() ([]string, error) {
	startPos := p.pos
	var prefixBuilder strings.Builder
	var suffixBuilder strings.Builder
	var templateVals []string
	templateSeen := false

	for !p.eof() {
		ch := p.peek()
		if ch == '{' || ch == '}' || ch == ',' || ch == '|' || ch == '\n' {
			break
		}

		if ch == '#' || ch == '@' {
			if p.pos+1 < len(p.src) && p.src[p.pos+1] == '(' {
				if templateSeen {
					return nil, p.errf("multiple templates in one entity are not allowed")
				}
				vals, err := p.parseTemplate()
				if err != nil {
					return nil, err
				}
				templateVals = vals
				templateSeen = true
				continue
			}
			return nil, p.errf("reserved character %q in directory name", string(ch))
		}

		if ch == '(' || ch == ')' {
			return nil, p.errf("reserved character %q in directory name", string(ch))
		}

		if !templateSeen {
			prefixBuilder.WriteRune(ch)
		} else {
			suffixBuilder.WriteRune(ch)
		}
		p.pos++
	}

	prefix := prefixBuilder.String()
	suffix := suffixBuilder.String()

	if !templateSeen {
		name := strings.TrimSpace(prefix)
		if name == "" {
			return nil, &ParseError{Pos: startPos, Msg: "empty directory name"}
		}
		return []string{name}, nil
	}

	expanded := make([]string, 0, len(templateVals))
	for _, v := range templateVals {
		name := strings.TrimSpace(prefix + v + suffix)
		if name == "" {
			return nil, &ParseError{Pos: startPos, Msg: "empty directory name after template expansion"}
		}
		expanded = append(expanded, name)
	}
	return expanded, nil
}

func (p *Parser) parseTemplate() ([]string, error) {
	if p.eof() {
		return nil, p.errf("expected template")
	}
	kind := p.peek()
	p.pos++
	if !p.peekIfExists('(') {
		return nil, p.errf("expected '(' after %q", string(kind))
	}
	p.pos++

	switch kind {
	case '#':
		return p.parseRangeFunc()
	case '@':
		return p.parseListFunc()
	default:
		return nil, p.errf("unknown template function %q", string(kind))
	}
}

func (p *Parser) parseRangeFunc() ([]string, error) {
	start, err := p.parseIntegerArg()
	if err != nil {
		return nil, err
	}
	if err := p.expectCommaInTemplate(); err != nil {
		return nil, err
	}
	end, err := p.parseIntegerArg()
	if err != nil {
		return nil, err
	}
	if err := p.expectCommaInTemplate(); err != nil {
		return nil, err
	}
	step, err := p.parseIntegerArg()
	if err != nil {
		return nil, err
	}

	p.skipTemplateWhitespace()
	if !p.peekIfExists(')') {
		return nil, p.errf("expected ')' for range function")
	}
	p.pos++

	if step == 0 {
		return nil, p.errf("range step must not be zero")
	}
	if start < end && step < 0 {
		return nil, p.errf("range step direction is invalid")
	}
	if start > end && step > 0 {
		return nil, p.errf("range step direction is invalid")
	}

	var out []string
	if step > 0 {
		for i := start; i <= end; i += step {
			out = append(out, strconv.Itoa(i))
		}
	} else {
		for i := start; i >= end; i += step {
			out = append(out, strconv.Itoa(i))
		}
	}
	if len(out) == 0 {
		return nil, p.errf("range generated no values")
	}
	return out, nil
}

func (p *Parser) parseListFunc() ([]string, error) {
	var items []string

	for {
		p.skipTemplateWhitespace()
		if p.peekIfExists(')') {
			break
		}

		item, err := p.parseListItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)

		p.skipTemplateWhitespace()
		if p.peekIfExists(')') {
			break
		}
		if err := p.expectCommaInTemplate(); err != nil {
			return nil, err
		}
	}

	if !p.peekIfExists(')') {
		return nil, p.errf("expected ')' for list function")
	}
	p.pos++

	if len(items) == 0 {
		return nil, p.errf("list function requires at least one item")
	}
	return items, nil
}

func (p *Parser) parseListItem() (string, error) {
	start := p.pos
	var b strings.Builder

	for !p.eof() {
		ch := p.peek()
		if ch == ',' || ch == ')' {
			break
		}
		if ch == '\n' || ch == '\r' {
			return "", p.errf("newline inside list item is not allowed")
		}
		if ch == '{' || ch == '}' || ch == '(' || ch == '#' || ch == '@' || ch == '|' {
			return "", p.errf("reserved character %q in list item", string(ch))
		}
		b.WriteRune(ch)
		p.pos++
	}

	item := strings.TrimSpace(b.String())
	if item == "" {
		return "", &ParseError{Pos: start, Msg: "empty list item"}
	}
	return item, nil
}

func (p *Parser) parseIntegerArg() (int, error) {
	p.skipTemplateWhitespace()
	start := p.pos

	if p.peekIfExists('-') {
		p.pos++
	}
	if p.eof() || !unicode.IsDigit(p.peek()) {
		return 0, &ParseError{Pos: start, Msg: "expected integer"}
	}

	for !p.eof() && unicode.IsDigit(p.peek()) {
		p.pos++
	}
	raw := string(p.src[start:p.pos])
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, &ParseError{Pos: start, Msg: "invalid integer"}
	}
	return v, nil
}

func (p *Parser) expectCommaInTemplate() error {
	p.skipTemplateWhitespace()
	if !p.peekIfExists(',') {
		return p.errf("expected ',' between template arguments")
	}
	p.pos++
	return nil
}

func (p *Parser) skipEntityLeadingWhitespace() {
	for !p.eof() {
		ch := p.peek()
		if ch == ',' || ch == '|' {
			break
		}
		if !unicode.IsSpace(ch) {
			break
		}
		p.pos++
	}
}

func (p *Parser) skipTemplateWhitespace() {
	for !p.eof() {
		ch := p.peek()
		if ch == '\n' || ch == '\r' || ch == '\t' || ch == ' ' {
			p.pos++
			continue
		}
		break
	}
}

func (p *Parser) consumeSeparator() bool {
	if p.eof() {
		return false
	}
	ch := p.peek()
	if ch == ',' || ch == '|' || ch == '\n' {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) allowImplicitSeparatorAfterRightBrace() bool {
	if p.pos <= 0 || p.eof() {
		return false
	}
	i := p.pos - 1
	for i >= 0 {
		ch := p.src[i]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			i--
			continue
		}
		break
	}
	if i < 0 || p.src[i] != '}' {
		return false
	}
	ch := p.peek()
	if ch == '{' || ch == '}' || ch == ',' || ch == '|' || ch == '\n' {
		return false
	}
	return true
}

func (p *Parser) eof() bool {
	return p.pos >= len(p.src)
}

func (p *Parser) peek() rune {
	return p.src[p.pos]
}

func (p *Parser) peekIfExists(ch rune) bool {
	return !p.eof() && p.peek() == ch
}

func (p *Parser) errf(format string, args ...any) error {
	return &ParseError{
		Pos: p.pos,
		Msg: fmt.Sprintf(format, args...),
	}
}

func mkdirTree(root string, nodes []*Node) error {
	for _, n := range nodes {
		path := filepath.Join(root, n.Name)
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("failed to create %q: %w", path, err)
		}
		if err := mkdirTree(path, n.Children); err != nil {
			return err
		}
	}
	return nil
}

func readInput(path string, inline string) (*InputSource, error) {
	if path != "" && inline != "" {
		return nil, errors.New("use either -f or inline DSL argument, not both")
	}
	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q: %w", path, err)
		}
		return &InputSource{Text: string(b), File: path}, nil
	}
	if inline != "" {
		return &InputSource{Text: inline}, nil
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to inspect stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, errors.New("no input provided; pass DSL as argument, -f file, or pipe stdin")
	}
	var stdinBuilder strings.Builder
	buf := make([]byte, 4096)
	for {
		n, readErr := os.Stdin.Read(buf)
		if n > 0 {
			stdinBuilder.Write(buf[:n])
		}
		if readErr != nil {
			break
		}
	}
	return &InputSource{Text: stdinBuilder.String()}, nil
}

func lineColAtRuneOffset(text string, runeOffset int) (int, int) {
	if runeOffset < 0 {
		return 1, 1
	}
	runes := []rune(text)
	if runeOffset > len(runes) {
		runeOffset = len(runes)
	}
	line := 1
	col := 1
	for i := 0; i < runeOffset; i++ {
		if runes[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func lineTextAt(text string, targetLine int) string {
	if targetLine < 1 {
		return ""
	}
	lines := strings.Split(text, "\n")
	if targetLine > len(lines) {
		return ""
	}
	return strings.TrimRight(lines[targetLine-1], "\r")
}

func printParseErrorWithLocation(src *InputSource, pe *ParseError) {
	if src != nil && src.File != "" {
		line, col := lineColAtRuneOffset(src.Text, pe.Pos)
		fmt.Fprintf(os.Stderr, "%s:%d:%d: parse error: %s\n", src.File, line, col, pe.Msg)

		lineText := lineTextAt(src.Text, line)
		if lineText != "" {
			fmt.Fprintln(os.Stderr, lineText)
			if col < 1 {
				col = 1
			}
			fmt.Fprintf(os.Stderr, "%s^\n", strings.Repeat(" ", col-1))
		}
		return
	}
	fmt.Fprintln(os.Stderr, "error:", pe.Error())
}

func main() {
	var (
		root string
		file string
	)
	flag.StringVar(&root, "root", ".", "root directory where tree will be created")
	flag.StringVar(&file, "f", "", "path to .dirp file")
	flag.Parse()

	inline := strings.Join(flag.Args(), " ")
	src, err := readInput(file, inline)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if !utf8.ValidString(src.Text) {
		fmt.Fprintln(os.Stderr, "error: input is not valid UTF-8")
		os.Exit(1)
	}

	parser := NewParser(src.Text)
	nodes, err := parser.ParseProgram()
	if err != nil {
		var pe *ParseError
		if errors.As(err, &pe) {
			printParseErrorWithLocation(src, pe)
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		os.Exit(1)
	}

	if err := mkdirTree(root, nodes); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
