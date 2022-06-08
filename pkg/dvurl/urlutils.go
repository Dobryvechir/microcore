/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvurl

func PrependUrl(url string, item string) string {
	if len(item) == 0 {
		return url
	}
	if item[0] != '/' {
		item = "/" + item
	}
	if len(url) == 0 {
		return item
	}
	if url[0] == '/' {
		return item + url
	}
	return item + "/" + url
}
