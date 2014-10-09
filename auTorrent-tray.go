package main

import (
    "os"
	"fmt"
	"time"
	"regexp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"github.com/xilp/systray"
)

var shows [][]string
var config = make(map[string]string)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	tray := systray.New("\\", "tray.ico", 6333)
	tray.OnClick(func() {
		os.Exit(0)
	})
	
	err := tray.Show("tray.ico", "auTorrent")
	if err != nil {
		println(err.Error())
	}

	go func() {
		fmt.Println("==auTorrent Initialised==")
		
		for i := 0; i > -1; i++ {
			loadConfig()
			buildShowDB()
			sleep, _ := strconv.Atoi(config["SLEEP_SEC"])
			
			checkForNewTorrents()
			fmt.Printf("(sleeping for %v sec)\n", sleep)
			time.Sleep(time.Duration(sleep)*1000000000)
		}
	
		err = tray.Stop()
		if err != nil {
			println(err.Error())
		}
		os.Exit(0)
	}()

	err = tray.Run()
	if err != nil {
		println(err.Error())
	}
}

func checkForNewTorrents() {
	for _, show := range shows {
		fmt.Printf("Searching for new episodes of \"%s\"\n", show[0])
		res := torrentSearch(show[0]);
		for i := 0; i < len(res); i++ {
			result := checkTorrentSuitability(res[i][0], res[i][1])
			if result == 1 {
				info := getSeasonInfo(res[i][0], show[1])
				if checkIfDownloaded(show[0], info[0], info[1]) == false {
					fmt.Printf("=> %s (by %s)\n", res[i][0], res[i][1])
					downloadTorrent(res[i])
					markAsDownloaded(show[0], info[0], info[1])
				}
			}
		}
	}
}

func torrentSearch(query string) [][]string {
	url := fmt.Sprintf("http://thepiratebay.se/search/%s/0/3/0", url.QueryEscape(query))
	src := httpGet(url)
	
	re := regexp.MustCompile("(?s)<a href=\"/torrent/(.+?)/(.+?)\".+?Details for (.+?)\".+?Browse (.+?)\"")
	rr := re.FindAllStringSubmatch(src, -1)
	
	out := make([][]string, len(rr))
	
	for i := 0; i < len(rr); i++ {
		out[i] = []string{rr[i][3], rr[i][4], rr[i][1], rr[i][2]}
	}
	
	return out
}

func checkTorrentSuitability(torrent string, author string) int {
	for _, parts := range shows {
		re := regexp.MustCompile(parts[1])
		rr := re.FindAllStringSubmatch(torrent, -1)
		start := strings.Split(parts[2], " ")
		authors := strings.Split(parts[3], " ")
		
		if len(rr) > 0 {
			s1, _ := strconv.Atoi(rr[0][1])
			e1, _ := strconv.Atoi(rr[0][2])
			s2, _ := strconv.Atoi(start[0])
			e2, _ := strconv.Atoi(start[1])
			
			if s1 > s2 || ( s1 == s2 && e1 > e2 ) {
				if stringInSlice(author, authors) || authors[0] == "*" {
					return 1
				}
			}
		}
	}
	
	return 0
}

func getSeasonInfo(torrent string, expr string) []string {
	re := regexp.MustCompile(expr)
	rs := re.FindAllStringSubmatch(torrent, -1)	
	return []string{rs[0][1], rs[0][2]}
}

func downloadTorrent(torrent_info []string) int {
	file := fmt.Sprintf("%s%s.torrent", config["SAVE_PATH"], torrent_info[3])
	link := fmt.Sprintf("http://torrents.thepiratebay.se/%s/%s.torrent", torrent_info[2], torrent_info[3])
	data := httpGet(link)
	ioutil.WriteFile(file, []byte(data), 0600)
	return 1
}

func checkIfDownloaded(show string, season string, episode string) bool {
	res, err := ioutil.ReadFile("download.log");
	check(err)
	lines := strings.Split(string(res[:]), "\r\n")
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if parts[0] == show && parts[1] == season && parts[2] == episode {
			return true
		}
	}
	return false
}

func markAsDownloaded(show string, season string, episode string) bool {
	if checkIfDownloaded(show, season, episode) == false {
		entry := fmt.Sprintf("%s|%s|%s\r\n", show, season, episode)
		f, err := os.OpenFile("download.log", os.O_APPEND|os.O_WRONLY, 0600)
		check(err)

		defer f.Close()

		if _, err = f.WriteString(entry); err != nil {
			panic(err)
		}
		
		return true
	} else {
		return false
	}
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func httpGet(url string) string {
	res, err := http.Get(url)
	check(err)
	
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	check(err)
	
	out := string(result[:])
	
	return out
}

func loadConfig() {
	res, err := ioutil.ReadFile("config.txt");
	check(err)
	
	lines := strings.Split(string(res[:]), "\r\n")
	for _, line := range lines {
		if line != "" && string(line[0]) != "#" {
			info := strings.Split(line, "=")
			config[info[0]] = info[1]
		}
	}
}

func buildShowDB() {
	res, err := ioutil.ReadFile("shows.txt");
	check(err)
	
	db := strings.Split(string(res[:]), "\r\n\r\n")
	shows = make([][]string, len(db))
	
	for i, chunk := range db {
		parts := strings.Split(chunk, "\r\n")
		shows[i] = []string{parts[0], parts[1], parts[2], parts[3]}
	}
}