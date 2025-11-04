package s3

type ListResponse struct {
	ContinuationToken string
	NumObjects        int
	Objects           []Object
}
