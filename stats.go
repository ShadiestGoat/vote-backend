package main

type Photo struct {
	For     int
	Against int
	ID      string
	Picture []byte
}

type Stat struct {
	PictureID string `json:"pictureID"`
	Likes     int    `json:"likes"`
}

type StatsT struct {
	Top10 []*Stat `json:"top10"`
	Total int     `json:"total"`
}

func GetStats() *StatsT {
	s := &StatsT{
		Top10: []*Stat{},
	}
	DBQueryRow(`SELECT COUNT(*) FROM votes`).Scan(&s.Total)

	resp, err := DBQuery(`SELECT photo, SUM(vote) AS total_votes
FROM votes
GROUP BY photo
ORDER BY total_votes DESC, photo DESC
FETCH FIRST 10 ROWS WITH TIES`)

	PanicIfErr(err)
	
	for resp.Next() {
		tmp := &Stat{
			
		}
		resp.Scan(&tmp.PictureID, &tmp.Likes)
		s.Top10 = append(s.Top10, tmp)
	}
	
	return s
}