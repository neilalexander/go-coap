package net

import (
  "github.com/yggdrasil-network/yggdrasil-go/src/config"
  "github.com/yggdrasil-network/yggdrasil-go/src/yggdrasil"
)

type YggdrasilNode struct {
  Core   *yggdrasil.Core
  Config *config.NodeConfig
}
