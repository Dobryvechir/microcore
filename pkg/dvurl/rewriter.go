package dvurl

type RewriteInfo struct {
	url string
	condition string
	options string
}

func createRewriteInfo(url string, condition string, options string, ids []string) *RewriteInfo {
	return &RewriteInfo{
		url: url,
		condition: condition,
		options: options,
	}
}