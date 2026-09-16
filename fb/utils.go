package fb

import "strings"

func WriteDocumentation(w Writer, docs string, indent string) {
	for line := range strings.Lines(docs) {
		line = strings.TrimSpace(line)

		if indent != "" {
			w.Write(indent)
		}

		if line == "" {
			w.Write("///\n")
			continue
		}

		w.Write("/// %s\n", line)
	}
}

func WriteAttributes(w Writer, attributes []string, indent string) {
	if len(attributes) == 0 {
		return
	}

	w.Write("%s", indent)
	w.Write("#[")

	for i, attribute := range attributes {
		if i > 0 {
			w.Write(", ")
		}

		w.Write("%s", attribute)
	}

	w.Write("]\n")
}

type Receiver uint8

const (
	None Receiver = iota
	Immutable
	Mutable
)

func WriteSignature(w Writer, name string, receiver Receiver, params []*Param, returns Type) bool {
	// Signature
	if name == "" {
		w.Write("func(")
	} else {
		w.Write("func %s(", name)
	}

	switch receiver {
	case Immutable:
		w.Write("self")
	case Mutable:
		w.Write("mut self")
	default:
	}

	for i, param := range params {
		if i > 0 || receiver != None {
			w.Write(", ")
		}

		if param.Name == "" {
			if name != "" {
				panic("fb.WriteSignature() - Param name cannot be empty when the function name isn't empty")
			}
		} else {
			w.Write("%s: ", param.Name)
		}

		param.Type.Write(w)
	}

	w.Write(")")

	// Returns
	if s, ok := returns.(*SimpleType); !ok || s.Text != "void" {
		w.Write(" ")
		returns.Write(w)

		return true
	}

	return false
}
