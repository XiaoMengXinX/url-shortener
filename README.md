# url-shortener

A serverless short url api on vercel powered by golang

### API Calling Reference

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	resp, _ := http.PostForm("https://neko.center/api/url",
		url.Values{
			"url": {"https://www.baidu.com/"}, // The url to be shortened
			"token": {"baidu"},                // Custom shorten url token (optional)
		})
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

Response:

```json
{
  "token": "baidu",
  "error": ""
}

```
