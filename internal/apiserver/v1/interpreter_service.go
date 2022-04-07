package v1

import (
	"context"
	"fmt"
	"path"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	linuxcncv1 "github.com/thetechnick/linuxcnc-ui/api/v1"
	"github.com/thetechnick/linuxcnc-ui/internal/rs274ngc"
)

type InterpreterServiceServer struct {
	linuxcncv1.UnimplementedInterpreterServiceServer

	Root        string // root of the folder structure to expose
	Interpreter interpreter
}

type interpreter interface {
	Parse(
		ctx context.Context,
		sink rs274ngc.InterpreterSink, filename string,
	) error
}

func (s *InterpreterServiceServer) Interpret(ctx context.Context, req *linuxcncv1.InterpreterRequest) (*linuxcncv1.InterpreterResponse, error) {
	filePath := path.Join(s.Root, req.File)
	res := &linuxcncv1.InterpreterResponse{}

	return res, interpret(ctx, s.Interpreter, filePath, res)
}

func interpret(
	ctx context.Context,
	interpreter interpreter,
	filePath string,
	res *linuxcncv1.InterpreterResponse,
) error {
	var instructions []proto.Message

	sink := &rs274ngc.InterpreterSinkFuncs{
		MessageFn: func(msg string) {
			instructions = append(instructions, &linuxcncv1.MessageInstruction{
				Message: msg,
			})
		},
		CommentFn: func(comment string) {
			instructions = append(instructions, &linuxcncv1.CommentInstruction{
				Comment: comment,
			})
		},
		ChangeToolFn: func(toolNumber int) {
			instructions = append(instructions, &linuxcncv1.ChangeToolInstruction{
				ToolNumber: int32(toolNumber),
			})
		},

		// Settings
		UseLengthUnitFn: func(unit rs274ngc.LengthUnit) {
			instructions = append(instructions, &linuxcncv1.UseLengthUnitInstruction{
				Unit: string(unit),
			})
		},
		UseToolLengthOffsetFn: func(pos rs274ngc.Position) {
			instructions = append(instructions, &linuxcncv1.UseToolLengthOffsetInstruction{
				Position: toV1Position(pos),
			})
		},
		SelectPlaneFn: func(plane rs274ngc.Plane) {
			instructions = append(instructions, &linuxcncv1.SelectPlaneInstruction{
				Plane: string(plane),
			})
		},
		SetXYRotationFn: func(position float64) {
			instructions = append(instructions, &linuxcncv1.SetXYRotationInstruction{
				Position: position,
			})
		},
		SetG5XOffsetFn: func(index int, pos rs274ngc.Position) {
			instructions = append(instructions, &linuxcncv1.SetG5XOffsetInstruction{
				Index:    int32(index),
				Position: toV1Position(pos),
			})
		},
		SetG92OffsetFn: func(pos rs274ngc.Position) {
			instructions = append(instructions, &linuxcncv1.SetG92OffsetInstruction{
				Position: toV1Position(pos),
			})
		},
		SetTraverseRateFn: func(rate float64) {
			instructions = append(instructions, &linuxcncv1.SetTraverseRateInstruction{
				Rate: rate,
			})
		},
		SetFeedModeFn: func(spindleNumber int, mode int) {
			instructions = append(instructions, &linuxcncv1.SetFeedModeInstruction{
				SpindleNumber: int32(spindleNumber),
				Mode:          int32(mode),
			})
		},
		SetFeedRateFn: func(rate float64) {
			instructions = append(instructions, &linuxcncv1.SetFeedRateInstruction{
				Rate: rate,
			})
		},

		// Movement
		StraightTraverseFn: func(lineNumber int, pos rs274ngc.Position) {
			instructions = append(instructions, &linuxcncv1.StraightTraverseInstruction{
				LineNumber: int32(lineNumber),
				Position:   toV1Position(pos),
			})
		},
		ArcFeedFn: func(m rs274ngc.ArcMove) {
			instructions = append(instructions, &linuxcncv1.ArcFeedInstruction{
				ArcMove: &linuxcncv1.ArcMove{
					FirstEnd:     m.FirstEnd,
					SecondEnd:    m.SecondEnd,
					FirstAxis:    m.FirstAxis,
					SecondAxis:   m.SecondAxis,
					Rotation:     int32(m.Rotation),
					AxisEndPoint: m.AxisEndPoint,

					A: m.A,
					B: m.B,
					C: m.C,
					U: m.A,
					V: m.B,
					W: m.C,
				},
			})
		},
		StraightFeedFn: func(lineNumber int, pos rs274ngc.Position) {
			instructions = append(instructions, &linuxcncv1.StraightFeedInstruction{
				LineNumber: int32(lineNumber),
				Position:   toV1Position(pos),
			})
		},
		StraightProbeFn: func(lineNumber int, pos rs274ngc.Position) {
			instructions = append(instructions, &linuxcncv1.StraightProbeInstruction{
				LineNumber: int32(lineNumber),
				Position:   toV1Position(pos),
			})
		},
		RigidTapFn: func(lineNumber int, X, Y, Z, scale float64) {
			instructions = append(instructions, &linuxcncv1.RigidTapInstruction{
				LineNumber: int32(lineNumber),
				X:          X,
				Y:          Y,
				Z:          Z,
				Scale:      scale,
			})
		},
		DwellFn: func(d time.Duration) {
			instructions = append(instructions, &linuxcncv1.DwellInstruction{
				Duration: durationpb.New(d),
			})
		},
	}

	if err := interpreter.Parse(ctx, sink, filePath); err != nil {
		return fmt.Errorf("parsing file: %w", err)
	}

	res.Instructions = make([]*anypb.Any, len(instructions))
	for i := range instructions {
		var err error
		res.Instructions[i], err = anypb.New(instructions[i])
		if err != nil {
			return fmt.Errorf("converting to anypb: %w", err)
		}
	}
	return nil
}

func toV1Position(pos rs274ngc.Position) *linuxcncv1.InterpreterPosition {
	return &linuxcncv1.InterpreterPosition{
		X: pos.X,
		Y: pos.Y,
		Z: pos.Z,
		A: pos.A,
		B: pos.B,
		C: pos.C,
		U: pos.U,
		V: pos.V,
		W: pos.W,
	}
}
