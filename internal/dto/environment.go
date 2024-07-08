package dto

// Env represents the environment configuration.
type Env struct {
	ServerPort                    string      // Porta do Server
	RedisAddr                     string      // Host do Redis
	RedisPort                     string      // Port do Redis
	MaxRequestsAllowedByIP        int         // Quantidade de requests permitida por IP
	MaxRequestsAllowedByToken     int         // Quantidade de requests permitida por Token
	OnRequestsExceededBlockBy     string      // O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
	TimeToReleaseRequestsIP       int         // As próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Valor em segundos
	TimeToReleaseRequestsToken    int         // As próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Valor em segundos
}
