// Code generated by "stringer -type busOp ."; DO NOT EDIT

package mem

import "fmt"

const _busOp_name = "opReadopWrite"

var _busOp_index = [...]uint8{0, 6, 13}

func (i busOp) String() string {
	if i < 0 || i >= busOp(len(_busOp_index)-1) {
		return fmt.Sprintf("busOp(%d)", i)
	}
	return _busOp_name[_busOp_index[i]:_busOp_index[i+1]]
}
