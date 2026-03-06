//go:build wireinject
// +build wireinject

package registry

import (
	"context"

	"github.com/ducthangng/GeoFleet/app/handler"
	"github.com/ducthangng/GeoFleet/app/internal/infrastructure"
	"github.com/google/wire"
)

func BuildKafkaListener(ctx context.Context) (*infrastructure.KafkaConsumer, error) {
	wire.Build(KafkaListenerSet)
	return nil, nil
}

func BuildTrackingHandler(ctx context.Context) (*handler.TrackingHandler, error) {
	wire.Build(TrackingHandlerSet)
	return nil, nil
}
