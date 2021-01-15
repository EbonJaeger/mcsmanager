package provider

// Download downloads a file from a given URL.
func (f File) Download(path string) error {
	return DownloadFile(f.URL, path)
}
