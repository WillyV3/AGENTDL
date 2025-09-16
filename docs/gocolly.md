Go
Skip to Main Content
Search packages or symbols

Why Gosubmenu dropdown icon
Learn
Docssubmenu dropdown icon
Packages
Communitysubmenu dropdown icon
Notice  The highest tagged major version is v2.
Discover Packages
 
github.com/gocolly/colly

Go
colly
package
module


Main
Details
unchecked Valid go.mod file 
checked Redistributable license 
checked Tagged version 
checked Stable version 
Learn more about best practices
Repository
github.com/gocolly/colly
Links
Open Source Insights Logo Open Source Insights

NewCollector(options)
 README ¶
Colly
Lightning Fast and Elegant Scraping Framework for Gophers

Colly provides a clean interface to write any kind of crawler/scraper/spider.

With Colly you can easily extract structured data from websites, which can be used for a wide range of applications, like data mining, data processing or archiving.

GoDoc Backers on Open Collective Sponsors on Open Collective build status report card view examples Code Coverage FOSSA Status Twitter URL

Features
Clean API
Fast (>1k request/sec on a single core)
Manages request delays and maximum concurrency per domain
Automatic cookie and session handling
Sync/async/parallel scraping
Caching
Automatic encoding of non-unicode responses
Robots.txt support
Distributed scraping
Configuration via environment variables
Extensions
Example
func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}
See examples folder for more detailed examples.

Installation
go get -u github.com/gocolly/colly/...
Bugs
Bugs or suggestions? Visit the issue tracker or join #colly on freenode

Other Projects Using Colly
Below is a list of public, open source projects that use Colly:

greenpeace/check-my-pages Scraping script to test the Spanish Greenpeace web archive
altsab/gowap Wappalyzer implementation in Go
jesuiscamille/goquotes A quotes scrapper, making your day a little better!
jivesearch/jivesearch A search engine that doesn't track you.
Leagify/colly-draft-prospects A scraper for future NFL Draft prospects.
lucasepe/go-ps4 Search playstation store for your favorite PS4 games using the command line.
If you are using Colly in a project please send a pull request to add it to the list.

Contributors
This project exists thanks to all the people who contribute. [Contribute]. 

Backers
Thank you to all our backers! 🙏 [Become a backer]



Sponsors
Support this project by becoming a sponsor. Your logo will show up here with a link to your website. [Become a sponsor]

         

License
FOSSA Status

Expand ▾
 Documentation ¶
Overview ¶
Package colly implements a HTTP scraping framework

Index ¶
Constants
Variables
func AllowURLRevisit() func(*Collector)
func AllowedDomains(domains ...string) func(*Collector)
func Async(a bool) func(*Collector)
func CacheDir(path string) func(*Collector)
func Debugger(d debug.Debugger) func(*Collector)
func DetectCharset() func(*Collector)
func DisallowedDomains(domains ...string) func(*Collector)
func DisallowedURLFilters(filters ...*regexp.Regexp) func(*Collector)
func ID(id uint32) func(*Collector)
func IgnoreRobotsTxt() func(*Collector)
func MaxBodySize(sizeInBytes int) func(*Collector)
func MaxDepth(depth int) func(*Collector)
func ParseHTTPErrorResponse() func(*Collector)
func SanitizeFileName(fileName string) string
func URLFilters(filters ...*regexp.Regexp) func(*Collector)
func UnmarshalHTML(v interface{}, s *goquery.Selection) error
func UserAgent(ua string) func(*Collector)
type Collector
func NewCollector(options ...func(*Collector)) *Collector
func (c *Collector) Appengine(ctx context.Context)
func (c *Collector) Clone() *Collector
func (c *Collector) Cookies(URL string) []*http.Cookie
func (c *Collector) DisableCookies()
func (c *Collector) Head(URL string) error
func (c *Collector) Init()
func (c *Collector) Limit(rule *LimitRule) error
func (c *Collector) Limits(rules []*LimitRule) error
func (c *Collector) OnError(f ErrorCallback)
func (c *Collector) OnHTML(goquerySelector string, f HTMLCallback)
func (c *Collector) OnHTMLDetach(goquerySelector string)
func (c *Collector) OnRequest(f RequestCallback)
func (c *Collector) OnResponse(f ResponseCallback)
func (c *Collector) OnScraped(f ScrapedCallback)
func (c *Collector) OnXML(xpathQuery string, f XMLCallback)
func (c *Collector) OnXMLDetach(xpathQuery string)
func (c *Collector) Post(URL string, requestData map[string]string) error
func (c *Collector) PostMultipart(URL string, requestData map[string][]byte) error
func (c *Collector) PostRaw(URL string, requestData []byte) error
func (c *Collector) Request(method, URL string, requestData io.Reader, ctx *Context, hdr http.Header) error
func (c *Collector) SetCookieJar(j *cookiejar.Jar)
func (c *Collector) SetCookies(URL string, cookies []*http.Cookie) error
func (c *Collector) SetDebugger(d debug.Debugger)
func (c *Collector) SetProxy(proxyURL string) error
func (c *Collector) SetProxyFunc(p ProxyFunc)
func (c *Collector) SetRequestTimeout(timeout time.Duration)
func (c *Collector) SetStorage(s storage.Storage) error
func (c *Collector) String() string
func (c *Collector) UnmarshalRequest(r []byte) (*Request, error)
func (c *Collector) Visit(URL string) error
func (c *Collector) Wait()
func (c *Collector) WithTransport(transport http.RoundTripper)
type Context
func NewContext() *Context
func (c *Context) ForEach(fn func(k string, v interface{}) interface{}) []interface{}
func (c *Context) Get(key string) string
func (c *Context) GetAny(key string) interface{}
func (c *Context) MarshalBinary() (_ []byte, _ error)
func (c *Context) Put(key string, value interface{})
func (c *Context) UnmarshalBinary(_ []byte) error
type ErrorCallback
type HTMLCallback
type HTMLElement
func NewHTMLElementFromSelectionNode(resp *Response, s *goquery.Selection, n *html.Node, idx int) *HTMLElement
func (h *HTMLElement) Attr(k string) string
func (h *HTMLElement) ChildAttr(goquerySelector, attrName string) string
func (h *HTMLElement) ChildAttrs(goquerySelector, attrName string) []string
func (h *HTMLElement) ChildText(goquerySelector string) string
func (h *HTMLElement) ForEach(goquerySelector string, callback func(int, *HTMLElement))
func (h *HTMLElement) ForEachWithBreak(goquerySelector string, callback func(int, *HTMLElement) bool)
func (h *HTMLElement) Unmarshal(v interface{}) error
type LimitRule
func (r *LimitRule) Init() error
func (r *LimitRule) Match(domain string) bool
type ProxyFunc
type Request
func (r *Request) Abort()
func (r *Request) AbsoluteURL(u string) string
func (r *Request) Do() error
func (r *Request) Marshal() ([]byte, error)
func (r *Request) New(method, URL string, body io.Reader) (*Request, error)
func (r *Request) Post(URL string, requestData map[string]string) error
func (r *Request) PostMultipart(URL string, requestData map[string][]byte) error
func (r *Request) PostRaw(URL string, requestData []byte) error
func (r *Request) Retry() error
func (r *Request) Visit(URL string) error
type RequestCallback
type Response
func (r *Response) FileName() string
func (r *Response) Save(fileName string) error
type ResponseCallback
type ScrapedCallback
type XMLCallback
type XMLElement
func NewXMLElementFromHTMLNode(resp *Response, s *html.Node) *XMLElement
func NewXMLElementFromXMLNode(resp *Response, s *xmlquery.Node) *XMLElement
func (h *XMLElement) Attr(k string) string
func (h *XMLElement) ChildAttr(xpathQuery, attrName string) string
func (h *XMLElement) ChildAttrs(xpathQuery, attrName string) []string
func (h *XMLElement) ChildText(xpathQuery string) string
func (h *XMLElement) ChildTexts(xpathQuery string) []string
Constants ¶
View Source
const ProxyURLKey key = iota
ProxyURLKey is the context key for the request proxy address.

Variables ¶
View Source
var (
	// ErrForbiddenDomain is the error thrown if visiting
	// a domain which is not allowed in AllowedDomains
	ErrForbiddenDomain = errors.New("Forbidden domain")
	// ErrMissingURL is the error type for missing URL errors
	ErrMissingURL = errors.New("Missing URL")
	// ErrMaxDepth is the error type for exceeding max depth
	ErrMaxDepth = errors.New("Max depth limit reached")
	// ErrForbiddenURL is the error thrown if visiting
	// a URL which is not allowed by URLFilters
	ErrForbiddenURL = errors.New("ForbiddenURL")

	// ErrNoURLFiltersMatch is the error thrown if visiting
	// a URL which is not allowed by URLFilters
	ErrNoURLFiltersMatch = errors.New("No URLFilters match")
	// ErrAlreadyVisited is the error type for already visited URLs
	ErrAlreadyVisited = errors.New("URL already visited")
	// ErrRobotsTxtBlocked is the error type for robots.txt errors
	ErrRobotsTxtBlocked = errors.New("URL blocked by robots.txt")
	// ErrNoCookieJar is the error type for missing cookie jar
	ErrNoCookieJar = errors.New("Cookie jar is not available")
	// ErrNoPattern is the error type for LimitRules without patterns
	ErrNoPattern = errors.New("No pattern defined in LimitRule")
)
Functions ¶
func AllowURLRevisit ¶
func AllowURLRevisit() func(*Collector)
AllowURLRevisit instructs the Collector to allow multiple downloads of the same URL

func AllowedDomains ¶
func AllowedDomains(domains ...string) func(*Collector)
AllowedDomains sets the domain whitelist used by the Collector.

func Async ¶
func Async(a bool) func(*Collector)
Async turns on asynchronous network requests.

func CacheDir ¶
func CacheDir(path string) func(*Collector)
CacheDir specifies the location where GET requests are cached as files.

func Debugger ¶
func Debugger(d debug.Debugger) func(*Collector)
Debugger sets the debugger used by the Collector.

func DetectCharset ¶
func DetectCharset() func(*Collector)
DetectCharset enables character encoding detection for non-utf8 response bodies without explicit charset declaration. This feature uses https://github.com/saintfish/chardet

func DisallowedDomains ¶
func DisallowedDomains(domains ...string) func(*Collector)
DisallowedDomains sets the domain blacklist used by the Collector.

func DisallowedURLFilters ¶
func DisallowedURLFilters(filters ...*regexp.Regexp) func(*Collector)
DisallowedURLFilters sets the list of regular expressions which restricts visiting URLs. If any of the rules matches to a URL the request will be stopped.

func ID ¶
func ID(id uint32) func(*Collector)
ID sets the unique identifier of the Collector.

func IgnoreRobotsTxt ¶
func IgnoreRobotsTxt() func(*Collector)
IgnoreRobotsTxt instructs the Collector to ignore any restrictions set by the target host's robots.txt file.

func MaxBodySize ¶
func MaxBodySize(sizeInBytes int) func(*Collector)
MaxBodySize sets the limit of the retrieved response body in bytes.

func MaxDepth ¶
func MaxDepth(depth int) func(*Collector)
MaxDepth limits the recursion depth of visited URLs.

func ParseHTTPErrorResponse ¶
func ParseHTTPErrorResponse() func(*Collector)
ParseHTTPErrorResponse allows parsing responses with HTTP errors

func SanitizeFileName ¶
func SanitizeFileName(fileName string) string
SanitizeFileName replaces dangerous characters in a string so the return value can be used as a safe file name.

func URLFilters ¶
func URLFilters(filters ...*regexp.Regexp) func(*Collector)
URLFilters sets the list of regular expressions which restricts visiting URLs. If any of the rules matches to a URL the request won't be stopped.

func UnmarshalHTML ¶
func UnmarshalHTML(v interface{}, s *goquery.Selection) error
UnmarshalHTML declaratively extracts text or attributes to a struct from HTML response using struct tags composed of css selectors. Allowed struct tags:

"selector" (required): CSS (goquery) selector of the desired data
"attr" (optional): Selects the matching element's attribute's value. Leave it blank or omit to get the text of the element.
Example struct declaration:

type Nested struct {
	String  string   `selector:"div > p"`
   Classes []string `selector:"li" attr:"class"`
	Struct  *Nested  `selector:"div > div"`
}
Supported types: struct, *struct, string, []string

func UserAgent ¶
func UserAgent(ua string) func(*Collector)
UserAgent sets the user agent used by the Collector.

Types ¶
type Collector ¶
type Collector struct {
	// UserAgent is the User-Agent string used by HTTP requests
	UserAgent string
	// MaxDepth limits the recursion depth of visited URLs.
	// Set it to 0 for infinite recursion (default).
	MaxDepth int
	// AllowedDomains is a domain whitelist.
	// Leave it blank to allow any domains to be visited
	AllowedDomains []string
	// DisallowedDomains is a domain blacklist.
	DisallowedDomains []string
	// DisallowedURLFilters is a list of regular expressions which restricts
	// visiting URLs. If any of the rules matches to a URL the
	// request will be stopped. DisallowedURLFilters will
	// be evaluated before URLFilters
	// Leave it blank to allow any URLs to be visited
	DisallowedURLFilters []*regexp.Regexp

	// Leave it blank to allow any URLs to be visited
	URLFilters []*regexp.Regexp

	// AllowURLRevisit allows multiple downloads of the same URL
	AllowURLRevisit bool
	// MaxBodySize is the limit of the retrieved response body in bytes.
	// 0 means unlimited.
	// The default value for MaxBodySize is 10MB (10 * 1024 * 1024 bytes).
	MaxBodySize int
	// CacheDir specifies a location where GET requests are cached as files.
	// When it's not defined, caching is disabled.
	CacheDir string
	// IgnoreRobotsTxt allows the Collector to ignore any restrictions set by
	// the target host's robots.txt file.  See http://www.robotstxt.org/ for more
	// information.
	IgnoreRobotsTxt bool
	// Async turns on asynchronous network communication. Use Collector.Wait() to
	// be sure all requests have been finished.
	Async bool
	// ParseHTTPErrorResponse allows parsing HTTP responses with non 2xx status codes.
	// By default, Colly parses only successful HTTP responses. Set ParseHTTPErrorResponse
	// to true to enable it.
	ParseHTTPErrorResponse bool
	// ID is the unique identifier of a collector
	ID uint32
	// DetectCharset can enable character encoding detection for non-utf8 response bodies
	// without explicit charset declaration. This feature uses https://github.com/saintfish/chardet
	DetectCharset bool
	// RedirectHandler allows control on how a redirect will be managed
	RedirectHandler func(req *http.Request, via []*http.Request) error
	// CheckHead performs a HEAD request before every GET to pre-validate the response
	CheckHead bool
	// contains filtered or unexported fields
}
Collector provides the scraper instance for a scraping job

func NewCollector ¶
func NewCollector(options ...func(*Collector)) *Collector
NewCollector creates a new Collector instance with default configuration

func (*Collector) Appengine ¶
func (c *Collector) Appengine(ctx context.Context)
Appengine will replace the Collector's backend http.Client With an Http.Client that is provided by appengine/urlfetch This function should be used when the scraper is run on Google App Engine. Example:

func startScraper(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  c := colly.NewCollector()
  c.Appengine(ctx)
   ...
  c.Visit("https://google.ca")
}
func (*Collector) Clone ¶
func (c *Collector) Clone() *Collector
Clone creates an exact copy of a Collector without callbacks. HTTP backend, robots.txt cache and cookie jar are shared between collectors.

func (*Collector) Cookies ¶
func (c *Collector) Cookies(URL string) []*http.Cookie
Cookies returns the cookies to send in a request for the given URL.

func (*Collector) DisableCookies ¶
func (c *Collector) DisableCookies()
DisableCookies turns off cookie handling

func (*Collector) Head ¶
added in v1.2.0
func (c *Collector) Head(URL string) error
Head starts a collector job by creating a HEAD request.

func (*Collector) Init ¶
func (c *Collector) Init()
Init initializes the Collector's private variables and sets default configuration for the Collector

func (*Collector) Limit ¶
func (c *Collector) Limit(rule *LimitRule) error
Limit adds a new LimitRule to the collector

func (*Collector) Limits ¶
func (c *Collector) Limits(rules []*LimitRule) error
Limits adds new LimitRules to the collector

func (*Collector) OnError ¶
func (c *Collector) OnError(f ErrorCallback)
OnError registers a function. Function will be executed if an error occurs during the HTTP request.

func (*Collector) OnHTML ¶
func (c *Collector) OnHTML(goquerySelector string, f HTMLCallback)
OnHTML registers a function. Function will be executed on every HTML element matched by the GoQuery Selector parameter. GoQuery Selector is a selector used by https://github.com/PuerkitoBio/goquery

func (*Collector) OnHTMLDetach ¶
func (c *Collector) OnHTMLDetach(goquerySelector string)
OnHTMLDetach deregister a function. Function will not be execute after detached

func (*Collector) OnRequest ¶
func (c *Collector) OnRequest(f RequestCallback)
OnRequest registers a function. Function will be executed on every request made by the Collector

func (*Collector) OnResponse ¶
func (c *Collector) OnResponse(f ResponseCallback)
OnResponse registers a function. Function will be executed on every response

func (*Collector) OnScraped ¶
func (c *Collector) OnScraped(f ScrapedCallback)
OnScraped registers a function. Function will be executed after OnHTML, as a final part of the scraping.

func (*Collector) OnXML ¶
func (c *Collector) OnXML(xpathQuery string, f XMLCallback)
OnXML registers a function. Function will be executed on every XML element matched by the xpath Query parameter. xpath Query is used by https://github.com/antchfx/xmlquery

func (*Collector) OnXMLDetach ¶
func (c *Collector) OnXMLDetach(xpathQuery string)
OnXMLDetach deregister a function. Function will not be execute after detached

func (*Collector) Post ¶
func (c *Collector) Post(URL string, requestData map[string]string) error
Post starts a collector job by creating a POST request. Post also calls the previously provided callbacks

func (*Collector) PostMultipart ¶
func (c *Collector) PostMultipart(URL string, requestData map[string][]byte) error
PostMultipart starts a collector job by creating a Multipart POST request with raw binary data. PostMultipart also calls the previously provided callbacks

func (*Collector) PostRaw ¶
func (c *Collector) PostRaw(URL string, requestData []byte) error
PostRaw starts a collector job by creating a POST request with raw binary data. Post also calls the previously provided callbacks

func (*Collector) Request ¶
func (c *Collector) Request(method, URL string, requestData io.Reader, ctx *Context, hdr http.Header) error
Request starts a collector job by creating a custom HTTP request where method, context, headers and request data can be specified. Set requestData, ctx, hdr parameters to nil if you don't want to use them. Valid methods:

"GET"
"HEAD"
"POST"
"PUT"
"DELETE"
"PATCH"
"OPTIONS"
func (*Collector) SetCookieJar ¶
func (c *Collector) SetCookieJar(j *cookiejar.Jar)
SetCookieJar overrides the previously set cookie jar

func (*Collector) SetCookies ¶
func (c *Collector) SetCookies(URL string, cookies []*http.Cookie) error
SetCookies handles the receipt of the cookies in a reply for the given URL

func (*Collector) SetDebugger ¶
func (c *Collector) SetDebugger(d debug.Debugger)
SetDebugger attaches a debugger to the collector

func (*Collector) SetProxy ¶
func (c *Collector) SetProxy(proxyURL string) error
SetProxy sets a proxy for the collector. This method overrides the previously used http.Transport if the type of the transport is not http.RoundTripper. The proxy type is determined by the URL scheme. "http" and "socks5" are supported. If the scheme is empty, "http" is assumed.

func (*Collector) SetProxyFunc ¶
func (c *Collector) SetProxyFunc(p ProxyFunc)
SetProxyFunc sets a custom proxy setter/switcher function. See built-in ProxyFuncs for more details. This method overrides the previously used http.Transport if the type of the transport is not http.RoundTripper. The proxy type is determined by the URL scheme. "http" and "socks5" are supported. If the scheme is empty, "http" is assumed.

func (*Collector) SetRequestTimeout ¶
func (c *Collector) SetRequestTimeout(timeout time.Duration)
SetRequestTimeout overrides the default timeout (10 seconds) for this collector

func (*Collector) SetStorage ¶
func (c *Collector) SetStorage(s storage.Storage) error
SetStorage overrides the default in-memory storage. Storage stores scraping related data like cookies and visited urls

func (*Collector) String ¶
func (c *Collector) String() string
String is the text representation of the collector. It contains useful debug information about the collector's internals

func (*Collector) UnmarshalRequest ¶
func (c *Collector) UnmarshalRequest(r []byte) (*Request, error)
UnmarshalRequest creates a Request from serialized data

func (*Collector) Visit ¶
func (c *Collector) Visit(URL string) error
Visit starts Collector's collecting job by creating a request to the URL specified in parameter. Visit also calls the previously provided callbacks

func (*Collector) Wait ¶
func (c *Collector) Wait()
Wait returns when the collector jobs are finished

func (*Collector) WithTransport ¶
func (c *Collector) WithTransport(transport http.RoundTripper)
WithTransport allows you to set a custom http.RoundTripper (transport)

type Context ¶
type Context struct {
	// contains filtered or unexported fields
}
Context provides a tiny layer for passing data between callbacks

func NewContext ¶
func NewContext() *Context
NewContext initializes a new Context instance

func (*Context) ForEach ¶
func (c *Context) ForEach(fn func(k string, v interface{}) interface{}) []interface{}
ForEach iterate context

func (*Context) Get ¶
func (c *Context) Get(key string) string
Get retrieves a string value from Context. Get returns an empty string if key not found

func (*Context) GetAny ¶
func (c *Context) GetAny(key string) interface{}
GetAny retrieves a value from Context. GetAny returns nil if key not found

func (*Context) MarshalBinary ¶
func (c *Context) MarshalBinary() (_ []byte, _ error)
MarshalBinary encodes Context value This function is used by request caching

func (*Context) Put ¶
func (c *Context) Put(key string, value interface{})
Put stores a value of any type in Context

func (*Context) UnmarshalBinary ¶
func (c *Context) UnmarshalBinary(_ []byte) error
UnmarshalBinary decodes Context value to nil This function is used by request caching

type ErrorCallback ¶
type ErrorCallback func(*Response, error)
ErrorCallback is a type alias for OnError callback functions

type HTMLCallback ¶
type HTMLCallback func(*HTMLElement)
HTMLCallback is a type alias for OnHTML callback functions

type HTMLElement ¶
type HTMLElement struct {
	// Name is the name of the tag
	Name string
	Text string

	// Request is the request object of the element's HTML document
	Request *Request
	// Response is the Response object of the element's HTML document
	Response *Response
	// DOM is the goquery parsed DOM object of the page. DOM is relative
	// to the current HTMLElement
	DOM *goquery.Selection
	// Index stores the position of the current element within all the elements matched by an OnHTML callback
	Index int
	// contains filtered or unexported fields
}
HTMLElement is the representation of a HTML tag.

func NewHTMLElementFromSelectionNode ¶
func NewHTMLElementFromSelectionNode(resp *Response, s *goquery.Selection, n *html.Node, idx int) *HTMLElement
NewHTMLElementFromSelectionNode creates a HTMLElement from a goquery.Selection Node.

func (*HTMLElement) Attr ¶
func (h *HTMLElement) Attr(k string) string
Attr returns the selected attribute of a HTMLElement or empty string if no attribute found

func (*HTMLElement) ChildAttr ¶
func (h *HTMLElement) ChildAttr(goquerySelector, attrName string) string
ChildAttr returns the stripped text content of the first matching element's attribute.

func (*HTMLElement) ChildAttrs ¶
func (h *HTMLElement) ChildAttrs(goquerySelector, attrName string) []string
ChildAttrs returns the stripped text content of all the matching element's attributes.

func (*HTMLElement) ChildText ¶
func (h *HTMLElement) ChildText(goquerySelector string) string
ChildText returns the concatenated and stripped text content of the matching elements.

func (*HTMLElement) ForEach ¶
func (h *HTMLElement) ForEach(goquerySelector string, callback func(int, *HTMLElement))
ForEach iterates over the elements matched by the first argument and calls the callback function on every HTMLElement match.

func (*HTMLElement) ForEachWithBreak ¶
added in v1.1.0
func (h *HTMLElement) ForEachWithBreak(goquerySelector string, callback func(int, *HTMLElement) bool)
ForEachWithBreak iterates over the elements matched by the first argument and calls the callback function on every HTMLElement match. It is identical to ForEach except that it is possible to break out of the loop by returning false in the callback function. It returns the current Selection object.

func (*HTMLElement) Unmarshal ¶
func (h *HTMLElement) Unmarshal(v interface{}) error
Unmarshal is a shorthand for colly.UnmarshalHTML

type LimitRule ¶
type LimitRule struct {
	// DomainRegexp is a regular expression to match against domains
	DomainRegexp string
	// DomainRegexp is a glob pattern to match against domains
	DomainGlob string
	// Delay is the duration to wait before creating a new request to the matching domains
	Delay time.Duration
	// RandomDelay is the extra randomized duration to wait added to Delay before creating a new request
	RandomDelay time.Duration
	// Parallelism is the number of the maximum allowed concurrent requests of the matching domains
	Parallelism int
	// contains filtered or unexported fields
}
LimitRule provides connection restrictions for domains. Both DomainRegexp and DomainGlob can be used to specify the included domains patterns, but at least one is required. There can be two kind of limitations:

Parallelism: Set limit for the number of concurrent requests to matching domains
Delay: Wait specified amount of time between requests (parallelism is 1 in this case)
func (*LimitRule) Init ¶
func (r *LimitRule) Init() error
Init initializes the private members of LimitRule

func (*LimitRule) Match ¶
func (r *LimitRule) Match(domain string) bool
Match checks that the domain parameter triggers the rule

type ProxyFunc ¶
type ProxyFunc func(*http.Request) (*url.URL, error)
ProxyFunc is a type alias for proxy setter functions.

type Request ¶
type Request struct {
	// URL is the parsed URL of the HTTP request
	URL *url.URL
	// Headers contains the Request's HTTP headers
	Headers *http.Header
	// Ctx is a context between a Request and a Response
	Ctx *Context
	// Depth is the number of the parents of the request
	Depth int
	// Method is the HTTP method of the request
	Method string
	// Body is the request body which is used on POST/PUT requests
	Body io.Reader
	// ResponseCharacterencoding is the character encoding of the response body.
	// Leave it blank to allow automatic character encoding of the response body.
	// It is empty by default and it can be set in OnRequest callback.
	ResponseCharacterEncoding string
	// ID is the Unique identifier of the request
	ID uint32

	// ProxyURL is the proxy address that handles the request
	ProxyURL string
	// contains filtered or unexported fields
}
Request is the representation of a HTTP request made by a Collector

func (*Request) Abort ¶
func (r *Request) Abort()
Abort cancels the HTTP request when called in an OnRequest callback

func (*Request) AbsoluteURL ¶
func (r *Request) AbsoluteURL(u string) string
AbsoluteURL returns with the resolved absolute URL of an URL chunk. AbsoluteURL returns empty string if the URL chunk is a fragment or could not be parsed

func (*Request) Do ¶
func (r *Request) Do() error
Do submits the request

func (*Request) Marshal ¶
func (r *Request) Marshal() ([]byte, error)
Marshal serializes the Request

func (*Request) New ¶
func (r *Request) New(method, URL string, body io.Reader) (*Request, error)
New creates a new request with the context of the original request

func (*Request) Post ¶
func (r *Request) Post(URL string, requestData map[string]string) error
Post continues a collector job by creating a POST request and preserves the Context of the previous request. Post also calls the previously provided callbacks

func (*Request) PostMultipart ¶
func (r *Request) PostMultipart(URL string, requestData map[string][]byte) error
PostMultipart starts a collector job by creating a Multipart POST request with raw binary data. PostMultipart also calls the previously provided. callbacks

func (*Request) PostRaw ¶
func (r *Request) PostRaw(URL string, requestData []byte) error
PostRaw starts a collector job by creating a POST request with raw binary data. PostRaw preserves the Context of the previous request and calls the previously provided callbacks

func (*Request) Retry ¶
func (r *Request) Retry() error
Retry submits HTTP request again with the same parameters

func (*Request) Visit ¶
func (r *Request) Visit(URL string) error
Visit continues Collector's collecting job by creating a request and preserves the Context of the previous request. Visit also calls the previously provided callbacks

type RequestCallback ¶
type RequestCallback func(*Request)
RequestCallback is a type alias for OnRequest callback functions

type Response ¶
type Response struct {
	// StatusCode is the status code of the Response
	StatusCode int
	// Body is the content of the Response
	Body []byte
	// Ctx is a context between a Request and a Response
	Ctx *Context
	// Request is the Request object of the response
	Request *Request
	// Headers contains the Response's HTTP headers
	Headers *http.Header
}
Response is the representation of a HTTP response made by a Collector

func (*Response) FileName ¶
func (r *Response) FileName() string
FileName returns the sanitized file name parsed from "Content-Disposition" header or from URL

func (*Response) Save ¶
func (r *Response) Save(fileName string) error
Save writes response body to disk

type ResponseCallback ¶
type ResponseCallback func(*Response)
ResponseCallback is a type alias for OnResponse callback functions

type ScrapedCallback ¶
type ScrapedCallback func(*Response)
ScrapedCallback is a type alias for OnScraped callback functions

type XMLCallback ¶
type XMLCallback func(*XMLElement)
XMLCallback is a type alias for OnXML callback functions

type XMLElement ¶
type XMLElement struct {
	// Name is the name of the tag
	Name string
	Text string

	// Request is the request object of the element's HTML document
	Request *Request
	// Response is the Response object of the element's HTML document
	Response *Response
	// DOM is the DOM object of the page. DOM is relative
	// to the current XMLElement and is either a html.Node or xmlquery.Node
	// based on how the XMLElement was created.
	DOM interface{}
	// contains filtered or unexported fields
}
XMLElement is the representation of a XML tag.

func NewXMLElementFromHTMLNode ¶
func NewXMLElementFromHTMLNode(resp *Response, s *html.Node) *XMLElement
NewXMLElementFromHTMLNode creates a XMLElement from a html.Node.

func NewXMLElementFromXMLNode ¶
func NewXMLElementFromXMLNode(resp *Response, s *xmlquery.Node) *XMLElement
NewXMLElementFromXMLNode creates a XMLElement from a xmlquery.Node.

func (*XMLElement) Attr ¶
func (h *XMLElement) Attr(k string) string
Attr returns the selected attribute of a HTMLElement or empty string if no attribute found

func (*XMLElement) ChildAttr ¶
func (h *XMLElement) ChildAttr(xpathQuery, attrName string) string
ChildAttr returns the stripped text content of the first matching element's attribute.

func (*XMLElement) ChildAttrs ¶
func (h *XMLElement) ChildAttrs(xpathQuery, attrName string) []string
ChildAttrs returns the stripped text content of all the matching element's attributes.

func (*XMLElement) ChildText ¶
func (h *XMLElement) ChildText(xpathQuery string) string
ChildText returns the concatenated and stripped text content of the matching elements.

func (*XMLElement) ChildTexts ¶
added in v1.1.0
func (h *XMLElement) ChildTexts(xpathQuery string) []string
ChildTexts returns an array of strings corresponding to child elements that match the xpath query. Each item in the array is the stripped text content of the corresponding matching child element.

 Source Files ¶
View all Source files
colly.go
context.go
htmlelement.go
http_backend.go
request.go
response.go
unmarshal.go
xmlelement.go
 Directories ¶
Expand all
_examples
cmd
debug
extensions
Package extensions implements various helper addons for Colly
proxy
queue
storage
Why Go
Use Cases
Case Studies
Get Started
Playground
Tour
Stack Overflow
Help
Packages
Standard Library
Sub-repositories
About Go Packages
About
Download
Blog
Issue Tracker
Release Notes
Brand Guidelines
Code of Conduct
Connect
Twitter
GitHub
Slack
r/golang
Meetup
Golang Weekly
Gopher in flight goggles
Copyright
Terms of Service
Privacy Policy
Report an Issue
System theme
Theme Toggle


Shortcuts Modal

Google logo
go.dev uses cookies from Google to deliver and enhance the quality of its services and to analyze traffic. Learn more.
Okay