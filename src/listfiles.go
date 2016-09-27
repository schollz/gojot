package gitsdees

func ListFiles(gitfolder string) []string {
	branchNames, _ := ListBranches(gitfolder)
	infos, _ := GetInfo(gitfolder, branchNames)
	foundDocuments := make(map[string]bool)
	documents := []string{}
	for _, info := range infos {
		if _, ok := foundDocuments[info.Document]; !ok {
			foundDocuments[info.Document] = true
			documents = append(documents, info.Document)
		}
	}
	return documents
}
