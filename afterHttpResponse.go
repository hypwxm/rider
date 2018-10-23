package rider

import (
	"time"
)

var afterHttpResponse = func(ctx Context, statusCode int, timeTaken time.Duration) {
	return
}
