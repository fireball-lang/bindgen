package fb

import (
	"fmt"
	"strings"
)

type Writer interface {
	Write(format string, args ...any)

	IsOutputCurrent(index int) bool
	GetOutputModule(index int) string
}

type stringWriter struct {
	parent Writer
	sb     strings.Builder
}

func (s *stringWriter) Write(format string, args ...any) {
	_, _ = fmt.Fprintf(&s.sb, format, args...)
}

func (s *stringWriter) IsOutputCurrent(index int) bool {
	return s.parent.IsOutputCurrent(index)
}

func (s *stringWriter) GetOutputModule(index int) string {
	return s.parent.GetOutputModule(index)
}

func (s *stringWriter) String() string {
	return s.sb.String()
}
