package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
	driveQuery = "'%s' in parents and mimeType = 'text/csv' and not appProperties has { key='development' and value = 'done' }"
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
	fl, err := driveService.Files.List().Q(query).Fields("files(id,name,mimeType,size,createdTime,modifiedTime,appProperties)").OrderBy("modifiedTime").PageSize(driveFetchSize).Do()
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
		fmt.Println(f.Properties)
		fmt.Println(f.AppProperties)
		fmt.Println(f.ModifiedTime)
		//_ = download(ctx, f)
		fmt.Println(time.Now())
		ff := &drive.File{}
		ff.AppProperties = make(map[string]string)
		ff.AppProperties["development"] = "done"
		_, err := driveService.Files.Update(f.Id, ff).Context(ctx).Do()
		if err != nil {
			panic(err)
		}
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
