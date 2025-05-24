package config

type Config struct {
	Payos     *Payos    `json:"payos"`
	Server    *Server   `json:"server"`
	DataBase  *DataBase `json:"database"`
	Bulbasaur *Server   `json:"bulbasaur"`
}

type Server struct {
	Host string `json:"hosts"`
	Port int    `json:"port"`
}

type Payos struct {
	ClientId   string `json:"client_id"`
	ApiKey     string `json:"api_key"`
	Checksum   string `json:"checksum_key"`
	CancelUrl  string `json:"cancel_url"`
	SuccessUrl string `json:"success_url"`
	TimeOut    int    `json:"timeout"`
}

type DataBase struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

// payos:
//   api_key: ${PAYOS_API_KEY}
//   api_secret: ${PAYOS_API_SECRET}
//   checksum_key: ${PAYOS_CHECKSUM_KEY}
//   cancel_url: https://api.payos.com
//   success_url: https://api.payos.com
