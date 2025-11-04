package s3

type DeleteResponse struct {
	DeletedKeys []string
	Errors      []ErrorResponse
}
