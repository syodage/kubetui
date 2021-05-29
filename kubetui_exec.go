package main

import (
	"fmt"
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

// =============================Kubectl=======================================

// kubectl [command] [TYPE] [NAME] [flags]
type Kubectl struct {
	Prog    string   // kubectl
	Command string   // get, describe
	Args    []string // aka TYPE, pods, pods pod1
	Names   []string // name1, name2 or TYPE1/name1, TYPE2/name2
	Flags   []string // -n namespace, -A, --all-namespace, --
}

func NewKubectl() *Kubectl {
	return &Kubectl{
		Prog:  "kubectl",
		Args:  []string{},
		Names: []string{},
		Flags: []string{},
	}
}

func (k *Kubectl) Build() []string {
	if k.Prog == "" || k.Command == "" {
		panic(fmt.Sprintf(`Either prog:[%v] or Command:[%v] is empty`, k.Prog, k.Command))
	}

	cmd := []string{k.Prog, k.Command}
	if len(k.Args) > 0 {
		cmd = append(cmd, k.Args...)
	}

	if len(k.Names) > 0 {
		cmd = append(cmd, k.Names...)
	}

	if len(k.Flags) > 0 {
		cmd = append(cmd, k.Flags...)
	}

	return cmd
}

func (k *Kubectl) Configs(configs ...string) *Kubectl {
	k.Command = "config"
	if configs != nil {
		k.Args = append(k.Args, configs...)
	}
	return k
}

func (k *Kubectl) Get() *Kubectl {
	k.Command = "get"
	return k
}

func (k *Kubectl) Namespaces() *Kubectl {
	k.Args = append(k.Args, "namespaces")
	return k
}

func (k *Kubectl) Pods(pods ...string) *Kubectl {
	k.Args = append(k.Args, "pods")
	if pods != nil {
		k.Args = append(k.Args, pods...)
	}
	return k
}

func (k *Kubectl) Nodes(nodes ...string) *Kubectl {
	k.Args = append(k.Args, "nodes")
	if nodes != nil {
		k.Args = append(k.Args, nodes...)
	}
	return k
}

func (k *Kubectl) Deployments(deploys ...string) *Kubectl {
	k.Args = append(k.Args, "deployments")
	if deploys != nil {
		k.Args = append(k.Args, deploys...)
	}
	return k
}

func (k *Kubectl) Services(services ...string) *Kubectl {
	k.Args = append(k.Args, "services")
	if services != nil {
		k.Args = append(k.Args, services...)
	}
	return k
}

func (k *Kubectl) Endpoints(endpoints ...string) *Kubectl {
	k.Args = append(k.Args, "endpoints")
	if endpoints != nil {
		k.Args = append(k.Args, endpoints...)
	}
	return k
}

func (k *Kubectl) All() *Kubectl {
	k.Args = append(k.Args, "all")
	return k
}

func (k *Kubectl) WithAllNamespaces() *Kubectl {
	k.Args = append(k.Args, "-A")
	return k
}

func (k *Kubectl) WithNamespace(ns string) *Kubectl {
	k.Flags = append(k.Flags, "-n", ns)
	return k
}
