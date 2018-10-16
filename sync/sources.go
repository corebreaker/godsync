package sync

func scanSourceDir(root, dir string, writer chan<- *tFileList) {

}

func ScanSources(root string, dirs []string) []*FileDesc {
	list := new(tFileList)
	wr := list.writer()
	defer close(wr)

	for _, dir := range dirs {
		scanSourceDir(root, dir, wr)
	}

	return list.list
}
