// Code generated by "stringer -type Data elf.go"; DO NOT EDIT

package elf

import "fmt"

const _Data_name = "DataNoneDataLittleDataBig"

var _Data_index = [...]uint8{0, 8, 18, 25}

func (i Data) String() string {
	if i >= Data(len(_Data_index)-1) {
		return fmt.Sprintf("Data(%d)", i)
	}
	return _Data_name[_Data_index[i]:_Data_index[i+1]]
}
