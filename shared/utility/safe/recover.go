package safe

import (
	"runtime"
	"shared/utility/glog"
)

func Recover() {
	if x := recover(); x != nil {
		glog.Errorf("recover panic: %v\n", x)
		for i := 0; i < 20; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				glog.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
	}
}
