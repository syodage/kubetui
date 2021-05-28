# KubeTUI

```t
 _   __      ___    _____        
| | / /     / _ \  (_   _)       
| |/ /_   _| |_) )___| |_   _ _  
|   <| | | |  _ </ __) | | | | | 
| |\ \ |_| | |_) > _)| | |_| | | 
|_| \_\___/|  __/\___)_|\___/ \_)
           | |                   
           |_|                   
                           v0.0.1
```

Terminal based UI to monitor Kubernete resources. KuebTUI use vim like keybinding to navigate and other operations. Easy learning curve for vim users.

KubeTUI uses [tview](https://github.com/rivo/tview) a Terminal UI library with rich, interactive widgets written in Go.

## Prerequisites

Install `kubectl`, `kubens` and `kubectx` CLI tools. These tools are used to connect with kubernetes api server. Basic instructions are explained below for more details chekcout the [official](https://kubernetes.io/docs/tasks/tools/) tools documentation. For `kubectx` follow [this](https://github.com/ahmetb/kubectx#installation)

### MacOs

```sh
# Install kubectl 
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
 
# or with homebrew
brew install kubectl

# Install kubectx which install kubens too
brew install kubectx
```

## Features

{TODO}

## Roadmap

- ContextView:
  - Show information like, active context, active namespace, resource usage, etc...
  - Show shortcuts (optional), maybe popup window to show all the valid keybinding for the active view
- Menu:
  - Navigate up and down using keyboard.
  - Start with context which shows all the availabe contexts.

## Development

### Run

Here are few commands that useful when developing the application.

```sh
# format code
go fmt ./...

# start terminal
go run main.go

# build the app
go build

# create vendor directory
go vendor

# Usefull alias
alias gf='go fmt ./...'
alias gr='go run ./'
```

### How `Main` view get update

- `NewKubertui(...)` start goroutine to listen for `KEvent`s
- `Menu` view fire `KEvent`s when user select menuItem by pressing enter(or space).
- When ever there is a new `KEvent`, `Kubetui` call `main.HandleStateChange(...)` with the event.
- `main.HandleStateChange(...)` call kubernetes api as needed and update the main view with new data.

### Debug

[Delve](https://github.com/go-delve/delve) is the easiest way to debug a go appliations. Far better than `fmt.Print(...)` statments.

#### Install Delve

```sh
# set Go env paths
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:PATH"

# on Go version 1.16 or later
go get github.com/go-delve/delve/cmd/dlv@latest

# Mac users: 
# Should install command line developer tools
xcode-select --install
# To passthrough authoization pop up
sudo /usr/sbin/DevToolsSecurity --enable
# might need to add user to the developer group
sudo dscl . append /Groups/_developer GroupMembership $(whoami)
# End Mac users
```

#### Debug with Delve

```sh
# start debug session, project path is optional if you are in the root
dlv debug
# To set break point at main function, `b` is alias for break 
break main.main 
# Add break point to a line number
break main.main:20
# Continute to the next debug point
continue
# Use print command to explore data
print <variable> or <expression>
# simply use 'r' command to restart the session
```

## Resources

- [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/)
- Kubectl [cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- kubectl [commands](https://github.com/dennyzhang/cheatsheet-kubernetes-A4)
- Color string easy with [aurora](https://github.com/logrusorgru/aurora)
- Go ascii art with [go-figure](https://github.com/common-nighthawk/go-figure)
- Image to Ascii [image2ascii](https://github.com/qeesung/image2ascii)
