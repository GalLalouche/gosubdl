This is a dead-simple and stupid Golang utility for downloading subtitles using the Subdl API. Mostly done as an experiment to see if Go is as bad as I think it is.
Conclusion: it's worse.

# Building
Since this relies on an API key, there are two ways of using it:
* Embed one during build:
```sh
go build -ldflags="-X 'gosubdl/requests.apiKeyDuringBuild=$SUBDL_API_KEY'"
```
* Define an environment variable during run (this will override the above)
