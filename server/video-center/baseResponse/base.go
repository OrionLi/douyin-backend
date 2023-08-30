package baseResponse

type Video struct {
	Id            int64  `json:"id"`
	User          User   `json:"user"`
	PlayUrl       string `json:"playUrl"`
	CoverUrl      string `json:"coverUrl"`
	FavoriteCount int64  `json:"favoriteCount"`
	CommentCount  int64  `json:"commentCount"`
	IsFavorite    bool   `json:"isFavorite"`
	Title         string `json:"title"`
}

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"followCount"`
	FollowerCount int64  `json:"followerCount"`
	IsFollow      bool   `json:"isFollow"`
}

type VideoArray []Video
