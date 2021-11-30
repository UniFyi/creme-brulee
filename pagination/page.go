package pagination

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	PageTimeFormat = "2006-01-02 15:04:05.000000"
	defaultPageSize = 3
)

type PageCursor struct {
	Num  uuid.UUID
	Time time.Time
}

type Summary struct {
	Current *PageCursor
	Next    *PageCursor
	Size       int
	NumResults int
}

type Pagination struct {
	Current    *string `json:"current"`
	Next       *string `json:"next"`
	Size       int     `json:"size"`
	NumResults int     `json:"numResults"`
}

type QP struct {
	PageNum  *string `form:"pageNum"`
	PageSize *int    `form:"pageSize"`
}

func ResolvePageSize(size *int) int {
	if size == nil {
		return defaultPageSize
	}
	if *size < 1 {
		// Negative page size or zero will result into usage of default page size
		return defaultPageSize
	}
	return *size
}

func FormatPageCursor(pageCursor *PageCursor) *string {
	if pageCursor != nil {
		formattedTime := pageCursor.Time.Format(PageTimeFormat)
		stringCursor := []byte(
			fmt.Sprintf("%v|%v", formattedTime, pageCursor.Num.String()),
		)
		result := base64.StdEncoding.EncodeToString(stringCursor)
		return &result
	}
	return nil
}

func FromPageSummary(pageSummary *Summary) *Pagination {
	return &Pagination{
		Current:    FormatPageCursor(pageSummary.Current),
		Next:       FormatPageCursor(pageSummary.Next),
		Size:       pageSummary.Size,
		NumResults: pageSummary.NumResults,
	}
}
