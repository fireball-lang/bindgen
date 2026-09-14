package fb

import (
	"strings"
)

type Decl interface {
	OutputIndex_() int
	Name_() string

	Write(w Writer)
}

// Alias

type Alias struct {
	OutputIndex int

	Documentation string

	Name string
	Type Type
}

func (a *Alias) OutputIndex_() int {
	return a.OutputIndex
}

func (a *Alias) Name_() string {
	return a.Name
}

func (a *Alias) Write(w Writer) {
	WriteDocumentation(w, a.Documentation, "")

	w.Write("pub type %s = ", a.Name)
	a.Type.Write(w)
	w.Write(";\n")
}

// Enum

type Case struct {
	Documentation string

	Name  string
	Value string
}

type Enum struct {
	OutputIndex int

	Documentation string

	Name string
	Type Type

	Cases []*Case

	Bitfield bool
}

func (e *Enum) Case(name string) *Case {
	for _, cas := range e.Cases {
		if cas.Name == name {
			return cas
		}
	}

	return nil
}

func (e *Enum) OutputIndex_() int {
	return e.OutputIndex
}

func (e *Enum) Name_() string {
	return e.Name
}

func (e *Enum) Write(w Writer) {
	WriteDocumentation(w, e.Documentation, "")

	w.Write("pub enum %s", e.Name)

	if e.Type != nil {
		w.Write(" : ")
		e.Type.Write(w)
	}

	w.Write(" {\n")

	for _, cas := range e.Cases {
		WriteDocumentation(w, cas.Documentation, "    ")

		if cas.Value == "" {
			w.Write("    %s,\n", cas.Name)
		} else {
			w.Write("    %s = %s,\n", cas.Name, cas.Value)
		}
	}

	w.Write("}\n")

	// Bitfield
	if e.Bitfield {
		underlying := stringWriter{parent: w}
		e.Type.Write(&underlying)

		w.Write(`
impl %[1]s : BitNot {
    type Result = %[1]s;

    pub func bit_not(self) Self {
        return ~(*self as %[2]s) as Self;
    }
}

impl %[1]s : BitOr[%[1]s] {
    type Result = %[1]s;

    pub func bit_or(self, rhs: Self) Self {
        return (*self as %[2]s | rhs as %[2]s) as Self;
    }
}

impl %[1]s : BitAnd[%[1]s] {
    type Result = %[1]s;

    pub func bit_and(self, rhs: Self) Self {
        return (*self as %[2]s & rhs as %[2]s) as Self;
    }
}
`, e.Name, underlying.String())
	}
}

// Struct

type Field struct {
	Documentation string

	Name string
	Type Type
}

type Struct struct {
	OutputIndex int

	Documentation string

	Name   string
	Fields []*Field

	Union bool
}

func (s *Struct) OutputIndex_() int {
	return s.OutputIndex
}

func (s *Struct) Name_() string {
	return s.Name
}

func (s *Struct) Write(w Writer) {
	WriteDocumentation(w, s.Documentation, "")

	if s.Union {
		w.Write("#[repr(Union)]\n")
	} else {
		w.Write("#[repr(C)]\n")
	}

	if len(s.Fields) == 0 {
		w.Write("pub struct %s {}\n", s.Name)
		return
	}

	w.Write("pub struct %s {\n", s.Name)

	for _, field := range s.Fields {
		WriteDocumentation(w, field.Documentation, "    ")

		w.Write("    pub %s: ", field.Name)
		field.Type.Write(w)
		w.Write(",\n")
	}

	w.Write("}\n")
}

// Func

type Param struct {
	Name string
	Type Type
}

type Func struct {
	OutputIndex int

	Documentation string

	Name     string
	LinkName string

	Params  []*Param
	Returns Type

	MethodName    string
	ReceiverIndex int
}

func (f *Func) OutputIndex_() int {
	return f.OutputIndex
}

func (f *Func) Name_() string {
	return f.Name
}

func (f *Func) Write(w Writer) {
	WriteDocumentation(w, f.Documentation, "")

	// Attributes
	w.Write("#[extern")

	if f.Name == f.LinkName {
		w.Write("]\n")
	} else {
		w.Write(", link_name(\"%s\")]\n", f.LinkName)
	}

	// Signature
	w.Write("pub func %s(", f.Name)

	for i, param := range f.Params {
		if i > 0 {
			w.Write(", ")
		}

		w.Write("%s: ", param.Name)
		param.Type.Write(w)
	}

	w.Write(")")

	// Returns
	if s, ok := f.Returns.(*SimpleType); !ok || s.Text != "void" {
		w.Write(" ")
		f.Returns.Write(w)
	}

	w.Write(";\n")
}

// utils

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
