package err

type InvalidTopicData struct {
	*CustomError
}

func NewInvalidTopicData(msg string) *InvalidTopicData {
	return &InvalidTopicData{
		CustomError: &CustomError{
			Msg: "Invalid TopicData: " + msg,
		},
	}
}
