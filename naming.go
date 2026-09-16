package bindgen

import (
	"strings"
	"unicode"

	"github.com/fireball-lang/bindgen/fb"
)

func CamelToSnakeCase(str string) string {
	runes := []rune(str)

	var sb strings.Builder
	last := rune(0)

	write := func(ch rune) {
		sb.WriteRune(ch)
		last = ch
	}

	for i, ch := range runes {
		// The D of a 1D/2D/3D suffix belongs to the digit
		if i > 0 && (runes[i-1] == '1' || runes[i-1] == '2' || runes[i-1] == '3') && ch == 'D' &&
			(i+1 == len(runes) || !unicode.IsLower(runes[i+1])) {
			write('d')
			continue
		}

		// 1D/2D/3D suffixes get their underscore before the digit
		if (ch == '1' || ch == '2' || ch == '3') && i+1 < len(runes) && runes[i+1] == 'D' &&
			(i+2 == len(runes) || !unicode.IsLower(runes[i+2])) {
			if last != '_' {
				write('_')
			}

			write(ch)
			continue
		}

		if !unicode.IsUpper(ch) {
			write(ch)
			continue
		}

		if i > 0 {
			prev := runes[i-1]

			if !unicode.IsUpper(prev) || i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				if last != '_' {
					write('_')
				}
			}
		}

		write(unicode.ToLower(ch))
	}

	return sb.String()
}

func SnakeToPascalCase(str string) string {
	var sb strings.Builder
	nextUpper := true

	for _, ch := range str {
		if ch == '_' {
			nextUpper = true
			continue
		}

		if nextUpper {
			sb.WriteRune(unicode.ToUpper(ch))
			nextUpper = false
		} else {
			sb.WriteRune(unicode.ToLower(ch))
		}

		if unicode.IsDigit(ch) {
			nextUpper = true
		}
	}

	return sb.String()
}

func GetDeclOrder(decl fb.Decl) int {
	switch decl.(type) {
	case *fb.Alias:
		return 1
	case *fb.Enum:
		return 2
	case *fb.Struct:
		return 3
	case *fb.Func:
		return 4

	default:
		panic("bindgen.GetDeclOrder() - Invalid declaration")
	}
}
