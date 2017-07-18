package squeezer

import "github.com/delandtj/netork/ifread"

func main() {
	fmt.Println("vim-go")
}

const (
	_ = iota
	1000mbit
	500mbit
	200mbit
	100mbit
	50mbit
	20mbit
	5mbit
	1mbit
	200kbit
	20kbit
)

func SqueezeBW(name string, bw int) error {

}
