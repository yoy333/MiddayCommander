package actions

// OpType identifies a file operation.
type OpType int

const (
	OpCopy OpType = iota
	OpMove
	OpDelete
	OpMkdir
	OpRename
)

// Progress reports the state of an ongoing file operation.
type Progress struct {
	Op         OpType
	TotalFiles int
	DoneFiles  int
	TotalBytes int64
	DoneBytes  int64
	Current    string // current file being processed
	Err        error
}
