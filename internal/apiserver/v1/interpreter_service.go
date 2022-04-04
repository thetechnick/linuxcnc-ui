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
// 	ctx context.Context,
// 	interpreter interpreter,
// 	req *linuxcncv1.InterpreterRequest,
// 	res *linuxcncv1.InterpreterResponse,
// ) error {
// 	var instructions []proto.Message

// 	sink := &rs274ngc.InterpreterSinkFuncs{
// 		MessageFn: func(msg string) {
// 			instructions = append(instructions, &linuxcncv1.MessageInstruction{
// 				Message: msg,
// 			})
// 		},
// 		CommentFn: func(comment string) {
// 			instructions = append(instructions, &linuxcncv1.CommentInstruction{
// 				Comment: comment,
// 			})
// 		},
// 	}

// 	if err := interpreter.Parse(ctx, sink, req.File); err != nil {
// 		return fmt.Errorf("parsing file: %w", err)
// 	}

// 	var err error
// 	res.Instructions = make([]*anypb.Any, len(instructions))
// 	for i := range instructions {
// 		res.Instructions[i], err = anypb.New(instructions[i])

// 	}

// 	anypb.New()
// 	return nil
// }
