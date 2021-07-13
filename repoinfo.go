package main

import(
	"os"
	"io/ioutil"
	"net/http"
	"fmt"
	"bytes"
	"log"
	"github.com/PuerkitoBio/goquery"
	"github.com/jessevdk/go-flags"
	"github.com/buger/jsonparser"
)

var (
	fork, steal, count, repo, max int
	m map[string]int = map[string]int{"": 0,}
	author string
	opts struct {
		Single string `short:"s" long:"single" description:"Show single repo info."`
		Author string `short:"a" long:"author"`
		Repo int `short:"r" long:"repo" description:"Repo num to check from top."`
	}
)

func GetBody(site string)[]byte{
	res, err := http.Get(site)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
func IsForkRepo(site string)bool {
	reader := bytes.NewReader(GetBody(site))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	keywords, _ := doc.Find("meta[name=\"octolytics-dimension-repository_is_fork\"]").Attr("content")
	if keywords[0]=='t' {
		fmt.Println("The repo is fork.")
		return true
	}
	fmt.Println("The repo is not fork.")
	return false
}

func ParseAuthorInfo(){
	for num := 1; num <= max; num++{
		site := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", author, num)
		body := GetBody(site)
		for i := 0; i < 100; i++{
			isfork, err := jsonparser.GetBoolean(body,fmt.Sprintf("[%d]",i),"fork")
			if err == nil {
				ctime, _ := jsonparser.GetString(body,fmt.Sprintf("[%d]",i),"created_at")
				utime, _ := jsonparser.GetString(body,fmt.Sprintf("[%d]",i),"updated_at")
				codetype, _ := jsonparser.GetString(body,fmt.Sprintf("[%d]",i),"language")
				if isfork {
					fork++
					if ctime == utime {steal++}
				}else{
					m[codetype]++
				}
				count++
			}
		}
	}
}

func PrintInfo(){
	fmt.Println("Total repo from",author,":",count)
	fmt.Println("|____Fork repo:", fork)
	fmt.Println("|    |____Non modified fork:", steal)
	fmt.Println("|____Non fork repo:", count - fork)
	fmt.Println("     |____Noncode repo:",m[""])
	for key, value := range m{
		if key != "" && value != 0{
			fmt.Println("         ",key,"repo:",value)
		}
	}
}

func main() {

	_,err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}
	if opts.Single != "" {
		IsForkRepo(opts.Single)
	} else if opts.Author != "" && opts.Repo != 0 {
		author, repo = opts.Author, opts.Repo
		max = repo/100 + 1
		ParseAuthorInfo()
		PrintInfo()
	}
}
