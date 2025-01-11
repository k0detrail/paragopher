package utils

func BoundingBoxOverlap(x1, y1, w1, h1 float32, x2, y2, w2, h2 float32) bool {
	if x1+w1 < x2 || x2+w2 < x1 {
		return false
	}
	if y1+h1 < y2 || y2+h2 < y1 {
		return false
	}
	return true
}

func HorizontalOverlap(x1, w1, x2, w2 float32) bool {
	left1 := x1 - w1/2
	right1 := x1 + w1/2
	left2 := x2 - w2/2
	right2 := x2 + w2/2
	return !(right1 < left2 || left1 > right2)
}
