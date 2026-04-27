package dirp

// Node は最終的に生成されるディレクトリ木の1要素。
type Node struct {
	Name     string
	Children []*Node
}

// Parser は Lexer から受け取った token 列を Node 木へ変換する。
// cur/peek の2トークン先読みで簡潔に構文判定する。
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

// NewParser は cur/peek を埋めた状態で parser を返す。
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken は 1 トークン進める。
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parse はトップレベルの entity 群を最後まで読む。
func (p *Parser) Parse() ([]*Node, error) {
	nodes := []*Node{}
	for p.curToken.Type != TokenEOF {
		parsed, err := p.parseEntity()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, parsed...)

		// 兄弟区切りは ',' と '|' を許可。
		if p.curToken.Type == TokenComma || p.curToken.Type == TokenPipe {
			p.nextToken()
		}
	}
	return nodes, nil
}

// parseEntity は単一 entity（必要ならテンプレ展開後に複数）を読む。
func (p *Parser) parseEntity() ([]*Node, error) {
	namePos := p.curToken.Pos
	if p.curToken.Type != TokenString && p.curToken.Type != TokenAt {
		if p.curToken.Type == TokenHash {
			return nil, NewError(ErrUnexpectedToken, p.curToken.Pos, "range template '#(...)' cannot appear at beginning of entity name")
		}
		return nil, NewError(ErrUnexpectedToken, p.curToken.Pos, "expected entity name, got token %q", p.curToken.Literal)
	}

	// lexer は template 記号を分割するため、ここで名前を再結合する。
	rawName, err := p.parseEntityName()
	if err != nil {
		return nil, err
	}

	// 例: api_#(1,3,1) -> api_1, api_2, api_3
	names, err := expandName(rawName, namePos)
	if err != nil {
		return nil, WrapError(ErrInvalidTemplate, namePos, err, "invalid entity name %q", rawName)
	}

	nodes := make([]*Node, 0, len(names))
	for _, name := range names {
		nodes = append(nodes, &Node{Name: name})
	}
	if p.curToken.Type == TokenLBrace {
		// 子ブロックを全て読み終わると curToken は '}' の次を指す。
		p.nextToken()
		children, err := p.parseUntilRBrace()
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			// 展開された各ノードが同じ子を共有すると片方の編集が他方へ波及するため、
			// 深いコピーで独立させる。
			n.Children = cloneNodes(children)
		}
	}

	return nodes, nil
}

// parseEntityName は entity 名を token から再構成する。
// 文頭 "@(...)" は許可し、文頭 "#(...)" は不許可とする。
func (p *Parser) parseEntityName() (string, error) {
	name := ""
	for {
		switch p.curToken.Type {
		case TokenString:
			name += p.curToken.Literal
			p.nextToken()
		case TokenAt:
			lit, err := p.parseTemplateLiteral()
			if err != nil {
				return "", err
			}
			name += lit
		case TokenHash:
			if name == "" {
				return "", NewError(ErrUnexpectedToken, p.curToken.Pos, "range template '#(...)' cannot appear at beginning of entity name")
			}
			lit, err := p.parseTemplateLiteral()
			if err != nil {
				return "", err
			}
			name += lit
		default:
			if name == "" {
				return "", NewError(ErrUnexpectedToken, p.curToken.Pos, "expected entity name")
			}
			return name, nil
		}
	}
}

func (p *Parser) parseTemplateLiteral() (string, error) {
	marker := p.curToken.Literal
	p.nextToken()
	if p.curToken.Type != TokenLParen {
		return "", NewError(ErrMissingDelimiter, p.curToken.Pos, "expected '(' after %s template marker", marker)
	}
	p.nextToken()

	// template 本体は TokenString / TokenComma のみ許可。
	body := ""
	for p.curToken.Type != TokenRParen && p.curToken.Type != TokenEOF {
		switch p.curToken.Type {
		case TokenString, TokenComma:
			body += p.curToken.Literal
		default:
			return "", NewError(ErrUnexpectedToken, p.curToken.Pos, "unexpected token %q inside template", p.curToken.Literal)
		}
		p.nextToken()
	}
	if p.curToken.Type != TokenRParen {
		return "", NewError(ErrMissingDelimiter, p.curToken.Pos, "missing ')' in template")
	}
	out := marker + "(" + body + ")"
	p.nextToken()
	return out, nil
}

// parseUntilRBrace は '}' までの entity 群を再帰的に読む。
func (p *Parser) parseUntilRBrace() ([]*Node, error) {
	nodes := []*Node{}
	for p.curToken.Type != TokenRBrace && p.curToken.Type != TokenEOF {
		parsed, err := p.parseEntity()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, parsed...)

		if p.curToken.Type == TokenComma || p.curToken.Type == TokenPipe {
			p.nextToken()
		}
	}

	if p.curToken.Type != TokenRBrace {
		return nil, NewError(ErrMissingDelimiter, p.curToken.Pos, "missing closing '}'")
	}

	p.nextToken()
	return nodes, nil
}

// cloneNodes は子ノード配列の深いコピーを作る。
func cloneNodes(in []*Node) []*Node {
	out := make([]*Node, 0, len(in))
	for _, n := range in {
		out = append(out, cloneNode(n))
	}
	return out
}

// cloneNode は単一ノードを再帰的に複製する。
func cloneNode(n *Node) *Node {
	c := &Node{Name: n.Name}
	if len(n.Children) > 0 {
		c.Children = cloneNodes(n.Children)
	}
	return c
}
