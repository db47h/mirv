// Code generated by "stringer -type Class elf.go"; DO NOT EDIT

package elf

import "fmt"

const _Class_name = "ClassNoneClass32Class64"

var _Class_index = [...]uint8{0, 9, 16, 23}

func (i Class) String() string {
	if i >= Class(len(_Class_index)-1) {
		return fmt.Sprintf("Class(%d)", i)
	}
	return _Class_name[_Class_index[i]:_Class_index[i+1]]
}
