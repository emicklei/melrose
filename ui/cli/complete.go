package cli

import (
	"bytes"
	"sort"
	"strings"
	"unicode"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/dsl"
)

func completeMe(line string, pos int) (head string, c []string, tail string) {
	start := pos

	if pos > 0 {
		// go back to last separator
		runes := []rune(line)
		for i := start - 1; i >= 0; i-- {
			r := runes[i]
			if !unicode.IsLetter(r) && r != '_' {
				break
			}
			start = i
		}
	}
	prefix := line[start:pos]
	if len(prefix) == 0 {
		return line[0:pos], c, line[pos:]
	}
	// vars first
	// for k, _ := range melrose.Context{}.VariableStorage.Variables() {  // broken
	// 	if strings.HasPrefix(k, prefix) {
	// 		c = append(c, k[len(prefix):])
	// 	}
	// }
	for k, f := range dsl.EvalFunctions(core.PlayContext{}) {
		// TODO start from closest (
		if strings.HasPrefix(k, prefix) {
			stripped := stripParameters(f.Template)
			c = append(c, stripped[len(prefix):])
		}
	}
	sort.Strings(c)
	return line[0:pos], c, line[pos:]
}

// to strip:  ${1:param}
func stripParameters(sample string) string {
	var buf bytes.Buffer
	inparam := false
	for _, each := range sample {
		if each == '$' {
			inparam = true
			continue
		}
		if each == '}' {
			inparam = false
			continue
		}
		if !inparam {
			buf.WriteRune(each)
		}
	}
	return buf.String()
}
