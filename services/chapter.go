package services

type Chapter struct {
	ID          int           `json:"id"`
	IDStr       string        `json:"id_str"`
	ClassID     int           `json:"class_id"`
	ClassIDStr  string        `json:"class_id_str"`
	OrderName   string        `json:"order_name"`
	Name        string        `json:"name"`
	Intro       string        `json:"intro"`
	PhaseNum    int           `json:"phase_num"`
	Status      int           `json:"status"`
	OrderNum    int           `json:"order_num"`
	IsFinished  int           `json:"is_finished"`
	UpdateTime  int           `json:"update_time"`
	LogID       string        `json:"log_id"`
	LogType     string        `json:"log_type"`
	ArticleList []ArticleInfo `json:"article_list"`
}
