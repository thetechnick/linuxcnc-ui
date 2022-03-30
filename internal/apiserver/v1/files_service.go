package v1

import (
	"context"
	"fmt"
	"os"
	"path"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	linuxcncv1 "github.com/thetechnick/linuxcnc-ui/api/v1"
)

type FilesServiceServer struct {
	linuxcncv1.UnimplementedFilesServiceServer

	Root string // root of the folder structure to expose
}

func (s *FilesServiceServer) List(
	ctx context.Context, req *linuxcncv1.FileListRequest,
) (*linuxcncv1.FileListResponse, error) {

	res := &linuxcncv1.FileListResponse{}
	if err := listFiles(path.Join(s.Root, req.Path), res); err != nil {
		// log
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func listFiles(path string, res *linuxcncv1.FileListResponse) error {
	folder, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file path: %w", err)
	}
	defer folder.Close()

	entries, err := folder.ReadDir(-1)
	if err != nil {
		return fmt.Errorf("reading directory content: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			res.Folders = append(res.Folders, &linuxcncv1.Folder{
				Name: entry.Name(),
			})
		} else {
			res.Files = append(res.Files, &linuxcncv1.File{
				Name: entry.Name(),
			})
		}
	}

	return nil
}
