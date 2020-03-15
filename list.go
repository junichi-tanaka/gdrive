package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var (
	driveService  *drive.Service
	driveFolderId string
)

const (
	// https://developers.google.com/drive/api/v3/search-files
	driveQuery = "'%s' in parents and mimeType = 'text/csv'"
	driveFetchSize = 100
)

func main() {
	flag.Parse()
	driveFolderId := flag.Arg(0)

	ctx := context.Background()

	var err error
	driveService, err = drive.NewService(ctx)
	if err != nil {
		panic(err)
	}

	query := fmt.Sprintf(driveQuery, driveFolderId)
	fl, err := driveService.Files.List().Q(query).PageSize(driveFetchSize).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok {
			if gerr.Code == http.StatusNotFound {
				fmt.Println("not found")
				os.Exit(1)
			}
		}
		panic(err)
	}

	for _, f := range fl.Files {
		fmt.Println(f.Name)
		_ = download(ctx, f)
	}
}

func download(ctx context.Context, file *drive.File) error {
	res, err := driveService.Files.Get(file.Id).Context(ctx).Download()
	if err != nil {
		return errors.Wrap(err, "failed to call download api")
	}
	defer res.Body.Close()

	filename := file.Name
	outFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy boby")
	}

	return nil
}
