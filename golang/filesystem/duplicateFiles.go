package filesystem

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
)

func FindDuplicateFiles(root *Folder) []*File {
	// First pass
	sizeBuckets := make(map[int64][]string)
	collectBySize(root, sizeBuckets)

	duplicates := []*File{}

	// Second pass
	for _, paths := range sizeBuckets {
		if len(paths) < 2 {
			continue
		}

		// 2a) sample-hash
		sampleBuckets := make(map[string][]string)
		for _, p := range paths {
			sig, err := sampleHash(p)
			if err != nil {
				log.Printf("sampleHash error for %s: %v", p, err)
				continue
			}
			sampleBuckets[sig] = append(sampleBuckets[sig], p)
		}

		// full-hash remaining groups
		for _, group := range sampleBuckets {
			if len(group) < 2 {
				continue
			}
			fullMap := make(map[string]string)
			for _, p := range group {
				fh, err := fullHash(p)
				if err != nil {
					log.Printf("fullHash error for %s: %v", p, err)
					continue
				}
				if existing, found := fullMap[fh]; found {
					duplicates = append(duplicates, &File{
						Name: filepath.Base(p),
						Path: existing,
					})
				} else {
					fullMap[fh] = p
				}
			}
		}
	}

	return duplicates
}

// collectBySize recurses through Folder, grouping files by their file size.
func collectBySize(item *Folder, buckets map[int64][]string) {
	for _, f := range item.Files {
		if info, err := os.Stat(f.Path); err == nil && info.Mode().IsRegular() {
			buckets[info.Size()] = append(buckets[info.Size()], f.Path)
		}
	}
	for _, sf := range item.Subfolders {
		collectBySize(sf, buckets)
	}
}

// sampleHash reads up to the first 4KB of a file and returns an MD5 signature.
func sampleHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	h := md5.Sum(buf[:n])
	return fmt.Sprintf("%x", h[:]), nil
}

// fullHash reads the entire file and returns its MD5 hash.
func fullHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return hashFile(md5.New(), f)
}

// hashFile streams data from r into the provided hash.Hash and returns the hex string.
func hashFile(h hash.Hash, r io.Reader) (string, error) {
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
