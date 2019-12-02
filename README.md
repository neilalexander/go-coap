# Yggdrasil-Go-CoAP

[![builds.sr.ht status](https://builds.sr.ht/~fnux/yggdrasil-go-coap.svg)](https://builds.sr.ht/~fnux/yggdrasil-go-coap?)

This project is a fork of
[github.com/go-ocf/go-coap](https://github.com/go-ocf/go-coap), adding
[Yggdrasil](https://yggdrasil-network.github.io/) support to the library.

This project lives at
[git.sr.ht/~fnux/yggdrasil-go-coap](https://git.sr.ht/~fnux/yggdrasil-go-coap),
a bug tracker being available at
[todo.sr.ht/~fnux/yggdrasil-go-coap](https://todo.sr.ht/~fnux/yggdrasil-go-coap).

## Contributions

You can send [me](https://fnux.ch/) [patches](https://git-send-email.io/) by
email. Feel free to directly contact [me](https://fnux.ch/) if you have issues
with this workflow.

## Yggdrasil support

In order to support Yggdrasil, the following structure and methods are
provided:

```
type YggdrasilNode struct {
	Core   *yggdrasil.Core
	Config *config.NodeConfig
}


func DialYggdrasil(node coapNet.YggdrasilNode, address string) (*ClientConn, error)
func ListenAndServeYggdrasil(node coapNet.YggdrasilNode, handler Handler) error
```

A complete example can be found on
[git.sr.ht/~fnux/yggdrasil-coap-nodes](https://git.sr.ht/~fnux/yggdrasil-coap-nodes).

## Custom CoAP Options

This fork also allow to use custom [CoAP
options](https://tools.ietf.org/html/rfc7252#section-3.1):

```
func SetOptionDef(oid OptionID, format string, minLen int, maxLen int)

// Example
var MyCustomCoAPOptionID coap.OptionID = 65000
coap.SetOptionDef(MyCustomCoAPOptionID, "string", 0, 255)

params := coap.MessageParams{ ... }
msg := c.NewMessage(params)

msg.SetOption(MyCustomCoAPOption, "spouik")
```
