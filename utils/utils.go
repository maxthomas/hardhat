package utils

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
)

func LoadTarFile(tarRdr *tar.Reader, data chan<- []byte) error {
	var err error
	for {
		_, err = tarRdr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error reading tar file: %v", err)
		}

		dataBytes, err := ioutil.ReadAll(tarRdr)
		if err != nil {
			return err
		}
		data <- dataBytes
	}
	return nil
}
