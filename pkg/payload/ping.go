package payload

type Ping struct {
	APIBase
	Data *struct {
		APIDataBase
		UID  string `json:"uid,omitempty"`
		Page any    `json:"page,omitempty"`
	} `json:"data,omitempty"`
}
