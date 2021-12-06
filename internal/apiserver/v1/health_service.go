package v1

import (
	"context"

	linuxcncv1 "github.com/thetechnick/linuxcnc-ui/api/v1"
)

type HealthServiceServer struct {
	linuxcncv1.UnimplementedHealthServiceServer
}

func (h *HealthServiceServer) Health(ctx context.Context, req *linuxcncv1.HealthRequest) (*linuxcncv1.HealthResponse, error) {
	return &linuxcncv1.HealthResponse{}, nil
}
