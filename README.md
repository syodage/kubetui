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

- InfoView:
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
dlv [COMMAND]

# Command could be one of bellow
debug # debug session, you can use the following commands:
help # Prints the help message.
run # Compile, run and attached in one step
restart # Restarts the process, killing the current one if it is running.
break (break <address> [-stack <n>|-goroutine|<variable name>]*) # Set a breakpoint. Example: break foo.go:13 or break main.main.
trace # Set a tracepoint. Syntax identical to breakpoints.
continue # Run until breakpoint or program termination.
step # Single step through program.
next # Step over to next source line.
threads # Print status of all traced threads.
thread $tid # Switch to another thread.
goroutines # Print status of all goroutines.
breakpoints # Print information on all active breakpoints.
print $var # Evaluate a variable.
info $type [regex] # Outputs information about the symbol table. An optional regex filters the list. Example info funcs unicode. Valid types are:
args # Prints the name and value of all arguments to the current function
funcs # Prints the name of all defined functions
locals # Prints the name and value of all local variables in the current context
sources # Prints the path of all source files
vars # Prints the name and value of all package variables in the app. Any variable that is not local or arg is considered a package variables
regs # Prints the contents of CPU registers.
stack [ <depth> [ <goroutine id> ] ] # Prints the stacktrace of the current goroutine, up to <depth>. <depth> defaults to 10, pass a second argument to print the stacktrace of a different goroutine.
test # Compile test bindary, start and attach
exit # Exit the debugger.
```

## Resources

- [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/)
- Kubectl [cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- kubectl [commands](https://github.com/dennyzhang/cheatsheet-kubernetes-A4)
- Color string easy with [aurora](https://github.com/logrusorgru/aurora)
- Go ascii art with [go-figure](https://github.com/common-nighthawk/go-figure)
- Image to Ascii [image2ascii](https://github.com/qeesung/image2ascii)
