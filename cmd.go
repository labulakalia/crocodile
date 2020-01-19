package main

import (
	"bytes"

	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	// "time"
)

func main() {
	a := []string{"id1","id2","id3"}
	WHERE id=
	fmt.Println(strings.Join(a, " OR "))
	// var a map[string]string


	// if a == nil {
	// 	goto Check
	// }
	
	// return 
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	// // reader := runcmd(ctx, "ping www.baidu.com -n 5")
	// reader := runapi(ctx, "http://www.google.com")
	// defer reader.Close()
	// var (
	// 	lastrecv []byte
	// 	out      = make([]byte, 1024)
	// 	buf      bytes.Buffer
	// )
	// for {
	// 	n, err := reader.Read(out)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		log.Println("read failed", err)
	// 		break
	// 	}
	// 	if n > 0 {
	// 		//fmt.Printf("%s", out[:n])
	// 		lastrecv = out[:n]
	// 		buf.Write(out[:n])
	// 	}
	// }
	// fmt.Println(buf.String())
	// code := strings.TrimRight(string(lastrecv), " ")
	// buf.Truncate(len(buf.Bytes()) - 4) // remove return code in run write
	// fmt.Println(buf.String())
	// fmt.Printf("exitcode:%s", code)
// Check:
// 		fmt.Println("goto")
}

func runcmd(ctx context.Context, cmd string) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		var exitCode = -1
		defer pw.Close()
		defer func() {
			pw.Write([]byte(fmt.Sprintf("%3d", exitCode))) // write exitCode,total 3 byte
		}()
		// tell the command to write to our pipe
		scmd := strings.Split(cmd, " ")
		cmd := exec.CommandContext(ctx, scmd[0], scmd[1:]...)
		cmd.Stdout = pw
		cmd.Stderr = pw
		err := cmd.Start()
		if err != nil {
			pw.Write([]byte(err.Error()))
			return
		}

		err = cmd.Wait()
		if err != nil {
			var customerr bytes.Buffer
			switch ctx.Err() {
			case context.DeadlineExceeded:
				customerr.WriteString("执行超时\n")
			case context.Canceled:
				customerr.WriteString("执行被终止\n")
			default:
				customerr.WriteString(err.Error())
			}
			pw.Write(customerr.Bytes())
			// try to get the exit code
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			}

		} else {
			exitCode = 0
		}

	}()
	return pr
}

// run http get url
func runapi(ctx context.Context, url string) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		var exitCode = -1
		defer pw.Close()
		defer func() {
			pw.Write([]byte(fmt.Sprintf("%3d", exitCode))) // write exitCode,total 3 byte
		}()
		// go1.13 use NewRequestWithContext
		// req, err := http.NewRequest(http.MethodGet, url, nil)
		// if err != nil {
		// 	pw.Write([]byte(err.Error()))
		// 	log.Println("NewRequest failed", err)
		// 	return
		// }
		// req = req.WithContext(ctx)

		// client := http.DefaultClient
		// resp, err := client.Do(req)
		// if err != nil {
		// 	var customerr bytes.Buffer
		// 	switch ctx.Err() {
		// 	case context.DeadlineExceeded:
		// 		customerr.WriteString("执行超时\n")
		// 	case context.Canceled:
		// 		customerr.WriteString("执行被终止\n")
		// 	default:
		// 		customerr.WriteString(err.Error())
		// 	}
		// 	pw.Write(customerr.Bytes())
		// 	return
		// }
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
		defer ts.Close()
		req,_ := http.NewRequest("GET", ts.URL, nil)
		resp,err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		// res,_ := ioutil.ReadAll(resp.Body)
		// fmt.Printf("%s\n",res)
		// return
		var out = make([]byte, 1)
		for {
			n, err := resp.Body.Read(out)
			if err != nil {
				if err == io.EOF {
					log.Println("read done")
					break
				}
				log.Print("read failed", err)
				return
			}
			if n > 0 {
				pw.Write(out[:n])
			}
		}
		if resp.StatusCode > 0 {
			exitCode = resp.StatusCode
		}
	}()
	return pr
}
