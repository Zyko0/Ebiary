//go:build js

package clipboard

import (
	"errors"
	"sync"
	"syscall/js"
)

// implementation from here: https://github.com/atotto/clipboard/pull/48/files

func readAll() (string, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	var result string
	var err error
	promise := js.Global().Get("navigator").Get("clipboard").Call("readText").Call("then", js.FuncOf(func(me js.Value, args []js.Value) interface{} {
		result = args[0].String()
		wg.Done()
		return nil
	}), js.FuncOf(func(me js.Value, args []js.Value) interface{} {
		err = errors.New(args[0].String())
		wg.Done()
		return nil
	}))

	if !promise.Truthy() {
		return "", errors.New("no promise received by js")
	}

	// Wait for promise to resolve
	wg.Wait()

	// Return value
	return result, err
}

func ReadText() string {
	str, err := readAll()
	if err != nil {
		return ""
	}
	return str
}

func ReadImage() {
	// TODO:
}

func writeAll(text string) error {
	var wg sync.WaitGroup
	wg.Add(1)
	var err error
	promise := js.Global().Get("navigator").Get("clipboard").Call("writeText", text).Call("then", js.FuncOf(func(me js.Value, args []js.Value) interface{} {
		wg.Done()
		return nil
	}), js.FuncOf(func(me js.Value, args []js.Value) interface{} {
		err = errors.New(args[0].String())
		wg.Done()
		return nil
	}))

	if !promise.Truthy() {
		return errors.New("No promise received by JS")
	}

	// Wait for promise to resolve
	wg.Wait()

	// Return value
	return err
}

func WriteText(txt string) {
	writeAll(txt)
}