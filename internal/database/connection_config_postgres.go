package database

type ConnectionConfig struct {
	HostTag         HostTag `json:"HostTag"`
	ContainsSecrets bool    `json:"ContainsSecrets"`
	Host            string  `json:"Host"`
	User            string  `json:"User"`
	Password        string  `json:"Password"`
	Port            string  `json:"Port"`
	Database        string  `json:"Database"`
	SSL             bool    `json:"SSL"`
}
