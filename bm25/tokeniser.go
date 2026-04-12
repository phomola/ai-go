package bm25

func isAlphanum(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= 128
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

// Token ...
type Token struct {
	Form       string
	IsAlphanum bool
}

// Tokenise ...
func Tokenise(text string) []Token {
	var tokens []Token
	i, l, s := 0, len(text), 0
	for i < l {
		c := text[i]
		if !isAlphanum(c) {
			if s < i {
				tokens = append(tokens, Token{text[s:i], true})
			}
			s = i + 1
			if !isWhitespace(c) {
				tokens = append(tokens, Token{text[i:s], false})
			}
		}
		i++
	}
	if s < i {
		tokens = append(tokens, Token{text[s:i], true})
	}
	return tokens
}

// GetTerms ...
func GetTerms(tokens []Token) []string {
	ts := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if t.IsAlphanum {
			ts = append(ts, t.Form)
		}
	}
	return ts
}
