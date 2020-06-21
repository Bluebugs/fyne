package internal

// PixelSize describe an object width and height in Pixels
type PixelSize struct {
	Width  int // The number of pixels along the X axis.
	Height int // The number of pixels along the Y axis.
}

// NewPixelSize return a newly allocated PixelSize of the specified dimensions.
func NewPixelSize(w int, h int) PixelSize {
	return PixelSize{w, h}
}
