// Code generated by "bitstringer -type=DataSrcTLB"; DO NOT EDIT

package perffile

import "strconv"

func (i DataSrcTLB) String() string {
	if i == 0 {
		return "NA"
	}
	s := ""
	if i&DataSrcTLBHardwareWalker != 0 {
		s += "HardwareWalker|"
	}
	if i&DataSrcTLBHit != 0 {
		s += "Hit|"
	}
	if i&DataSrcTLBL1 != 0 {
		s += "L1|"
	}
	if i&DataSrcTLBL2 != 0 {
		s += "L2|"
	}
	if i&DataSrcTLBMiss != 0 {
		s += "Miss|"
	}
	if i&DataSrcTLBOSFaultHandler != 0 {
		s += "OSFaultHandler|"
	}
	i &^= 63
	if i == 0 {
		return s[:len(s)-1]
	}
	return s + "0x" + strconv.FormatUint(uint64(i), 16)
}
