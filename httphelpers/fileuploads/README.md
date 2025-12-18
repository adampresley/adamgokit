# File Uploads

This package provides helper methods to simplify handling file uploads from `multipart/form-data` requests.

## UploadFileToDir

**UploadFileToDir** is a convenience function that reads a file from a request and saves it to a specified directory. It uses a template to generate the destination filename.

```go
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	destDir := "/path/to/uploads"
	fileNameTemplate := "{{.randomString}}-{{.fileName}}"

	fileUpload, err := fileuploads.UploadFileToDir(
		"my-file",
		r,
		destDir,
		fileNameTemplate,
		fileuploads.WithMaxFileSize(20<<20), // 20MB
	)

	if err != nil {
		// Handle error
		return
	}

	fmt.Printf("File saved to %s\n", fileUpload.SavedFilePath)
}
```

### Filename Templating

The `destFileNameTemplate` parameter uses Go's `text/template` engine. The following variables are available to use in your template:

-   `baseFileName`: The base name of the file, which is the file name without any directory path (e.g., `document.pdf`).
-   `ext`: The extension of the file (e.g., `.pdf`).
-   `fileName`: The original name of the uploaded file as provided by the client (e.g., `my-folder/document.pdf`).
-   `size`: The size of the file in bytes.
-   `randomString`: A random string of alphanumeric digits, useful for preventing filename collisions. The default size is 10 characters.

## ReadUploadedFile

For more control over the uploaded file, you can use **ReadUploadedFile**. This function parses the form and provides the file data to a callback function for custom processing, without automatically saving it to disk.

```go
func ProcessUploadHandler(w http.ResponseWriter, r *http.Request) {
	err := fileuploads.ReadUploadedFile("my-file", r, func(file fileuploads.FileUpload, options *fileuploads.UploadOptions) error {
		// You have access to the file via file.File (multipart.File)
		// Process the file here, e.g., upload to S3, scan for viruses, etc.

		fmt.Printf("Processing file %s\n", file.FileName)
		return nil
	})

	if err != nil {
		// Handle error
	}
}
```

## Options

You can customize the upload behavior by passing in one or more option functions.

### WithMaxFileSize

**WithMaxFileSize** overrides the maximum allowed file size. The default is 10MB (`10 << 20`).

```go
// 50 MB limit
fileuploads.WithMaxFileSize(50 << 20)
```

### WithRandomStringSize

**WithRandomStringSize** sets the length of the random string available in the filename template. The default is 10.

```go
fileuploads.WithRandomStringSize(15)
```

