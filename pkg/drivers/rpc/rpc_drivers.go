package rpc

type RpcClientDriver interface {
	GetConfigRaw() ([]byte, error)
}
