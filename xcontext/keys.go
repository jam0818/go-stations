package xcontext

import (
	"context"
	"github.com/TechBowl-japan/go-stations/model"
)

type envinfoKey struct{}

func SetOSInfo(ctx context.Context, info model.EnvInfo) context.Context {
	return context.WithValue(ctx, envinfoKey{}, info)
}

func GetOSInfo(ctx context.Context) model.EnvInfo {
	envInfo, ok := ctx.Value(model.EnvinfoKey{}).(model.EnvInfo)
	if !ok {
		envInfo = model.EnvInfo{
			OS:      "",
			Browser: "",
		}
	}
	return envInfo
}
