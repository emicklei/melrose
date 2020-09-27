package core

import (
	"bytes"
	"fmt"
	"sort"
)

type Inspection struct {
	Context    Context
	Type       string
	Text       string
	Properties map[string]interface{}
}

const maxTextLength = 40

func NewInspect(ctx Context, value interface{}) Inspection {
	i := Inspection{
		Context:    ctx,
		Type:       fmt.Sprintf("%T", value),
		Properties: map[string]interface{}{},
	}
	if s, ok := value.(Storable); ok {
		i.Text = s.Storex()
	}
	if p, ok := value.(Inspectable); ok {
		p.Inspect(i)
	}
	if s, ok := value.(Sequenceable); ok {
		s.S().Inspect(i) // show props as sequence
	}
	// default
	if len(i.Text) == 0 {
		i.Text = fmt.Sprintf("%v", value)
	}
	return i
}

// Markdown returns a markdown formatted string with inspection details
func (i Inspection) Markdown() string {
	var b bytes.Buffer
	title := i.Text
	if len(title) > maxTextLength {
		title = title[:maxTextLength] + "..."
	}
	fmt.Fprintf(&b, "`%s`\n", title)
	// sort keys
	keys := []string{}
	for k := range i.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := i.Properties[k]
		fmt.Fprintf(&b, "- %v %s\n", v, k)
	}
	return b.String()
}

func (i Inspection) String() string {
	var b bytes.Buffer
	title := i.Text
	fmt.Fprintf(&b, "%s ", title)
	// sort keys
	keys := []string{}
	for k := range i.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// last minute prop
	keys = append(keys, "type")
	i.Properties["type"] = i.Type

	for _, k := range keys {
		v := i.Properties[k]
		fmt.Fprintf(&b, "\033[94m%s:\033[0m%v ", k, v)
	}
	return b.String()
}
