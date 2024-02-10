package main

import (
	"log"
	"os/exec"
	"time"
)

func setCreationTime(filepath string, time time.Time) {
	date := time.Format("2006-01-02 15:04:05")
	command := "[System.IO.File]::SetCreationTime('" + filepath + "', '" + date + "')"
	cmd := exec.Command("powershell", "-NoProfile", command)
	err := cmd.Run()
	if err != nil {
		log.Printf("Could not set CreationDate")
		log.Printf(err.Error())
	}
}
