package gitdomain

type Querier interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
	QueryWithEnv(env []string, executable string, args ...string) (string, error)
}
