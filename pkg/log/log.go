package log

import (
	"fmt"
	"runtime"
)

func Error(err error) {
	dir, file := caller()
	fmt.Printf("Error: at %s:%d : %s\n", dir, file, err.Error())
}

func Dump(w ...interface{}) {
	dir, file := caller()
	fmt.Printf("Dump: at %s:%d\n", dir, file)
	for _, v := range w {
		fmt.Printf("%s\n", v)
	}
}

func caller() (string, int) {
	_, dir, file, _ := runtime.Caller(2)
	return dir, file
}
