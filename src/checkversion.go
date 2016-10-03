package sdees

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func CheckNewVersion(dir string, version string, build string, osType string) {
	logger.Debug("Current executable path: %s", dir)
	if version == "dev" {
		updateDevVersion(dir, version, build, osType)
	} else {
		updateDownloadVersion(dir, version, build, osType)
	}
}

func updateDevVersion(dir string, version string, build string, osType string) {
	logger.Debug("Updating dev version of sdees")
	url := "https://api.github.com/repos/schollz/sdees/commits"
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	var j GithubCommitsJSON
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		log.Fatal(err)
	}
	logger.Debug(j[0].Sha)

	os.Exit(0)
}

func updateDownloadVersion(dir string, version string, build string, osType string) {
	newVersion, versionName := checkGithub(version)
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
	downloadVersion := versionName + "/sdees_" + osType + ".zip"
	os.Remove("sdees_" + osType + ".zip")
	fmt.Printf("\nDownloading %s...", downloadVersion)
	cmd := exec.Command("wget", "https://github.com/schollz/sdees/releases/download/"+downloadVersion)
	_, err := cmd.Output()
	if err != nil {
		logger.Error("Problem downloading, do you have internet?")
		log.Fatal(err)
	}

	logger.Debug("Removing old version: %s", dir)
	err = os.Remove(dir)
	if err != nil {
		logger.Error("Problem removing file, do you need sudo?")
		log.Fatal(err)
	}

	logger.Debug("Unzipping new version")
	cmd = exec.Command("unzip", "sdees_"+osType+".zip")
	_, err = cmd.Output()
	if err != nil {
		logger.Error("Problem unzipping, do you have zip?")
		log.Fatal(err)
	}

	logger.Debug("Cleaning...")
	os.Remove("sdees_" + osType + ".zip")
	fmt.Printf("\n\nsdees Version %s installed!\n", versionName)
	os.Exit(0)
}

func checkGithub(version string) (bool, string) {
	url := "https://api.github.com/repos/schollz/sdees/releases/latest"
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
		return false, ""
	}
	majorMinorWeb := []int{}
	for i := 0; i < 3; i++ {
		i, _ := strconv.Atoi(versions[i])
		majorMinorWeb = append(majorMinorWeb, i)
	}

	versions = strings.Split(version, ".")
	if len(versions) != 3 {
		return false, ""
	}
	majorMinor := []int{}
	for i := 0; i < 3; i++ {
		i, _ := strconv.Atoi(versions[i])
		majorMinor = append(majorMinor, i)
	}

	newVersionAvailable := false
	for i := range majorMinor {
		if majorMinorWeb[i] > majorMinor[i] {
			newVersionAvailable = true
			break
		}
	}

	return newVersionAvailable, newVersion
}
