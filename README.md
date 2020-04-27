Golang Web Crawler

## Compile
`go build -o crawly crawly.go`

## Usage
Must complete, at least, <URL> and <Log File> input params
```
./crawly -u <URL> -l <Log File> -h <Auth Header> -t <Number of threads> <Cookie1=Value1> <Cookie2=Value2> ...
Info: Cookie must be set in 'Name=Value' format
```
