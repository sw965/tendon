package tendon

type Alignment int

const (
	AlignStart  Alignment = iota // 上端 または 左端
	AlignCenter                  // 中央
	AlignEnd                     // 下端 または 右端
)
