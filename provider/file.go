package provider

import (
	"github.com/EbonJaeger/mcsmanager/util"
)

// Update downloads a file from a given URL.
func (f File) Update(path string) error {
	return util.DownloadFile(f.URL, path)
}
