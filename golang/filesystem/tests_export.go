package filesystem

func SaveCompositeDetailsForTest(c *Folder) {
	saveCompositeDetails(c)
}

func PopulateKeywordsFromStoredJsonFileForTest(c *Folder) {
	populateKeywordsFromStoredJsonFile(c)
}

func DeleteCompositeDetailsFileForTest(name string) error {
	return deleteCompositeDetailsFile(name)
}
