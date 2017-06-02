package common

import (
	"fmt"
	"runtime"
)

const VerBinary = "0.0.4"

func VerString(app string) string {
	return fmt.Sprintf("%s v%s (built w/%s)", app, VerBinary, runtime.Version())
}
