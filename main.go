package main

import (
	"encoding/base64"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
)

type ACME struct {
	Certificates []Certificates
}

type Certificates struct {
	Domain      Domain
	Certificate string
	Key         string
}

type Domain struct {
	Main string
}

type Resolver map[string]ACME

var sugar *zap.SugaredLogger

const (
	PathToAcmeJson = "/usr/share/tatcf/acme.json"
	PathToCertDir  = "/var/lib/tatcf"
	PathSeparator  = "/"
)

func init() {
	// initialize logger
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalln("logger could not be shut down reasonably ")
		}
	}(logger)
	sugar = logger.Sugar()

	// Check if acme.json file fileExists
	_, err := os.OpenFile(PathToAcmeJson, os.O_RDONLY, 0400)
	if !fileExists(err) {
		sugar.Fatal("no acme.json file found")
	}
	sugar.Info("acme.file found")

	// Check if dir for certs fileExists and if not create
	_, err = os.Stat(PathToCertDir)
	if !fileExists(err) {
		err = os.MkdirAll(PathToCertDir, 0644)
		if err != nil {
			sugar.Fatal("path to certificate dir cannot be created")
		}
	}
	sugar.Info("certificate dir created")
}

func main() {

	watchChan := make(chan bool, 1)
	updateChan := make(chan bool, 1)

	updateChan <- true
	for {
		select {
		case _ = <-watchChan:
			watch()
			updateChan <- true

		case _ = <-updateChan:
			updateCertificates()
			watchChan <- true
		}
	}
}

func updateCertificates() {
	sugar.Info("updating certificates")
	data, err := ioutil.ReadFile(PathToAcmeJson)
	if err != nil {
		sugar.Fatal("error when opening file: ", err)
	}

	var result Resolver

	err = json.Unmarshal(data, &result)
	if err != nil {
		sugar.Fatal("error during Unmarshal(): ", err)
	}

	for resolverName, _ := range result {
		for _, certificate := range result[resolverName].Certificates {
			processCertificate(certificate)
		}
	}

}

func watch() {
	sugar.Info("watching file for changes")
	err := watchFile(PathToAcmeJson)
	if err != nil {
		sugar.Fatal("error while watching file", err)
	}
	sugar.Info("file has been changed")

}

func processCertificate(certificate Certificates) {
	path := PathToCertDir + "/" + certificate.Domain.Main
	createDirForDomain(path, certificate.Domain.Main)

	var encoded []byte
	var fileName string

	encoded, _ = base64.StdEncoding.DecodeString(certificate.Key)
	fileName = certificate.Domain.Main + ".key"
	updateCertificateFile(path, fileName, encoded)

	encoded, _ = base64.StdEncoding.DecodeString(certificate.Certificate)
	fileName = certificate.Domain.Main + ".crt"
	updateCertificateFile(path, fileName, encoded)

}
