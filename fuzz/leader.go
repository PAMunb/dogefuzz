package fuzz

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

type fuzzerLeader struct {
	blackboxFuzzer        interfaces.Fuzzer
	greyboxFuzzer         interfaces.Fuzzer
	directedGreyboxFuzzer interfaces.Fuzzer

	improvedGreyboxFuzzer         interfaces.Fuzzer
	customDirectedGreyboxFuzzer   interfaces.Fuzzer
	improvedDirectedGreyboxFuzzer interfaces.Fuzzer
}

func NewFuzzerLeader(e env) *fuzzerLeader {
	return &fuzzerLeader{
		blackboxFuzzer:                e.BlackboxFuzzer(),
		greyboxFuzzer:                 e.GreyboxFuzzer(),
		directedGreyboxFuzzer:         e.DirectedGreyboxFuzzer(),
		improvedGreyboxFuzzer:         e.ImprovedGreyboxFuzzer(),
		customDirectedGreyboxFuzzer:   e.CustomDirectedGreyboxFuzzer(),
		improvedDirectedGreyboxFuzzer: e.ImprovedDirectedGreyboxFuzzer(),
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

	case common.IMPROVED_GREYBOX_FUZZING:
		return l.improvedGreyboxFuzzer, nil
	case common.CUSTOM_DIRECTED_GREYBOX_FUZZING:
		return l.customDirectedGreyboxFuzzer, nil
	case common.IMPROVED_DIRECTED_GREYBOX_FUZZING:
		return l.improvedDirectedGreyboxFuzzer, nil
	default:
		return nil, ErrFuzzerTypeNotFound
	}
}
