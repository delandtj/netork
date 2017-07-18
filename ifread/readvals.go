package ifread

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type IfStat struct {
	rxb, txb, rxp, txp uint64
}

func readval(path string) (uint64, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	// we assume that's always ok ?
	val := strings.Split(string(contents), "\n")
	digit, _ := strconv.ParseUint(val[0], 10, 64)
	return digit, nil
}

func readvals(i string) (IfStat, error) {
	v := IfStat{}
	var err error
	// repeat .. fsck it only 4
	if v.rxb, err = readval("/sys/class/net/" + i + "/statistics/rx_bytes"); err != nil {
		errfound = true
	}
	if v.txb, err = readval("/sys/class/net/" + i + "/statistics/tx_bytes"); err != nil {
		errfound = true
	}
	if v.rxp, err = readval("/sys/class/net/" + i + "/statistics/rx_packets"); err != nil {
		errfound = true
	}
	if v.txp, err = readval("/sys/class/net/" + i + "/statistics/tx_packets"); err != nil {
		errfound = true
	}
	if errfound {
		return nil, err
	}
	return v, nil
}
