package main

import (
	"time"

	"github.com/Dextication/snowflake"
)

var (
	BASE_ID_TIME  = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	BASE_ID_STAMP = BASE_ID_TIME.UnixMilli()
	SnowNode      *snowflake.Node
)

func InitSnowflake() {
	node, err := snowflake.NewNode(0, BASE_ID_TIME, 41, 11, 11)
	PanicIfErr(err)
	SnowNode = node
}

// func SnowToTime(id string) time.Time {
// 	i, _ := strconv.ParseInt(id, 10, 64)

// 	timestamp := (i >> 22) + BASE_ID_STAMP

// 	return time.UnixMilli(timestamp)
// }

// func TimeToSnow(time time.Time) string {
// 	stamp := time.UnixMilli()
// 	stamp -= BASE_ID_STAMP

// 	return fmt.Sprint(stamp << 22)
// }
