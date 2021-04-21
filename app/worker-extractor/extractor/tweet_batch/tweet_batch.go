package tweet_batch

type TweetBatch struct {
	Id          int64   `avro:"id" json:"id"`
	ExtractedAt string  `avro:"extracted_at" json:"extracted_at"`
	Size        int     `avro:"size" json:"size"`
	Items       []Tweet `avro:"items" json:"items"`
}

type Tweet struct {
	Id            int64  `avro:"id" json:"id"`
	CreatedAt     string `avro:"created_at" json:"created_at"`
	Text          string `avro:"text" json:"text"`
	UserId        int64  `avro:"user_id" json:"user_id"`
	Retweeted     bool   `avro:"retweeted" json:"retweeted"`
	ReplyCount    int    `avro:"reply_count" json:"reply_count"`
	RetweetCount  int    `avro:"retweet_count" json:"retweet_count"`
	FavoriteCount int    `avro:"favorite_count" json:"favorite_count"`
	Lang          string `avro:"lang" json:"lang"`
	User          User   `avro:"user" json:"user"`
}

type User struct {
	Id             int64  `avro:"id" json:"id"`
	Location       string `avro:"location" json:"location"`
	FollowersCount int    `avro:"follower_count" json:"follower_count"`
	Verified       bool   `avro:"verified" json:"verified"`
}
