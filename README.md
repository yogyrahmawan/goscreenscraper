goscreenscraper
===============

Go Screen Scrapper is library for screen scraping. It's support get and post method. 

INSTALL
========

go get github.com/lesvolanolas/goscreenscraper

USAGE
========

```
import (
	"fmt"
	"goscreenscraper"
)

func main() {
  //example call without redirect, post method 
	m := make(map[string]string)
	m["log"] = "shyrin"
	m["pwd"] = "your_own_password"
	m["wp-submit"] = "Log In"
	m["redirect_to"] = "http://www.tourismjournal.net/wp-admin/"
	m["testcookie"] = "1"
	req2 := screenscraper.NewRequest("POST","http://www.tourismjournal.net","/wp-login.php","http://www.tourismjournal.net/wp-login.php?2yniv1rfyct10sdlt1d1r",m,nil)
	resp2,_ := req2.CallWithoutRedirect()

	fmt.Println(resp2.Cookies)

  //example call and redirect, get method
	req3 := screenscraper.NewRequest("GET","http://www.tourismjournal.net","/wp-admin/","",nil,resp2.Cookies)
	resp3,_ := req3.CallAndRedirect()

	fmt.Println(string(resp3.RawContent))
}
```

Release Notes 
===============

See Notes Here on release.txt

Author
=======

Yogy Rahmawan




