package v1

import (
	"context"

	linuxcncv1 "github.com/thetechnick/linuxcnc-ui/api/v1"
	"github.com/thetechnick/linuxcnc-ui/internal/rs274ngc"
)

type InterpreterServiceServer struct {
	linuxcncv1.UnimplementedInterpreterServiceServer

	Interpreter interpreter
}

type interpreter interface {
	Parse(
		ctx context.Context,
		sink rs274ngc.InterpreterSink, filename string,
	) error
}

func (s *InterpreterServiceServer) Interpret(ctx context.Context, req *linuxcncv1.InterpreterRequest) (*linuxcncv1.InterpreterResponse, error) {
	return nil, nil
}

// func interpret(
// 	req *linuxcncv1.InterpreterRequest,
// 	res *linuxcncv1.InterpreterResponse,
// ) error {

// 	anypb.New()

// 	return nil
// }
