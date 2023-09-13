// +build windows

package lib

import (
	"os"
)

func replaceMe(oldf, newf string) error {
	err := os.Rename(newf, newf+"~")
	if err != nil {
		return err
	}
	return os.Rename(oldf, newf)
}
