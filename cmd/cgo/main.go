package main

/*
// 告诉 Go 链接器链接 Mac 自带的 zlib 库
#cgo LDFLAGS: -lz
#include <zlib.h>
*/
import "C"
import "fmt"

func main() {
	// 调用 zlib 库中的 zlibVersion 函数
	// 该函数返回一个 C 语言的 char* 指针
	cVer := C.zlibVersion()

	// 将 C 的 char* 转换为 Go 的 string
	goVer := C.GoString(cVer)

	fmt.Printf("🎉 成功调用 Mac 内置 zlib！版本号为: %s\n", goVer)
}
