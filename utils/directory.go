package utils

import (
	"fmt"
	"os"
	"path"

	"google.golang.org/api/drive/v3"
)

func ChekcIfDirectoryExists(dir string) {
	os.Mkdir(dir, os.ModePerm)
}

func CreateDirectory(service *drive.Service) (*drive.File, error) {
	return service.Files.Create(&drive.File{
		Name:        driveDirectoryName,
		Kind:        "drive#file",
		MimeType:    "application/vnd.google-apps.folder",
		Description: fmt.Sprintf("%v Container Folder", driveDirectoryName),
	}).Do()
}

func GetDirectoryId(service *drive.Service) string {
	res, err := service.Files.List().PageSize(1000).Do()
	HandleError("Unable to retrieve Drive client", err)

	for _, file := range res.Files {
		if file.Name == driveDirectoryName {
			return file.Id
		}
	}

	dir, err := CreateDirectory(service)
	HandleError("Canont create directory on Google Drive", err)

	return dir.Id
}

func UploadAllFiles(service *drive.Service, id string) {

	ChekcIfDirectoryExists(inputDirectory)

	files, err := os.ReadDir(inputDirectory)
	HandleError(fmt.Sprintf("Cannot open directory %v", inputDirectory), err)

	for _, file := range files {
		var fullPath string = path.Join(inputDirectory, file.Name())
		_, err := UploadFile(service, id, fullPath)
		HandleError(fmt.Sprintf("Cannot upload file %v", fullPath), err)
	}
}

func UploadFile(service *drive.Service, id string, filepath string) (*drive.File, error) {
	file, err := os.Open(filepath)
	HandleError(fmt.Sprintf("Cannot open file %v", filepath), err)

	info, err := file.Stat()
	HandleError(fmt.Sprintf("Cannot retrieve file stats %v", filepath), err)
	defer file.Close()

	return service.Files.Create(&drive.File{
		Name:    info.Name(),
		Parents: []string{id},
	}).
		Media(file).
		ProgressUpdater(func(now, size int64) {
			fmt.Printf("%d, %d\r", now, size)
		}).
		Do()
}
