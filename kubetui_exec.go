package main

import (
	"log"
	"os/exec"
	"regexp"
)

type Version struct {
	Kubernetes string
	Kubectl    string
}

func GetVersion() *Version {
	out, err := exec.Command("kubectl", "version", "--short").Output()
	if err != nil {
		log.Panic(err)
	}
	r := regexp.MustCompile(`.*: (.*)`)
	sout := string(out)
	vs := r.FindAllStringSubmatch(sout, -1)
	if vs != nil {
		return &Version{
			Kubernetes: vs[0][1],
			Kubectl:    vs[1][1],
		}
	}
	return &Version{
		Kubernetes: sout,
		Kubectl:    sout,
	}
}
