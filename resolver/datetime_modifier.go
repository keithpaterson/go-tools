package resolver

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// supported operation types

const (
	opAdd operation = iota
	opSubtract
)

type operation int

// supported modifier fields

const (
	fieldYear field = iota
	fieldMonth
	fieldWeek
	fieldDay
	fieldHour
	fieldMinute
	fieldSecond
)

type field int

var (
	opType    = map[string]operation{"+": opAdd, "-": opSubtract}
	fieldType = map[string]field{
		"Y": fieldYear, "M": fieldMonth, "W": fieldWeek, "D": fieldDay,
		"h": fieldHour, "m": fieldMinute, "s": fieldSecond}
)

type timeValueModifier struct {
	op    operation
	field field
	delta int
}

func (m timeValueModifier) adjust(from time.Time) time.Time {
	if m.delta == 0 {
		return from
	}

	delta := m.delta
	if m.op == opSubtract {
		delta *= -1
	}

	switch m.field {
	case fieldYear:
		return from.AddDate(delta, 0, 0)
	case fieldMonth:
		return from.AddDate(0, delta, 0)
	case fieldWeek:
		return from.AddDate(0, 0, 7*delta)
	case fieldDay:
		return from.AddDate(0, 0, delta)
	case fieldHour:
		return from.Add(time.Hour * time.Duration(delta))
	case fieldMinute:
		return from.Add(time.Minute * time.Duration(delta))
	case fieldSecond:
		return from.Add(time.Second * time.Duration(delta))
	}

	// this is an error condition (how did we get here?)
	return from
}

const (
	errInvalidModifier = "invalid date/time modifier"
)

func (r *dateTimeResolver) toModifier(opval string, modval string) (timeValueModifier, error) {
	modifier := timeValueModifier{}
	var err error

	if opval == "" && modval == "" {
		// valid: delta-zero modification
		return modifier, nil
	}
	if opval == "" || modval == "" {
		// error: both must be present
		return modifier, fmt.Errorf("invalid date/time operation: ('%s', '%s')", opval, modval)
	}

	var ok bool
	var op operation
	if op, ok = opType[opval]; !ok {
		return modifier, fmt.Errorf("invalid operation '%s'", opval)
	}

	modval = strings.TrimSpace(modval)
	matches := regField.FindStringSubmatch(modval)
	if len(matches) != 3 {
		// didn't match enough information; this is an invalid specification
		return modifier, fmt.Errorf("%s '%s'", errInvalidModifier, modval)
	}

	var delta int
	if delta, err = strconv.Atoi(matches[1]); err != nil {
		return modifier, fmt.Errorf("%s delta '%s'", errInvalidModifier, matches[1])
	}

	var field field
	if field, ok = fieldType[matches[2]]; !ok {
		return modifier, fmt.Errorf("%s field'%s'", errInvalidModifier, matches[2])
	}

	modifier.op = op
	modifier.field = field
	modifier.delta = delta

	return modifier, nil
}
