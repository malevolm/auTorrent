auTorrent 0.1.0
=========
a minimal torrent autodownloader for tv shows. meant to be used with a torrent client which watches a directory for new torrents, such as uTorrent

auTorrent is written to be platform independent, however you must make sure your shows.txt and config.txt contain the right type of newline for your system. on windows you must use CR+LF, on everything else its just LF. the files in this repo are CR+LF by default

auTorrent-tray
==============
* no console window, docks in system tray
* no interaction, click the tray icon once to close
* if building auTorrent-tray, use the following command
 
> go build -ldflags -H=windowsgui auTorrent-tray.go

about
==========
* desired tv shows are specified via shows.txt
* sleep settings & save path for torrents in config.txt
* download.log tracks your progression through tv series', deleting it will lead to duplicates unless you also update your shows.txt
* only supported torrent search is thepiratebay.se, should suffice for 99% of tv shows
* changes in shows.txt and config.txt are reflected upon next iteration (no need to restart app)

shows.txt
=========
the first line is the show's name
> My First Show

the second is a regexp for the show torrent name, must match your desired release, and capture the season and episode as submatches. you can prepend your expression with (?i) to toggle the case-insensitivity flag (recommended)
> (?i)my.first.show.s(\d+)e(\d+).+?720p

the third line specifies the last episode you've seen, the program will download the next episode after the one you specify here. in this example it will start with season 4 episode 1
> 4 0

the fourth line is a whitelist for reputable torrent uploaders, separated by spaces (eg. "eztv ettv Drarbg"). if you dont want a whitelist, just put a wildcard

> *

and that's it. each block should be separated by a newline (currently expects \r\n)
