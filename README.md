# Yggdrasil-Go-CoAP

[![builds.sr.ht status](https://builds.sr.ht/~fnux/yggdrasil-go-coap.svg)](https://builds.sr.ht/~fnux/yggdrasil-go-coap?)

This project is a fork of
[github.com/go-ocf/go-coap](https://github.com/go-ocf/go-coap), adding
[Yggdrasil](https://yggdrasil-network.github.io/) support to the library.

This project lives at
[git.sr.ht/~fnux/yggdrasil-go-coap](https://git.sr.ht/~fnux/yggdrasil-go-coap),
any other source (i.e. github.com) is not to be considered maintained by its
original author.

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

### Examples

A complete example can be found on
[git.sr.ht/~fnux/yggdrasil-coap-nodes](https://git.sr.ht/~fnux/yggdrasil-coap-nodes).
