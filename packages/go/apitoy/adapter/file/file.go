package file

type Adapter struct {
	tempDir          string
	ingestFilePrefix string
}

func NewAdapter(tempDir string) Adapter {
	return Adapter{
		tempDir:          tempDir,
		ingestFilePrefix: "bh",
	}
}
