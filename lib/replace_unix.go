// +build linux

package lib

import (
	"os"
)

func replaceMe(oldf, newf string) error {
	return os.Rename(oldf, newf)
}
