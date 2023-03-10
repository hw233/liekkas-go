package manager

import (
	"fmt"
	"github.com/antlabs/timer"
	"time"
)

//
//import (
//	"github.com/antlabs/timer"
//	"testing"
//	"time"
//)
//
func main() {
	newTimer := timer.NewTimer()

	go newTimer.Run()
	newTimer.ScheduleFunc(time.Second, func() {
		fmt.Println("test")
	})

}
