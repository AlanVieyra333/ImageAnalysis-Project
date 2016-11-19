package tools

import (
	"fmt"
	"os"
)

func RemovePath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
	} else {
		err = os.Mkdir(path, 0600)
		if err != nil {
			err = os.Mkdir(path, 0777)
			if err != nil {
				err = os.Mkdir(path, 0600)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
