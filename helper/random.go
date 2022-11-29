package helper

import "math/rand"

// 取随机数
func RangeRand(min, max int) int {
	if min > max {
		panic("the min is greater than max!")
	}
	dif := max - min
	return min + rand.Intn(dif)
}
