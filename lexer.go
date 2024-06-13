package main

func isChar(word byte) bool {
	if word >= 65 && word < 91 || word >= 97 && word <= 122 {
		return true
	}
	return false
}

func isNumber(word byte) bool {
	if word >= 48 && word < 58 {
		return true
	}
	return false
}

func return_kind(word byte) kind {
	if isChar(word) {
		return Identifier
	}
	if isNumber(word) {
		return Number
	}
	switch word {
	case 42: // *
		return Star
	case 94: // ^
		return Raise
	case 126: // ~
		return Tilde
	case 62: // >
		return Great
	case 60: // <
		return Less
	case 61: // =
		return Equal
	case 124: // |
		return Pipe
	case 45: // -
		return Hyphen
	case 46: // .
		return Dot
	case 64: // @
		return At
	}
	return Nil
}

func GetTokens(version string) []Token {
	var tokens []Token
	i := 0
	for i < len(version) {
		ver := version[i]
		if ver != 32 { // 32 is ascii for blank
			tkind := return_kind(ver)
			if tkind == Number {
				tokens = append(tokens, Token{
					Value: string(ver),
					Kind:  tkind,
				})

				for i+1 < len(version) && (return_kind(version[i+1]) == Number || return_kind(version[i+1]) == Dot || version[i+1] == 'x') {
					tokens[len(tokens)-1].Value += string(version[i+1])
					i += 1
				}
			} else if tkind == Identifier {
				tokens = append(tokens, Token{
					Value: "",
					Kind:  tkind,
				})
				for tkind == Identifier || tkind == Number {
					tokens[len(tokens)-1].Value += string(ver)
					i += 1
					if i < len(version) {
						ver = version[i]
						tkind = return_kind(ver)
					} else {
						break
					}
				}
				i -= 1
			} else if tkind == Equal {
				if tokens[len(tokens)-1].Kind == Great {
					tokens[len(tokens)-1].Value = ">="
					tokens[len(tokens)-1].Kind = GreatEqual
				} else if tokens[len(tokens)-1].Kind == Less {
					tokens[len(tokens)-1].Value = "<="
					tokens[len(tokens)-1].Kind = LessEqual
				}
			} else if tkind == Pipe {
				if return_kind(version[i+1]) == Pipe {
					tokens = append(tokens, Token{
						Value: "||",
						Kind:  tkind,
					})
					i += 1
				}
			} else {
				tokens = append(tokens, Token{
					Value: string(ver),
					Kind:  tkind,
				})
			}
		}
		i += 1
	}
	return tokens
}
