Golang Web Crawler

## Compile
`go build -o crawly crawly.go`

## Usage
```
./crawly -u <URL> -l <Log File> -a <Auth Header> -d <Depth> -t <Number of threads> <Cookie1=Value1> <Cookie2=Value2> ...

-a string
	Authorization Basic Header (optional)
-d int
	Crawling depth (0 -> only input url, 1 -> infinite) (default 1)
-l string
	Log file name (optional)
-t int
	Number of threads (optional) (default 10)
-u string
	URL to crawl (mandatory)

You can set cookies as the last argument like Cookie1=Value1 Cookie2=Value2
```
