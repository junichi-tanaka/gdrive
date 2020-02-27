package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

func main() {
	flag.Parse()
	ctx := context.Background()

	driveService, err := drive.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create drive.Service %v", err)
	}

	// id of source file
	id := flag.Arg(0)
	src, err := driveService.Files.Get(id).Context(ctx).Do()
	if err != nil {
		log.Fatalf("failed to get %v", err)
	}
	fmt.Println(src)

	dstFolder := flag.Arg(1)
	dstName := flag.Arg(2)
	file := &drive.File{Name: dstName, Description: dstName, DriveId: src.DriveId, Parents: []string{dstFolder}}
	res, err := driveService.Files.Copy(id, file).Context(ctx).Do()
	if err != nil {
		log.Fatalf("failed to copy %v", err)
	}
	fmt.Println(res)
}
