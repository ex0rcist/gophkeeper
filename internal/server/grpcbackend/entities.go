package grpcbackend

import (
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/entities"
)

type GRPCServerAddress entities.Address

type GRPCServerAddressDependencies struct {
	config.Dependency
}

func NewGRPCServerAddress(deps GRPCServerAddressDependencies) GRPCServerAddress {
	return GRPCServerAddress(deps.Config.Address)
}
