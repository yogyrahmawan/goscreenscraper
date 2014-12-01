package goscreenscraper

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func Test_CallAndRedirect(Test *testing.T) {
	req := NewRequest("GET","http://www.tourismjournal.net","/wp-admin/","",nil,nil)
	resp,err := req.CallAndRedirect()

	if err != nil {
		fmt.Println("error ", err)
	}else {
		assert := assert.New(Test)
		assert.NotNil(resp)
		assert.Equal("http://www.tourismjournal.net", resp.LastUrl , "Redirect to main page")
		assert.NotNil(resp.RawContent)
		assert.NotNil(resp.Cookies)
	}
}

func TestCallWithoutRedirect(Test *testing.T) {
	m := make(map[string]string)
	m["log"] = "shyrin"
	m["pwd"] = "*"
	m["wp-submit"] = "Log In"
	m["redirect_to"] = "http://www.tourismjournal.net/wp-admin/"
	m["testcookie"] = "1"
	
	req := NewRequest("POST","http://www.tourismjournal.net","/wp-login.php","http://www.tourismjournal.net/wp-login.php?2yniv1rfyct10sdlt1d1r",m,nil)
	resp,err := req.CallWithoutRedirect()

	if err != nil {
		fmt.Println("error ", err)
	}else {
		assert := assert.New(Test)
		assert.NotNil(resp)
		assert.Equal("http://www.tourismjournal.net/wp-login.php", resp.LastUrl)
		assert.NotNil(resp.Cookies)
		assert.NotNil(resp.RawContent)
	}
}
