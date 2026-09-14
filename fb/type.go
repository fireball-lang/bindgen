package fb

type Type interface {
	isType()

	Write(w Writer)
}

// Simple

type SimpleType struct {
	Text string
}

func (s *SimpleType) isType() {}

func (s *SimpleType) Write(w Writer) {
	w.Write(s.Text)
}

// DeclType

type DeclType struct {
	Decl Decl
}

func (d *DeclType) isType() {}

func (d *DeclType) Write(w Writer) {
	if !w.IsOutputCurrent(d.Decl.OutputIndex_()) {
		w.Write(w.GetOutputModule(d.Decl.OutputIndex_()))
	}

	w.Write(d.Decl.Name_())
}

// ArrayType

type ArrayType struct {
	Size    uint32
	Element Type
}

func (a *ArrayType) isType() {}

func (a *ArrayType) Write(w Writer) {
	w.Write("[%d]", a.Size)
	a.Element.Write(w)
}

// PointerType

type PointerType struct {
	Mutable bool
	Pointee Type
}

func (p *PointerType) isType() {}

func (p *PointerType) Write(w Writer) {
	if p.Mutable {
		w.Write("mut ")
	}

	w.Write("*")
	p.Pointee.Write(w)
}

// FuncType

type FuncType struct {
	Params  []Param
	Returns Type
}

func (f *FuncType) isType() {}

func (f *FuncType) Write(w Writer) {
	w.Write("func(")

	for i, param := range f.Params {
		if i > 0 {
			w.Write(", ")
		}

		if param.Name != "" {
			w.Write("%s: ", param.Name)
		}

		param.Type.Write(w)
	}

	w.Write(")")

	if s, ok := f.Returns.(*SimpleType); !ok || s.Text != "void" {
		w.Write(" ")
		f.Returns.Write(w)
	}
}
