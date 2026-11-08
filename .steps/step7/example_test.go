// This file is in package avatar_test, not package avatar. It sees only the
// exported names, the same as any other package that imports avatar.
package avatar_test

import (
	"fmt"

	"arena/avatar"
)

// go test runs an Example and compares what it prints with the Output comment.
func ExampleDirection() {
	dx, dy := avatar.Direction("w")
	fmt.Println(dx, dy)
	// Output: 0 -1
}
