package revision

var (
	gitRevision string
)

func Revision() string {
	if len(gitRevision) > 6 {
		return gitRevision[0:6]
	}

	return gitRevision
}
