package rider

import (
	"testing"
	"time"
	"fmt"
)

func TestRiderJwt(t *testing.T) {
	handleFunc := RiderJwt("rider", time.Hour)

}
