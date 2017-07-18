package ifread

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/VividCortex/ewma"
)

type NicIO struct {
	name  string
	rates ifaceRates
	ewma  ifaceEWMA
	sync.RWMutex
}
type ifaceRates struct {
	rxb, txb, rxp, txp func(uint64) uint64
}

type ifaceEWMA struct {
	rxb, txb, rxp, txp ewma.MovingAverage
}

// Delta is a small closure over the counters, returning the delta against previous
// first = initial value
func Delta(first uint64) func(uint64) uint64 {
	keep := first
	return func(delta uint64) uint64 {
		v := delta - keep
		keep = delta
		return v
	}
}

// MakeNic initalises the NicIO type
// It'll be used in a receiver go channel
func MakeNic(name string) (*NicIO, error) {
	var nic NicIO
	nic.name = name
	r, err := readvals(name)
	if err != nil {
		return nil, err
	}
	nic.Lock()
	defer nic.Unlock()
	nic.rates.rxb = Delta(r.rxb)
	nic.rates.txb = Delta(r.txb)
	nic.rates.rxp = Delta(r.rxp)
	nic.rates.txp = Delta(r.txp)
	nic.ewma.rxb = ewma.NewMovingAverage()
	nic.ewma.txb = ewma.NewMovingAverage()
	nic.ewma.rxp = ewma.NewMovingAverage()
	nic.ewma.txp = ewma.NewMovingAverage()
	return &nic, nil
}

// ReadNic put all NIC data in struct
func (nic *NicIO) ReadNic() error {
	ifstat, err := readvals(nic.name)
	// fmt.Println("%+v", ifstat)
	if err == nil {
		// mind that calling nic.rate.xxx updates rate
		nic.ewma.rxb.Add(float64(nic.rates.rxb(ifstat.rxb)))
		nic.ewma.txb.Add(float64(nic.rates.txb(ifstat.txb)))
		nic.ewma.rxp.Add(float64(nic.rates.rxp(ifstat.rxp)))
		nic.ewma.txp.Add(float64(nic.rates.txp(ifstat.txp)))
		return nil
	}
	return err
}

// ListNics retuns the list of nics
// TODO: make it regexp aware
func ListNics() []string {
	var ifaces []string
	l, err := ioutil.ReadDir("/sys/class/net")
	if err != nil {
		errors.New("Can't read /sys/class/net. Is /sys mounted? Bailing...")
		os.Exit(1)
	}
	for _, iface := range l {
		ifaces = append(ifaces, iface.Name())
	}
	return ifaces
}
