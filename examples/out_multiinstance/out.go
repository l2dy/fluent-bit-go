package main

import (
	"C"
	"fmt"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/fluent/fluent-bit-go/output/log"
)

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(def, "multiinstance", "Testing multiple instances.")
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	id := output.FLBPluginConfigKey(plugin, "id")
	log.PErrorf(plugin, "[multiinstance] id = %q", id)
	// Set the context to point to any Go variable
	output.FLBPluginSetContext(plugin, id)

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	log.PError("[multiinstance] Flush called for unknown instance")
	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	// Type assert context back into the original type for the Go variable
	id := output.FLBPluginGetContext(ctx).(string)
	log.PErrorf("[multiinstance] Flush called for id: %s", id)

	dec := output.NewDecoder(data, int(length))

	count := 0
	for {
		ret, ts, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		// Print record keys and values
		timestamp := ts.(output.FLBTime)
		fmt.Printf("[%d] %s: [%s, {", count, C.GoString(tag), timestamp.String())

		for k, v := range record {
			fmt.Printf("\"%s\": %v, ", k, v)
		}
		fmt.Printf("}\n")
		count++
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func main() {
}
