package dto

// Env represents the environment configuration.
type Env struct {
	ServerPort                    string      // Porta do Server
	RedisAddr                     string      // Host e Port do Redis
	MaxRequestsAllowedByIP        int         // Quantidade de requests permitida por IP
	MaxRequestsAllowedByToken     int         // Quantidade de requests permitida por Token
	LimitRequestsByIP             bool        // Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso. Se true é por IP senão é por Token
	OnRequestsExceededBlockBy     string      // O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
	TimeToReleaseRequests         int         // As próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Valor em segundos
}
