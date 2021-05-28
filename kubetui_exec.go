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
	arr := []string{"kubectl", "version", "--short"}
	sout := executeCmd(arr)
	r := regexp.MustCompile(`.*: (.*)`)
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

func GetCurrentNamespace() string {
	arr := []string{"kubens", "-c"}
	return executeCmd(arr)
}

func GetCurrentContext() string {
	arr := []string{"kubectl", "config", "current-context"}
	return executeCmd(arr)
}

func GetClusterName() string {
	return "Minikube"
}

func executeCmd(arr []string) string {
	out, err := exec.Command(arr[0], arr[1:]...).Output()
	if err != nil {
		log.Panic(err)
	}
	return string(out)
}
