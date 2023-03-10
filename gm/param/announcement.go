package param

type AddAnnouncement struct {
	Type          int32  `json:"type"`
	Module        int32  `json:"module"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Tag           int32  `json:"tag"`
	Image         string `json:"image"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	ShowStartTime int64  `json:"show_start_time"`
	ShowEndTime   int64  `json:"show_end_time"`
	Priority      int32  `json:"priority"`
}

type UpdateAnnouncement struct {
	Id            int64  `json:"id"`
	Type          int32  `json:"type"`
	Module        int32  `json:"module"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Tag           int32  `json:"tag"`
	Image         string `json:"image"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	ShowStartTime int64  `json:"show_start_time"`
	ShowEndTime   int64  `json:"show_end_time"`
	Priority      int32  `json:"priority"`
}

type FetchAnnouncement struct {
	Module int32 `json:"module"`
}

type DeleteAnnouncement struct {
	Id int64 `json:"id"`
}

type AddBanner struct {
	Type          int32  `json:"type"`
	Module        int32  `json:"module"`
	Image         string `json:"image"`
	Jump          string `json:"jump"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	ShowStartTime int64  `json:"show_start_time"`
	ShowEndTime   int64  `json:"show_end_time"`
	Priority      int32  `json:"priority"`
}

type UpdateBanner struct {
	Id            int64  `json:"id"`
	Type          int32  `json:"type"`
	Module        int32  `json:"module"`
	Image         string `json:"image"`
	Jump          string `json:"jump"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	ShowStartTime int64  `json:"show_start_time"`
	ShowEndTime   int64  `json:"show_end_time"`
	Priority      int32  `json:"priority"`
}

type FetchBanner struct {
	Module int32 `json:"module"`
}

type DeleteBanner struct {
	Id int64 `json:"id"`
}

type AddCaution struct {
	Content   string `json:"content"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type UpdateCaution struct {
	Id        int64  `json:"id"`
	Content   string `json:"content"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type DeleteCaution struct {
	Id int64 `json:"id"`
}
