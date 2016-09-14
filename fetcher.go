package main

import (
	"./requeststore"
	"fmt"
	"log"
	"net/http"
	"os"
	//	"bytes"
	//	"io"
	"bufio"
	"flag"
	"github.com/cheggaaa/pb"
	"strconv"
	"strings"
)

/*
func headersToWriter(h http.Header, w io.Writer) error {
        if err := h.Write(w); err != nil {
                return err
        }
        // ReadMIMEHeader expects a trailing newline
        _, err := w.Write([]byte("\r\n"))
        return err
}


func requestStr(h http.Request) string {
        hb := &bytes.Buffer{}
        hb.Write([]byte(fmt.Sprintf("%s %s %s\r\n", h.Method, h.URL.String(), h.Proto)))
        headersToWriter(h.Header, hb)

        return hb.String()
}
*/
const (
	defaultDir = "./storedata"
)

var (
	verbose bool
	//retries      int
	//hashSum      bool
	//useDisk      bool
	//private      bool
	dir string
	//dumpHttp bool
	noprivate            bool
	bypassprivateconfirm bool
	requestdata          bool
)

func init() {
	flag.StringVar(&dir, "dir", defaultDir, "the dir to store cache data in, implies -disk")
	flag.BoolVar(&noprivate, "no-private", true, "Do not store private responses")
	flag.BoolVar(&bypassprivateconfirm, "as-script", false, "Runs under script and ignores specific errors")
	flag.BoolVar(&requestdata, "log-request", false, "Print request details")

	flag.Parse()

}
func main() {
	//useDisk := true
	//dir := "./storedata"

	var store requeststore.Store

	log.Printf("storing cached resources in %s", dir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatal(err)
	}
	var err error
	store, err = requeststore.NewDiskStore(dir)
	if err != nil {
		log.Fatal(err)
	}

	//storeQueue := requeststore.NewMemoryStore()

	client := &http.Client{}
	/*
		req, _ = http.NewRequest("GET", reqUrl, nil)
		req.Header.Add("User-Agent", "Microsoft BITS/7.8")
		req.Header.Add("Cache-Control", "max-age=259200")

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp)
			fmt.Println(store.StoreResponse(requeststore.NewResponseFromHttp(resp, resp.Body), resp.Request.Method+":"+resp.Request.URL.String(), true))
		}
	*/

	/*
		stats, names := store.WalkRequests()
		fmt.Println(names)
		fmt.Println(stats)
	*/
	reader := bufio.NewReader(os.Stdin)

	allRequests := store.RetrieveAllRequests()
	for _, request := range allRequests {

		if request.Method == "GET" || request.Method == "HEAD" {
			if len(request.Header.Get("Range")) > 0 {
				tmpHead := request.Header
				tmpHead.Del("Range")
				request.Header = tmpHead
			}
			if requestdata {
				fmt.Println(request)
			}
			fmt.Println("ORIGINAL KEY =>", request.Method+":"+request.URL.String())
			fmt.Println("ORIGNAL KEY HASH =>", store.HashKey(request.Method+":"+request.URL.String()))
			fmt.Println("NEW KEY =>", request.Method+":"+request.URL.Scheme+"://"+"msupdates.ngtech.internal"+request.URL.Path)
			fmt.Println("NEW KEY HASH =>", store.HashKey(request.Method+":"+request.URL.Scheme+"://"+"msupdates.ngtech.internal"+request.URL.Path))

			switch {
			case request.Method == "GET":
				vfsResp, err := store.RetrieveResponse(request.Method + ":" + request.URL.Scheme + "://" + "msupdates.ngtech.internal" + request.URL.Path)
				if err == nil {
					fmt.Println("found on disk", vfsResp)
					vfsResp.Close()
					continue
				}
			case request.Method == "HEAD":
				vfsResp, err := store.RetrieveResponseHeader(request.Method + ":" + request.URL.Scheme + "://" + "msupdates.ngtech.internal" + request.URL.Path)
				if err == nil {
					fmt.Println("found on disk", vfsResp)
					//vfsResp.Close()
					continue
				}
			default:
				continue
			}
			resp, err := client.Do(&request)

			if err != nil {
				fmt.Println(err)
			} else {

				fmt.Println(resp)

				if noprivate {
					switch {
					case strings.Contains(strings.ToLower(resp.Header.Get("Cache-Control")), "private"):
						resp.Body.Close()
						continue
					case strings.Contains(strings.ToLower(resp.Header.Get("Cache-Control")), "no-store"):
						resp.Body.Close()
						continue
					case strings.Contains(strings.ToLower(resp.Header.Get("Cache-Control")), "must-revalidate"):
						fmt.Println("Starting to store response, has a must-revalidate")
						//resp.Body.Close()
						//continue
					case strings.Contains(strings.ToLower(resp.Header.Get("Cache-Control")), "public"):
						fmt.Println("Starting to store response, has a public Cache-Control")
						//Ignore
					default:
						fmt.Println("Starting to store response, no Cache-Control or non-private")
						//Ignore
					}
				}

				var bar *pb.ProgressBar
				var storeResp error
				switch {
				case request.Method == "GET":
					if len(resp.Header.Get("Content-Length")) > 0 {
						i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
						bar = pb.New(i).SetUnits(pb.U_BYTES)
						bar.Start()

						// https://github.com/fujiwara/shapeio
						// Add some rate limiting around the rest.Body
						// create proxy reader
						rd := bar.NewProxyReader(resp.Body)
						// and copy from reader
						//io.Copy(file, rd)
						storeResp = store.StoreResponse(requeststore.NewResponseFromHttp(resp, rd), resp.Request.Method+":"+resp.Request.URL.Scheme+"://"+"msupdates.ngtech.internal"+resp.Request.URL.Path, false)
					} else {
						storeResp = store.StoreResponse(requeststore.NewResponseFromHttp(resp, resp.Body), resp.Request.Method+":"+resp.Request.URL.Scheme+"://"+"msupdates.ngtech.internal"+resp.Request.URL.Path, false)
					}
				case request.Method == "HEAD":
					storeResp = store.StoreResponseHeader(requeststore.NewResponseFromHttp(resp, resp.Body), resp.Request.Method+":"+resp.Request.URL.Scheme+"://"+"msupdates.ngtech.internal"+resp.Request.URL.Path, false)
				default:
					continue
				}
				//	storeResp = store.StoreResponse(requeststore.NewResponseFromHttp(resp, resp.Body), resp.Request.Method+":"+resp.Request.URL.Scheme+"://"+"msupdates.ngtech.internal"+resp.Request.URL.Path, false)
				fmt.Println(storeResp)
				if storeResp != nil && storeResp == requeststore.ErrFoundInStorePrivate {
					fmt.Println("Press anything to continue:")
					if !bypassprivateconfirm {
						_, _ = reader.ReadString('\n')
					}
				}
			}
		} else {
			fmt.Println("Request was cancled due to the method", request)
		}

	}
}
