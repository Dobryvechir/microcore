/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

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
