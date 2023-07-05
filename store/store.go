package store

type Parameters struct {
	bigUrl string
	key    string
}

var storedURLs []Parameters

// storedURLS := []Parameters{}

func saveURL(bigUrl string, key string) bool {
	storedURLs = append(storedURLs, Parameters{
		bigUrl: bigUrl,
		key:    key,
	})
	return true
}

func retriveURL() string {
	return ""
}
