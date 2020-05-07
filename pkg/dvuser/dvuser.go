package dvuser

type UserInfo struct {
	Registered bool
	Name       string
	Email      string
	Roles      []string
}

func GetCurrentUserInfo() UserInfo {
	return UserInfo{}
}

func GetSourceInfo(sourceRefs []string, baseInfo map[string]string) map[string]string {
	return baseInfo
}
