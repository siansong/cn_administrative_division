package main

import (
	"fmt"
	"net/url"
	"sync"

	// "log"
	"net/http" //http utils
	// "runtime"
	"strconv"  // string convert pkg
	"strings"
	"time"
	

	"github.com/PuerkitoBio/goquery" // html文档解析
	"github.com/djimenez/iconv-go"   // 字符编码间转换
	"github.com/go-pg/pg"            // postgresql orm
	"github.com/go-pg/pg/orm"        // 同上
)

// BaseURL ...
const BaseURL = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2017/"

const ignoreCode = "[ignore]"

// Level .
type Level string

const (
	// LVProvince ...
	LVProvince Level = "省"
	// LVCity ...
	LVCity Level = "市"
	// LVCounty ...
	LVCounty Level = "区/县"
	// LVTown ...
	LVTown Level = "街道/乡镇"
	// LVVillage ...
	LVVillage Level = "社区/村"
)

// Area 行政区划
type Area struct {
	Code       string `sql:",pk"` // 编码
	Name       string // 名字
	Href       string // html链接
	ParentCode string // 父级code
	Level      Level  // 级别
	// Path       string // eg:10/1001/100103 差不多自己体会吧。。。
	// FullName   string // 全路径名字
	//todo path,full name
}

// parseRelativeHref2BaseUrl 解析出相对于BaseUrl的相对路径
func parseRelativeHref2BaseURL(href string, refererURL string) string {
	if strings.Index(href, "/") == 0 {
		panic("应该没有这种情况吧，有的话看下")
	}

	relativeBase := refererURL[0:strings.LastIndex(refererURL, "/")]
	fullHref := relativeBase + "/" + href
	start := strings.LastIndex(fullHref, BaseURL) + len(BaseURL)
	relativeHref := fullHref[start:]
	if strings.Index(relativeHref, "/") == 0 {
		relativeHref = relativeHref[1:]
	}
	return relativeHref
}

// 解析为Area
type areaParser func(selection *goquery.Selection, refererUrl string) *Area

var provinceParser areaParser = func(selection *goquery.Selection, refererUrl string) *Area {
	a := selection.Find("a")
	href, exists := a.Attr("href")
	if !exists {
		html, _ := selection.Html()
		fmt.Println("不存在href属性的链接", html)
		return nil
	}
	name := gbk2utf8(a.Text())
	// html, _ := selection.Html()
	// fmt.Println(html)

	code := href[0:strings.LastIndex(href, ".")]
	return &Area{
		Name:  name,
		Code:  code,
		Href:  parseRelativeHref2BaseURL(href, refererUrl),
		Level: LVProvince,
	}
}

// city county town parser
var cityCountyTownParser areaParser = func(selection *goquery.Selection, refererUrl string) *Area {
	area := Area{}
	selection.Find("td").Each(func(i int, se *goquery.Selection) {
		if i == 0 {
			area.Code = se.Text()
		} else if i == 1 {
			area.Name = gbk2utf8(se.Text())

			aSelection := se.Find("a")
			if aSelection.Length() > 0 {
				href, _ := aSelection.Attr("href")
				area.Href = parseRelativeHref2BaseURL(href, refererUrl)
			}
		}
	})
	return &area
}

var villageParser areaParser = func(selection *goquery.Selection, refererUrl string) *Area {
	area := Area{}
	selection.Find("td").Each(func(i int, se *goquery.Selection) {
		if i == 0 {
			area.Code = se.Text()
		} else if i == 2 {
			area.Name = gbk2utf8(se.Text())
		}
	})
	return &area
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func panicIfWithMsg(err error, msg string) {
	if err != nil {
		fmt.Println("Err msg:", msg)
		panic(err)
	}
}

func ifThenPanic(flag bool, errMsg string) {
	if flag {
		panic(errMsg)
	}
}

func areaExistsInArray(area *Area, arr *[]Area) bool {
	for _, old := range *arr {
		if old.Code == area.Code {
			return true
		}
	}
	return false
}

func isStringBlank(str string) bool {
	return &str == nil || len(strings.Trim(str, " ")) == 0
}

// connect and try create table
func connTryCreateTb() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "pg",
		Password: "123",
		Database: "cad",
		PoolSize: 25, //runtime.NumCPU(): 8 (my pc)
	})

	// db.DropTable(&Area{}, &orm.DropTableOptions{
	// 	IfExists: true,
	// })

	db.CreateTable(&Area{}, &orm.CreateTableOptions{
		IfNotExists: true,
	})

	return db
}

// gbk2utf8 gbk字符串转utf8
func gbk2utf8(text string) string {
	gbkBytes := []byte(text)
	// fmt.Println("gbk len:", len(gbkBytes))
	utf8Val := make([]byte, len(gbkBytes)/2*3) // gbk中文2字节, utf8中文一般3字节，这个要考证下 XXX
	iconv.Convert(gbkBytes, utf8Val, "gbk", "utf-8")
	return string(utf8Val)
}

// func requestUrlRetryIfNone200(href string)

// 解析出Code 和 Name
func pageSpider(href string, selector string, parser areaParser) (areas []Area) {
	startT := time.Now()
	url, _ := url.Parse(href)

	var resp *http.Response
	var httpErr error

	for i := 0; i < 5; i++ {
		tmpResp, e := http.DefaultClient.Do(&http.Request{
			Method: "GET",
			URL:    url,
			Header: map[string][]string{
				// "Cookie": []string{"AD_RS_COOKIE=20081945"},
				"Host": []string{"www.stats.gov.cn"},
			},
		})
		httpErr = e
		if e != nil {
			continue
		}

		resp = tmpResp

		if tmpResp.StatusCode == 200 {
			defer tmpResp.Body.Close()
			break
		} else {
			tmpResp.Body.Close()
		}

	}
	panicIf(httpErr)
	ifThenPanic(resp.StatusCode != 200, "【ERR】repCode != 200, Got:"+strconv.Itoa(resp.StatusCode)+", url: "+href)
	// fmt.Println("http time:", time.Now().Sub(startT))
	fmt.Println("[http done]", time.Now().Sub(startT), href)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	panicIfWithMsg(err, href)
	doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
		area := parser(selection, href)
		if area != nil {
			areas = append(areas, *area)
		}
	})
	return
}

func findByLevelAndParent(lv Level, parentCode string) []Area {
	var areas []Area
	db := connTryCreateTb()
	defer db.Close()
	q := db.Model(&areas).
		Where("level = ?", lv)
	if parentCode != ignoreCode {
		q.Where("parent_code = ? ", parentCode)
	}
	q.Select()
	return areas
}

// fetchAllProvinces ...
func fetchAllProvinces() (areas []Area) {
	db := connTryCreateTb()
	defer db.Close()

	nowDbProvinces := findByLevelAndParent(LVProvince, ignoreCode)

	spiderOutput := pageSpider(BaseURL+"/index.html", "table.provincetable tr.provincetr td", provinceParser)

	var newProvinces []Area
	for _, area := range spiderOutput {
		if areaExistsInArray(&area, &nowDbProvinces) {
			continue
		}
		area.Level = LVProvince
		newProvinces = append(newProvinces, area)
	}

	if len(newProvinces) > 0 {
		err := db.Insert(&newProvinces)
		panicIf(err)
	}
	return
}

func fetchAllCities() {
	db := connTryCreateTb()
	defer db.Close()

	nowDbCities := findByLevelAndParent(LVCity, ignoreCode)

	startT := time.Now()
	var wg sync.WaitGroup

	provinces := findByLevelAndParent(LVProvince, ignoreCode)
	// fmt.Println("db datas", provinces)

	for _, p := range provinces {
		wg.Add(1)
		go func(p Area) {
			fetchCitiesOfProvince(p, db, nowDbCities)
			wg.Done()
		}(p)
	}
	wg.Wait()
	fmt.Println("[Fetch cities done]", time.Now().Sub(startT))
}

func fetchCitiesOfProvince(province Area, db *pg.DB, nowDbCities []Area) {
	fmt.Println("[Fetch city]", province)
	spiderOutput := pageSpider(BaseURL+province.Href, "table.citytable tr.citytr", cityCountyTownParser)

	var datas []Area
	for _, a := range spiderOutput {
		if areaExistsInArray(&a, &nowDbCities) {
			continue
		}
		a.Level = LVCity
		a.ParentCode = province.Code
		datas = append(datas, a)
	}

	if len(datas) > 0 {
		err := db.Insert(&datas)
		panicIf(err)
	}
}

func fetchAllCounties() {
	db := connTryCreateTb()
	defer db.Close()

	cities := findByLevelAndParent(LVCity, ignoreCode)

	if len(cities) == 0 {
		panic("没有City数据，先抓取City数据")
	}

	var wg sync.WaitGroup
	for _, city := range cities {
		wg.Add(1)
		go func(city Area) {
			fetchCountiesOfCity(city, db)
			wg.Done()
		}(city)
	}
	wg.Wait()
}

// 抓取某个City下的所有County
func fetchCountiesOfCity(city Area, db *pg.DB) {
	startT := time.Now()
	nowDbCountiesOfCity := findByLevelAndParent(LVCounty, city.Code)

	spiderOutput := pageSpider(BaseURL+city.Href, "table.countytable tr.countytr", cityCountyTownParser)

	var newCounties []Area
	for _, a := range spiderOutput {
		if areaExistsInArray(&a, &nowDbCountiesOfCity) {
			continue
		}
		a.Level = LVCounty
		a.ParentCode = city.Code
		newCounties = append(newCounties, a)
	}

	// fmt.Println("counties:", city, newCounties)
	if len(newCounties) > 0 {
		err := db.Insert(&newCounties)
		panicIf(err)
	}
	fmt.Println("[Fetch counties of city] done,", time.Now().Sub(startT), city)
}

func fetchAllTowns() {
	db := connTryCreateTb()
	defer db.Close()

	counties := findByLevelAndParent(LVCounty, ignoreCode) // 约3000+rows

	if len(counties) == 0 {
		panic("没有County数据，先抓取County数据")
	}

	var wg sync.WaitGroup
	for _, county := range counties {
		wg.Add(1)
		go func(ct Area) {
			fetchTownsOfCounty(ct, db)
			wg.Done()
		}(county)
	}
	wg.Wait()
}

func fetchTownsOfCounty(county Area, db *pg.DB) {
	if isStringBlank(county.Href) {
		fmt.Println("!!!county href is blank, skiped", county)
		return
	}

	startT := time.Now()
	nowDbTownsOfCounty := findByLevelAndParent(LVTown, county.Code)

	spiderOutput := pageSpider(BaseURL+county.Href, "table.towntable tr.towntr", cityCountyTownParser)

	var newTowns []Area
	for _, a := range spiderOutput {
		if areaExistsInArray(&a, &nowDbTownsOfCounty) {
			continue
		}
		a.Level = LVTown
		a.ParentCode = county.Code
		newTowns = append(newTowns, a)
	}

	// fmt.Println("counties:", city, newTowns)
	if len(newTowns) > 0 {
		err := db.Insert(&newTowns)
		panicIf(err)
	}
	fmt.Println("[Fetch towns of county] done,", time.Now().Sub(startT), county)
}

func fetchVillagesOfTown(town Area, db *pg.DB) {
	if isStringBlank(town.Href) {
		fmt.Println("!!!town href is blank, skiped", town)
		return
	}

	startT := time.Now()
	// nowDbVillagesOfTown := findByLevelAndParent(LVVillage, town.Code)

	spiderOutput := pageSpider(BaseURL+town.Href, "table.villagetable tr.villagetr", villageParser)

	var newVillages []Area
	for _, a := range spiderOutput {
		// if areaExistsInArray(&a, &nowDbVillagesOfTown) {
		// 	continue
		// }
		a.Level = LVVillage
		a.ParentCode = town.Code
		newVillages = append(newVillages, a)
	}

	// fmt.Println("counties:", city, newTowns)
	if len(newVillages) > 0 {
		err := db.Insert(&newVillages)
		if err != nil {
			var codes []string
			for _, v := range newVillages {
				codes = append(codes, v.Code)
			}
			fmt.Println("codes:", codes)
			panic(err)
		}
	}
	fmt.Println("[Fetch villages of town] done,", time.Now().Sub(startT), town)
}

func fetchAllVillages() {
	// 怎么搞呢？？
	// 这样吧，反正goroutines也NB，不如就老办法，一个乡镇开一个，先试试，初步估算下差不多会有4w+个页面吧，估计行数60w左右。。。
	// 数据规模：province:31, city:312, county:3000+, town:4w+, city/province: 10, county/province:100, town/city:120
	// counties(k) >> towns(m) >> villages  : k
	// [dep]city > county > town > village
	// [dep]100*400
	// counties > group by city 

	startT := time.Now()
	fmt.Println("[Fetch all villages start]")
	
	db := connTryCreateTb()
	defer db.Close()

	counties := findByLevelAndParent(LVCounty, ignoreCode) // 约3000+rows
	cityCountyMap := make(map[string][]Area) // key:cityCode, val: counties array
	if len(counties) == 0 {
		panic("没有County数据，先抓取County数据")
	}
	for _, county := range counties {
		cityCountyMap[county.ParentCode] = append(cityCountyMap[county.ParentCode], county)
	}

	var wg sync.WaitGroup
	for cityCode, countiesOfCity := range cityCountyMap {
		fmt.Println("[Fetch villages of city] ", cityCode)

		var countieIds []string
		for _, county := range countiesOfCity {
			countieIds = append(countieIds, county.Code)
		}

		var townsOfCity []Area
		// Where("id in (?)", pg.In(ids)).
		err := db.Model(&townsOfCity).
			Where("parent_code in (?)", pg.In(countieIds)).
			Select()
		panicIf(err)

		for _, town := range townsOfCity {
			wg.Add(1)
			go func(t Area) {
				fetchVillagesOfTown(t, db)
				// fmt.Println("Action mock:fetchVillagesOfTown", t)
				wg.Done()
			}(town)
		}
		wg.Wait()
	}

	fmt.Println("[Fetch all villages done]", time.Now().Sub(startT))

}

func demoEncoding() {
	utf8Str := "中文"

	utf8bytes := []byte(utf8Str)

	fmt.Println(len(utf8Str), len(utf8bytes))

}

func main() {
	fetchAllProvinces()
	fetchAllCities()
	fetchAllCounties()
	fetchAllTowns()

	// fetchAllVillages()
}
