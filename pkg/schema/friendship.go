package schema

type Friendship struct {
	User   int64 `json:"id"`
	Friend int64 `json:"friend"`
}
