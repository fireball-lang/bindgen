package fb

type Decl interface {
	OutputIndex_() int
	Name_() string

	Write(w Writer)
}

// Alias

type Alias struct {
	OutputIndex int

	Documentation string
	Attributes    []string

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
	WriteAttributes(w, a.Attributes, "")

	w.Write("pub type %s = ", a.Name)
	a.Type.Write(w)
	w.Write(";\n")
}

// Enum

type Case struct {
	Documentation string
	Attributes    []string

	Name  string
	Value string
}

type Enum struct {
	OutputIndex int

	Documentation string
	Attributes    []string

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
	WriteAttributes(w, e.Attributes, "")

	w.Write("pub enum %s", e.Name)

	if e.Type != nil {
		w.Write(" : ")
		e.Type.Write(w)
	}

	w.Write(" {\n")

	for _, cas := range e.Cases {
		WriteDocumentation(w, cas.Documentation, "    ")
		WriteAttributes(w, cas.Attributes, "    ")

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

type Layout uint8

const (
	C Layout = iota
	Union
	Fireball
)

type Field struct {
	Documentation string
	Attributes    []string

	Public bool
	Name   string
	Type   Type
}

type Struct struct {
	OutputIndex int

	Documentation string
	Attributes    []string

	Name   string
	Fields []*Field

	Layout Layout
}

func (s *Struct) OutputIndex_() int {
	return s.OutputIndex
}

func (s *Struct) Name_() string {
	return s.Name
}

func (s *Struct) Write(w Writer) {
	attributes := s.Attributes

	switch s.Layout {
	case C:
		attributes = append([]string{"repr(C)"}, attributes...)
	case Union:
		attributes = append([]string{"repr(Union)"}, attributes...)
	case Fireball:
	}

	WriteDocumentation(w, s.Documentation, "")
	WriteAttributes(w, attributes, "")

	if len(s.Fields) == 0 {
		w.Write("pub struct %s {}\n", s.Name)
		return
	}

	w.Write("pub struct %s {\n", s.Name)

	for _, field := range s.Fields {
		WriteDocumentation(w, field.Documentation, "    ")
		WriteAttributes(w, field.Attributes, "    ")

		if field.Public {
			w.Write("    pub %s: ", field.Name)
		} else {
			w.Write("    %s: ", field.Name)
		}

		field.Type.Write(w)
		w.Write(",\n")
	}

	w.Write("}\n")
}

// Const

type Const struct {
	OutputIndex int

	Documentation string
	Attributes    []string

	Name  string
	Type  Type
	Value string
}

func (c *Const) OutputIndex_() int {
	return c.OutputIndex
}

func (c *Const) Name_() string {
	return c.Name
}

func (c *Const) Write(w Writer) {
	WriteDocumentation(w, c.Documentation, "")
	WriteAttributes(w, c.Attributes, "")

	w.Write("pub const %s: ", c.Name)
	c.Type.Write(w)
	w.Write(" = %s;\n", c.Value)
}

// Func

type Param struct {
	Name string
	Type Type
}

type Func struct {
	OutputIndex int

	Documentation string
	Attributes    []string

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
	WriteAttributes(w, f.Attributes, "")

	// Attributes
	w.Write("#[extern")

	if f.Name == f.LinkName {
		w.Write("]\n")
	} else {
		w.Write(", link_name(\"%s\")]\n", f.LinkName)
	}

	// Signature
	w.Write("pub ")
	WriteSignature(w, f.Name, None, f.Params, f.Returns)
	w.Write(";\n")
}
