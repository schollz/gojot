package gojot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type GithubCommitsJSON []struct {
	Sha    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL          string `json:"url"`
		CommentCount int    `json:"comment_count"`
	} `json:"commit"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Author      struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Committer struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"committer"`
	Parents []struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"parents"`
}

type GithubJson struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

func CheckNewVersion(dir string, version string, lastcommit string, osType string) {
	logger.Debug("Current executable path: %s", dir)
	if version == "dev" {
		updateDevVersion(lastcommit)
	} else {
		updateDownloadVersion(dir, version, lastcommit, osType)
	}
}

func updateDevVersion(lastcommit string) error {
	logger.Debug("Updating dev version of gojot")
	url := "https://api.github.com/repos/schollz/gojot/commits"
	r, err := http.Get(url)
	if err != nil {
		return errors.New("Couldn't call Github API for getting new date")
	}
	defer r.Body.Close()
	var j GithubCommitsJSON
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		return errors.New("Couldn't decode Github API")
	}

	if len(j) == 0 {
		return errors.New("No data form Github!")
	}
	currentCommit, err := ParseDate(strings.Replace(lastcommit, "'", "", -1))
	if err != nil {
		return errors.New("Couldn't parse Github API Commit date")
	}
	logger.Debug("Github: %s, Current: %s", j[0].Commit.Author.Date.String(), currentCommit.String())
	if currentCommit.Sub(j[0].Commit.Author.Date).Hours() < 0 {
		fmt.Println("New version of gojot available! Run\n\n\tgo get -u github.com/schollz/gojot\n\not download.")
	}
	return nil
}

func updateDownloadVersion(dir string, version string, lastcommit string, osType string) {
	newVersion, versionName, _ := checkGithub(version)
	if !newVersion {
		logger.Debug("Current version is up to date: %s / %s", version, versionName)
		return
	}
	var yesnoall string
	fmt.Printf("\nVersion %s is available. Download? (y/n) ", versionName)
	fmt.Scanln(&yesnoall)
	if yesnoall == "n" {
		return
	}
	err := mkdir("tmp11")
	if err != nil {
		logger.Debug("Couldn't make tmp directory!")
		return
	}
	cwd, _ := os.Getwd()
	os.Chdir("tmp11")
	defer os.Chdir(cwd)

	// Download
	downloadVersion := versionName
	downloadName := "gojot_" + osType + ".zip"
	fmt.Printf("\nDownloading %s/%s...", downloadVersion, downloadName)
	err = DownloadFile(downloadName, "https://github.com/schollz/gojot/releases/download/"+downloadVersion+"/"+downloadName)
	if err != nil {
		logger.Debug("Problem downloading file: %s", "https://github.com/schollz/gojot/releases/download/"+downloadVersion+"/"+downloadName)
	}

	// Unzip
	logger.Debug("Unzipping new version")
	err = Unzip(downloadName, "./")
	if err != nil {
		logger.Debug("Problem unzipping file")
	}
	os.Remove(downloadName)

	// Move file
	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		if strings.Contains(f.Name(), "gojot") {
			logger.Debug("Moving %s to %s", f.Name(), dir)
			err = os.Remove(dir)
			if err != nil {
				logger.Debug("Error deleting file: %s", err.Error())
			}
			err = CopyFile(f.Name(), dir)
			if err != nil {
				logger.Debug("Error moving file: %s", err.Error())
			}
		}
	}

	// Clean
	logger.Debug("Cleaning...")
	os.Chdir(cwd)
	os.RemoveAll("tmp11")

	// Done!
	fmt.Printf("\n\ngojot Version %s installed!\n", versionName)
	os.Exit(0)
}

func checkGithub(version string) (bool, string, []int) {
	url := "https://api.github.com/repos/schollz/gojot/releases/latest"
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	var j GithubJson
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		log.Fatal(err)
	}
	newVersion := j.TagName
	versions := strings.Split(newVersion, ".")
	if len(versions) != 3 {
		return false, "", []int{}
	}
	majorMinorWeb := []int{}
	for i := 0; i < 3; i++ {
		num, _ := strconv.Atoi(versions[i])
		majorMinorWeb = append(majorMinorWeb, num)
	}

	versions = strings.Split(version, ".")
	if len(versions) != 3 {
		return false, "", []int{}
	}
	majorMinor := []int{}
	for i := 0; i < 3; i++ {
		num, _ := strconv.Atoi(versions[i])
		majorMinor = append(majorMinor, num)
	}

	newVersionAvailable := false
	for i := range majorMinor {
		if majorMinorWeb[i] > majorMinor[i] {
			newVersionAvailable = true
			break
		}
	}

	return newVersionAvailable, newVersion, majorMinorWeb
}
