package tendon

import (
	"math"
)

type EasingFunc func(current, target float64) float64

// (ratio*100)ずつ目標値に近づき、snapDictの距離で吸着
func NewLerpEasingFunc(ratio, snapDist float64) EasingFunc {
	return func(current, target float64) float64 {
		if math.Abs(current-target) < snapDist {
			return target
		}
		return current + (target-current)*ratio
	}
}

func NewLinearEasingFunc(speed float64) EasingFunc {
	return func(current, target float64) float64 {
		diff := target - current
		if math.Abs(diff) <= speed {
			return target
		}
		if diff > 0 {
			return current + speed
		}
		return current - speed
	}
}
