// Code generated by "bitstringer -type=DataSrcOp"; DO NOT EDIT

package perffile

import "strconv"

func (i DataSrcOp) String() string {
	if i == 0 {
		return "NA"
	}
	s := ""
	if i&DataSrcOpExec != 0 {
		s += "Exec|"
	}
	if i&DataSrcOpLoad != 0 {
		s += "Load|"
	}
	if i&DataSrcOpPrefetch != 0 {
		s += "Prefetch|"
	}
	if i&DataSrcOpStore != 0 {
		s += "Store|"
	}
	i &^= 15
	if i == 0 {
		return s[:len(s)-1]
	}
	return s + "0x" + strconv.FormatUint(uint64(i), 16)
}
