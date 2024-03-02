package buffered

type BufferedSender interface {
	UsageSource
	ChangeSource
	Buffer
}
