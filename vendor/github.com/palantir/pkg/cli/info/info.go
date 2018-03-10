// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package info

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/palantir/pkg/cli"
)

type Provider func() (interface{}, error)

type Printer func(interface{})

type Info struct {
	Flag       string
	Name       string
	Value      Provider
	PrintAlone Printer // default is just %v
	PrintAmong Printer // default is like name: %v
}

func Just(value interface{}) Provider {
	return func() (interface{}, error) {
		return value, nil
	}
}

func Print(ctx cli.Context, infos []Info) error {
	switch countFlags(ctx, infos) {
	case 0:
		return PrintAll(ctx, infos)
	case 1:
		return PrintOne(ctx, infos)
	default:
		return errorMany(ctx, infos)
	}
}

func countFlags(ctx cli.Context, infos []Info) uint {
	nflags := uint(0)
	for _, info := range infos {
		if ctx.Bool(info.Flag) {
			nflags++
		}
	}
	return nflags
}

func PrintAll(ctx cli.Context, infos []Info) error {
	for _, info := range infos {
		val, err := info.Value()
		if err != nil {
			return err
		}
		if info.PrintAmong != nil {
			info.PrintAmong(val)
		} else {
			ctx.Printf("%v: %v\n", info.Name, val)
		}
	}
	return nil
}

func PrintOne(ctx cli.Context, infos []Info) error {
	for _, info := range infos {
		if ctx.Bool(info.Flag) {
			val, err := info.Value()
			if err != nil {
				return err
			}
			if info.PrintAlone != nil {
				info.PrintAlone(val)
			} else {
				ctx.Println(val)
			}
		}
	}
	return nil
}

func errorMany(ctx cli.Context, infos []Info) error {
	set := []string{}
	for _, info := range infos {
		if ctx.Bool(info.Flag) {
			set = append(set, info.Flag)
		}
	}
	if len(set) == 2 {
		return fmt.Errorf("Cannot specify both %v and %v", set[0], set[1])
	}
	return fmt.Errorf("Cannot specify more than one of these flags: %v", strings.Join(set, ", "))
}

func PrintSliceAlone(ctx cli.Context) Printer {
	return func(v interface{}) {
		slice := reflect.ValueOf(v)
		for i := 0; i < slice.Len(); i++ {
			ctx.Printf("%v\n", slice.Index(i).Interface())
		}
		if slice.Len() == 0 && ctx.IsTerminal() {
			ctx.Println("(none)")
		}
	}
}

func PrintSliceAmong(ctx cli.Context, header string) Printer {
	return func(v interface{}) {
		slice := reflect.ValueOf(v)
		strs := make([]string, 0, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			strs = append(strs, fmt.Sprintf("%v", slice.Index(i).Interface()))
		}
		if slice.Len() == 0 {
			ctx.Printf("%v: (none)\n", header)
		} else {
			ctx.Printf("%v: %v\n", header, strings.Join(strs, ", "))
		}
	}
}
