package goscreenscraper

import (
	"log"
	"net/http"
	"strings"
	"net/url"
	"io"
	"io/ioutil"
	"strconv"
	"net/http/cookiejar"
	"compress/gzip"
)

//constant for method post and get
const (
	METHOD_POST = "POST"
	METHOD_GET = "GET"
)

/**
 * HTTP Request will be used for creating request
 * method : method POST or GET
 * postContents : contents that will be sent through request
 * baseUrl : base url of site (ex : http://www.yoursite.com) . It must be start with http://
 * referer : the previous url
 * cookies : response cookies from previous page/login page  
**/
type HTTPRequest struct {
	method    string
	targetUrl string
	postContents map[string]string
	baseUrl string
	referer string
	cookies []*http.Cookie 
}

/**
 * HTTP Response 
 * RawContent : array of byte from resp.Body
 * lastUrl : landing page(redirect page)
 * Cookies : cookies from response
**/
type HttpReponse struct {
	RawContent []byte
	LastUrl	string
	Cookies []*http.Cookie
}

//initialise request
func NewRequest(method string, baseUrl string,targetUrl string, referer string, formData map[string]string, cookies []*http.Cookie) *HTTPRequest {
	req := new(HTTPRequest)
	req.method = method
	req.targetUrl = targetUrl
	req.referer = referer
	req.postContents = formData
	req.baseUrl = baseUrl
	req.cookies = cookies

	return req
}

//initialise response
func NewResponse() *HttpReponse {
	w := new(HttpReponse)
	return w
}

/**
 * CallWithoutRedirect method using roundtrip transport for sending data
 * This method doesn't allow cookies when sending data
 * It's also doesn't support redirect. Plese refer doc of http.Transport roundtrip
 * I usually use this method to get cookie after successfull login
 **/
func (httpRequest *HTTPRequest) CallWithoutRedirect() (*HttpReponse,error) {
	fixedUrl := getFullUrl(httpRequest.baseUrl ,httpRequest.targetUrl)
	
	httpRequest.targetUrl = fixedUrl

	var formData url.Values

	if httpRequest.postContents != nil {
		formData = generateFormData(httpRequest.postContents)
	}

	reqConn,err := httpRequest.createRequestConn(formData)
	
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client := &http.Transport{}
	resp, err := client.RoundTrip(reqConn)
	
	if err != nil {
		log.Println(err)
        return nil,err
    } 
    
    defer resp.Body.Close()

    parsedResponse,err := httpRequest.parseResponse(resp)

    if err != nil {
    	log.Println(err)
   		return nil,err 	
    }
	
	return parsedResponse, nil
}

//populate HTTResponse
func (httpRequest *HTTPRequest) parseResponse(resp *http.Response) (*HttpReponse,error){
	response := NewResponse()
	
	var contents []byte
	var er error

	if resp.Header.Get("Content-Encoding") == "gzip"{
		greader,err := gzip.NewReader(resp.Body)
        
		if err != nil {
			log.Println(err)
			return nil,err
		}

        contents,er = ioutil.ReadAll(greader)	

        if er != nil {
        	log.Println(err)
        	return nil,er
        }

		defer greader.Close()
	}else{
		contents,er = ioutil.ReadAll(resp.Body)

		if er != nil {
			log.Println(er)
			return nil,er
		}
	}

	response.RawContent = contents
	response.Cookies = resp.Cookies()
	response.LastUrl = resp.Request.URL.String()

	return response,nil 
}

/**
 * CallAndRedirect use client for send request. 
 * It is implement cookieJar so if you want to send request using cookies, Use this method
 * It is support redirect page
 **/
func (httpRequest *HTTPRequest) CallAndRedirect() (*HttpReponse,error) {
	fixedUrl := getFullUrl(httpRequest.baseUrl ,httpRequest.targetUrl)
	
	httpRequest.targetUrl = fixedUrl

	var formData url.Values

	if httpRequest.postContents != nil {
		formData = generateFormData(httpRequest.postContents)
	}

	reqConn,err := httpRequest.createRequestConn(formData)

	if err != nil {
		log.Println("error : ", err)
		return nil,err
	}

	jar, _ := cookiejar.New(nil)

	if httpRequest.cookies != nil{
		jar.SetCookies(reqConn.URL, httpRequest.cookies)
	}

	client := &http.Client{
		Jar : jar,
	}
	
	resp,err := client.Do(reqConn)	

	if err != nil {
		log.Printf("%s", err)
        return nil,err
    } 
        
    defer resp.Body.Close()

	parsedResponse,err := httpRequest.parseResponse(resp)
    
	if err != nil {
		log.Println(err)
   		return nil,err 	
    }
	
	return parsedResponse, nil
}

//get full url from baseUrl + targetUrl
func getFullUrl(baseUrl string, targetUrl string) string {
	startChar := targetUrl[0:1]

	if startChar == "/" {
		return baseUrl + targetUrl
	} else {
		return targetUrl
	}
}

//create httpRequest
func (httpRequest *HTTPRequest) createRequestConn(formData url.Values) (*http.Request,error) {
	var reader io.Reader

	formDataStr := ""
	generatedUrl := httpRequest.targetUrl

	reader = nil
	
	if httpRequest.method == METHOD_POST && formData != nil{
		formDataStr = formData.Encode()
		reader = strings.NewReader(formDataStr)
	}else if httpRequest.method == METHOD_GET && formData != nil{
		baseUrl,_ := url.Parse(httpRequest.targetUrl)

		baseUrl.RawQuery = formData.Encode()

		generatedUrl = baseUrl.String()
	}

	req, err := http.NewRequest(httpRequest.method, generatedUrl, reader)

	if err != nil {
		log.Printf("%s", err)
		return nil,err
	}else{
		httpRequest.setRequestHeader(req,formDataStr)
	}

	return req,nil
}

//set request header
func(httpRequest *HTTPRequest) setRequestHeader(request *http.Request,formDataStr string){
	request.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 6.1; rv:11.0) Gecko/20100101 Firefox/11.0")
	request.Header.Add("Accept-Encoding", "gzip, deflate")

	var newHost string

	newHost = strings.Replace(httpRequest.baseUrl,"http://" , "", -1)
	newHost = strings.Replace(newHost,"https://" , "", -1)

	request.Header.Add("Host", newHost)
	request.Header.Add("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Language", "en-us,en;q=0.5")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Set("Content-Length", strconv.Itoa(len(formDataStr)))

	if  httpRequest.referer != "" {
		request.Header.Set("Referer", httpRequest.referer)
	}

	request.Header.Set("Content-Type", "text/html; charset=UTF-8")

	if METHOD_POST == httpRequest.method {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} 
}

//populate form data value 
func generateFormData(values map[string]string) url.Values{
	parameters := url.Values{}

	for k,v := range values{
		parameters.Add(k, v)	
	}

	return parameters
}