package log

import (
	"fmt"
	"runtime"
)

func Error(err error) {
	_, dir, file, _ := runtime.Caller(1)
	fmt.Printf("Error: at %s:%d : %s\n", dir, file, err.Error())
}
