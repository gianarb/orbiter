package hook

import "github.com/Sirupsen/logrus"

type Channel struct {
	c chan *logrus.Entry
}

func NewChannelHook(cc chan *logrus.Entry) Channel {
	return Channel{
		c: cc,
	}
}

func (channel Channel) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (channel Channel) Fire(entry *logrus.Entry) error {
	select {
	case channel.c <- entry:
	default:
	}

	return nil
}
