package rider

type Error struct {
	Error      string
	StatusCode int
	StatusText string
	Stack      string
}
