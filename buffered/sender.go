package buffered

type Sender interface {
	UsageSource
	ChangeSource
	Buffer
}
