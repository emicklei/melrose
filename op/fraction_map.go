package op

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type FractionMap struct {
	fraction core.HasValue
	target   core.Sequenceable
}

func NewFractionMap(fraction core.HasValue, target core.Sequenceable) FractionMap {
	return FractionMap{fraction: fraction, target: target}
}

func (f FractionMap) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "fractionmap(%s,%s)", core.Storex(f.fraction), core.Storex(f.target))
	return b.String()
}

// S is part of core.Sequenceable
func (f FractionMap) S() core.Sequence {
	frac := core.String(f.fraction)
	if len(frac) == 0 {
		notify.Warnf("invalid fraction type detected, %v", f.fraction)
		return core.EmptySequence
	}
	mapping, err := parseIndexFractions(frac)
	if err != nil {
		notify.Warnf("invalid fraction mapping detected, %v", err)
		return core.EmptySequence
	}
	if len(mapping) == 0 {
		return core.EmptySequence
	}
	source := f.target.S().Notes
	target := [][]core.Note{}
	for _, entry := range mapping {
		if entry.at <= 0 || entry.at > len(source) {
			// invalid offset, skip
			continue
		}
		eachGroup := source[entry.at-1] // at is one-based
		newGroup := []core.Note{}
		for _, eachNote := range eachGroup {
			newGroup = append(newGroup, eachNote.WithFraction(float32(1.0/float32(entry.inverseFraction)), entry.dotted))
		}
		target = append(target, newGroup)
	}
	return core.Sequence{Notes: target}
}

type int2fractionAndDotted struct {
	at              int
	inverseFraction int
	dotted          bool
}

// 1:1 2:.2 3:8.
func parseIndexFractions(s string) (m []int2fractionAndDotted, err error) {
	entries := strings.Fields(strings.ReplaceAll(s, ",", " "))
	for _, each := range entries {
		if !strings.Contains(each, ":") {
			return m, fmt.Errorf("mapping must use ':' as separated, got %s", each)
		}
		kv := strings.Split(each, ":")
		ik, err := strconv.Atoi(strings.TrimSpace(kv[0]))
		if err != nil {
			return m, err
		}
		rh := strings.TrimSpace(kv[1])
		dotted := strings.Contains(rh, ".")
		if dotted {
			rh = strings.Replace(rh, ".", "", -1)
		}
		iv, err := strconv.Atoi(rh)
		if err != nil {
			return m, err
		}
		if ik < 1 {
			return m, fmt.Errorf("index must be >= 1, got %d", ik)
		}
		// TODO move this to Note.ValidateFraction
		if !core.ContainsInt([]int{1, 2, 4, 8, 16}, iv) {
			return m, fmt.Errorf("fraction must be in [1,2,4,8,16], got %d", iv)
		}
		m = append(m, int2fractionAndDotted{at: ik, inverseFraction: iv, dotted: dotted})
	}
	return
}
