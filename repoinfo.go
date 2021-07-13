package main

import(
	"io/ioutil"
	"net/http"
	"fmt"
	"bytes"
	"log"
	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
)

func IsForkRepo(site string)bool {
	res, err := http.Get(site)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	keywords, _ := doc.Find("meta[name=\"octolytics-dimension-repository_is_fork\"]").Attr("content")
	if keywords[0]=='t'{return true}
	return false
}
func main() {
	/*if IsForkRepo("https://github.com/audreyt/q-moedict"){
		fmt.Println("is fork!")
	}else{
		fmt.Println("not fork!")
	}*/
	var fork, steal, count int
	var m map[string]int = map[string]int{"": 0,}
	author := "audreyt"
	repo := 331
	max := repo/100 + 1
	for num := 1; num <= max; num++{
		res, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", author, num))
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
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
