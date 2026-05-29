package concurrent

import (
	"fmt"
	"net/http"
	"runtime"
	"testing"
	"unsafe"
)

type Reader interface {
	read()
}

type IOReader struct{}

func (r *IOReader) read() {}

func WithIOReader() Reader {
	return &IOReader{}
}

type FileReader struct{}

func (r FileReader) read() {}

func WithFileReader() Reader {
	return FileReader{}
}

func TestXContext(t *testing.T) {
	fmt.Println(unsafe.Sizeof(IOReader{}))
	fmt.Println(unsafe.Sizeof(&IOReader{}))
	ioReader := WithIOReader()
	fmt.Println(unsafe.Sizeof(ioReader))
	fileReader := WithFileReader()
	fmt.Println(unsafe.Sizeof(fileReader))
}

func TestContext(t *testing.T) {
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.NumGoroutine())
	cli := http.Client{}
	fmt.Println(runtime.NumGoroutine())
	_, _ = cli.Get("https://www.baidu.com")
	fmt.Println(runtime.NumGoroutine())
}
