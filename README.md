# blobby

blobby is a JSON blob analyzer. It breaks down JSON by field so that you can visualize the makeup of it.

Inspiration for this tool came from a ticket I worked on recently where we wanted to improve our FCP (first content paint) by reducing the size of the data we were fetching in NextJS' `GetServerSideProps`. I wanted to see which fields were the biggest contributors to a fairly chunky JSON blob and then (assuming it wasn't first paint critical for SEO reasons for example) would be left to fetch on the client later on.

## Development

1. Install [Golang](https://go.dev/doc/install)
2. Install dependencies with `go get .` (at time of writing there are no dependencies outside of the Go standard library so technically not necessary)
3. Build with `go build`
4. Run with either `./blobby` (the generated executable) or `go run main.go`
