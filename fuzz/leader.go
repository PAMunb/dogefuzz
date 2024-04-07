package fuzz

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

type fuzzerLeader struct {
	blackboxFuzzer             interfaces.Fuzzer
	greyboxFuzzer              interfaces.Fuzzer
	directedGreyboxFuzzer      interfaces.Fuzzer
	otherDirectedGreyboxFuzzer interfaces.Fuzzer
}

func NewFuzzerLeader(e env) *fuzzerLeader {
	return &fuzzerLeader{
		blackboxFuzzer:             e.BlackboxFuzzer(),
		greyboxFuzzer:              e.GreyboxFuzzer(),
		directedGreyboxFuzzer:      e.DirectedGreyboxFuzzer(),
		otherDirectedGreyboxFuzzer: e.OtherDirectedGreyboxFuzzer(),
	}
}

func (l *fuzzerLeader) GetFuzzerStrategy(typ common.FuzzingType) (interfaces.Fuzzer, error) {
	switch typ {
	case common.BLACKBOX_FUZZING:
		return l.blackboxFuzzer, nil
	case common.GREYBOX_FUZZING:
		return l.greyboxFuzzer, nil
	case common.DIRECTED_GREYBOX_FUZZING:
		return l.directedGreyboxFuzzer, nil
	case common.OTHER_DIRECTED_GREYBOX_FUZZING:
		return l.otherDirectedGreyboxFuzzer, nil
	default:
		return nil, ErrFuzzerTypeNotFound
	}
}
