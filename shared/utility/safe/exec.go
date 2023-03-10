package safe

import "time"

// 很多地方函数不返回错误，但是很重要，多执行几次函数，减少函数意外报错的概率
func Exec(times int, f func(int) error) {
	for i := 0; i < times; i++ {
		err := f(i)
		if err == nil {
			return
		}

		time.Sleep(time.Duration(i+1) * time.Second)
	}
}
