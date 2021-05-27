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

## Development

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
```

## Features



## Roadmap

- ContextView:
  - Show information like, active context, active namespace, resource usage, etc...
  - Show shortcuts (optional), maybe popup window to show all the valid keybinding for the active view
- Menu:
  - Navigate up and down using keyboard.
  - Start with context which shows all the availabe contexts.

## Resources

- [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/)
- Kubectl [cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- kubectl [commands](https://github.com/dennyzhang/cheatsheet-kubernetes-A4)
- Color string easy with [aurora](https://github.com/logrusorgru/aurora)
- Go ascii art with [go-figure](https://github.com/common-nighthawk/go-figure)
- Image to Ascii [image2ascii](https://github.com/qeesung/image2ascii)
