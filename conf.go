package workerpool

import "time"

// Conf workerpool configure
type Conf struct {
	RTimes   int8              // retry times
	Interval time.Duration     // retry interval
	PanicHD  func(interface{}) // panic handler
	Retry    func(error) bool  // whether retry +1, true retry, false return error
}
