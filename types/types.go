package types

type Distance struct {
	Values float64 `json:"value"`
	OBUID  int     `json:"obuID"`
	Unix   int64   `json:"unix"`
}

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}
