package main

import (
	"bytes"
	"errors"
	"os"
	"time"
)

func fileExists(err error) bool {
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func updateCertificateFile(path string, fileName string, encoded []byte) {
	file, err := os.OpenFile(path+PathSeparator+fileName, os.O_RDWR, 0644)
	if !fileExists(err) {
		file, err = os.OpenFile(path+PathSeparator+fileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			sugar.Fatal("could not create file: ", fileName)
		}
	} else if err != nil {
		sugar.Fatal("could not open file: ", fileName)
	}
	//compare file
	filestat, err := file.Stat()
	if err != nil {
		sugar.Fatal("could not get fileInfo from file: ", fileName)
	}
	currentCertificate := make([]byte, filestat.Size())
	_, err = file.Read(currentCertificate)
	if !bytes.Equal(currentCertificate, encoded) {
		// overwrite file
		_, err = file.Write(encoded)
		if err != nil {
			sugar.Fatal("could not write to file: ", fileName)
		}
		sugar.Info("updated file: ", fileName)
	}

	err = file.Close()
	if err != nil {
		sugar.Fatal("error while closing file: ", fileName)
	}
}

func createDirForDomain(path, domain string) {
	// Check if dir for domain fileExists
	_, err := os.Stat(path)
	if !fileExists(err) {
		err = os.Mkdir(path, 0744)
		if err != nil {
			sugar.Fatal("error while creating dir for domain: ", domain, err)
		}
	} else if err != nil {
		sugar.Fatal("error while checking if dir for domain fileExists", domain)
	}
}
