package provider

// Update downloads a file from a given URL.
func (f File) Update(path string) error {
	return DownloadFile(f.URL, path)
}
