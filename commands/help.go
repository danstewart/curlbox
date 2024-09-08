package commands

import (
	"fmt"
)

func Help() {
	fmt.Println(`curlbox - A curl management tool

Create a new curlbox at the given path:

    curlbox create path/to/curlbox

Run a script, loading any vars in the curlbox tree:

    curlbox run path/to/script`)
}
