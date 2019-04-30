package types

// Shape is a struct parsed from the data type for RTI shapes demo
type Shape struct {
	Color     string `json:"color"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Shapesize int    `json:"shapesize"`
}

// ShapeArray is a shape struct with arrays
type ShapeArray struct {
	Color     string `json:"color"`
	X         [100]int    `json:"x"`
	Y         [100]int    `json:"y"`
	Shapesize int    `json:"shapesize"`
}

// ShapeSlice is a shape struct with slices
type ShapeSlice struct {
	Color     string `json:"color"`
	X         []int    `json:"x"`
	Y         []int    `json:"y"`
	Shapesize int    `json:"shapesize"`
}
