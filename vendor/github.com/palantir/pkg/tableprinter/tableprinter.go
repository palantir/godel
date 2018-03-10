// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tableprinter implements a pretty printer that writes rows and columns
// as a formatted table.
package tableprinter

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

// A ColumnGetter returns the selected ColumnGetter value of a row.
type ColumnGetter func(row interface{}) string

// A Printer is a formatter for writing tables.
//
// The header of the table is written by formatting the column names. The body
// of the table is written by formatting each row using the selected columns in
// the given order.
type Printer struct {
	tw                     *tabwriter.Writer
	get                    map[string]ColumnGetter // ColumnGetters by column name
	sorted                 bool                    // whether or not to sort the body
	includeHeader          bool                    // whether or not to print the column headers
	includeHeaderSeparator bool                    // whether or not to print '-' characters under the column headers (only used if includeHeader is true)
}

// New returns a new Printer backed by the given writers. The ColumnGetters
// are used to select the columns from a row.
func New(tw *tabwriter.Writer, cgetters map[string]ColumnGetter, sorted, includeHeader, includeHeaderSeparator bool) *Printer {
	return &Printer{
		tw:                     tw,
		get:                    cgetters,
		sorted:                 sorted,
		includeHeader:          includeHeader,
		includeHeaderSeparator: includeHeaderSeparator,
	}
}

// Print writes the rows as a table, selecting on the specified column names.
func (p *Printer) Print(columns []string, rows []interface{}) error {
	if err := p.validateColumns(columns); err != nil {
		return err
	}
	if p.includeHeader {
		if err := p.printHeader(columns); err != nil {
			return err
		}
	}
	if err := p.printBody(columns, rows); err != nil {
		return err
	}
	return p.tw.Flush()
}

// validateColumns checks column names for undefined ColumnGetters.
func (p *Printer) validateColumns(columns []string) error {
	// Collect unique undefined column names
	var undefined []string
	visited := make(map[string]bool)
	for _, n := range columns {
		if _, ok := p.get[n]; !ok && !visited[n] {
			undefined = append(undefined, n)
			visited[n] = true
		}
	}

	if len(undefined) > 0 {
		var defined []string
		for g := range p.get {
			defined = append(defined, g)
		}
		sort.Strings(defined)
		return fmt.Errorf(
			"undefined column(s) %v: this tableprinter only defines column(s) %v",
			undefined, defined)
	}
	return nil
}

// printHeader writes a header for the table enumerating the columns.
func (p *Printer) printHeader(columns []string) error {
	if _, err := fmt.Fprintln(p.tw, strings.Join(columns, "\t")); err != nil {
		return fmt.Errorf("failed to write header for columns %v: %v", columns, err)
	}
	if !p.includeHeaderSeparator {
		return nil
	}
	separators := make([]string, len(columns))
	for i, v := range columns {
		separators[i] = strings.Repeat("-", len(v))
	}
	if _, err := fmt.Fprintln(p.tw, strings.Join(separators, "\t")); err != nil {
		return fmt.Errorf("failed to write header separator for columns %v: %v", columns, err)
	}
	return nil
}

// printBody writes each row by selecting the columns in the given order.
func (p *Printer) printBody(columns []string, rows []interface{}) error {
	for _, line := range p.lines(columns, rows) {
		if _, err := fmt.Fprintln(p.tw, line); err != nil {
			return fmt.Errorf("failed to write line %q: %v", line, err)
		}
	}
	return nil
}

func (p *Printer) lines(columns []string, rows []interface{}) []string {
	// Assemble one line per row
	var lines []string

	// Holds the selected values: declared outside of loop to avoid reallocation
	tokens := make([]string, len(columns))

	// Collect lines
	for _, row := range rows {
		for c, col := range columns {
			tokens[c] = p.get[col](row)
		}
		lines = append(lines, strings.Join(tokens, "\t"))
	}
	if p.sorted {
		sort.Strings(lines)
	}
	return lines
}
