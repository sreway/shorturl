package b

import "os"

func f() {
	os.Exit(0) // no error; not call of os.Exit in main func
}
