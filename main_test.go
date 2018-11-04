package main

import (
	"fmt"
	"testing"

	// "github.com/go-pg/pg"            // postgresql orm
	"github.com/go-pg/pg/orm" // 同上
)

var testDataProvinceNames = []string{"北京市", "天津市", "河北省", "山西省", "内蒙古自治区", "辽宁省", "吉林省", "黑龙江省", "上海市", "江苏省", "浙江省", "安徽省", "福建省", "江西省", "山东省", "河南省", "湖北省", "湖南省", "广东省", "广西壮族自治区", "海南省", "重庆市", "四川省", "贵州省", "云南省", "西藏自治区", "陕西省", "甘肃省", "青海省", "宁夏回族自治区", "新疆维吾尔自治区"}
var testDataProvinceNamesMap = map[string]int{
	"北京市": 233, "天津市": 233, "河北省": 233, "山西省": 233, "内蒙古自治区": 233, "辽宁省": 233, "吉林省": 233, "黑龙江省": 233, "上海市": 233, "江苏省": 233, "浙江省": 233, "安徽省": 233, "福建省": 233, "江西省": 233, "山东省": 233, "河南省": 233, "湖北省": 233, "湖南省": 233, "广东省": 233, "广西壮族自治区": 233, "海南省": 233, "重庆市": 233, "四川省": 233, "贵州省": 233, "云南省": 233, "西藏自治区": 233, "陕西省": 233, "甘肃省": 233, "青海省": 233, "宁夏回族自治区": 233, "新疆维吾尔自治区": 233,
}

var testDataHeiBeiSheng = Area{Code: "13", Name: "河北省"}

// 河北省辖属市
var testDataHeBeiCities = []Area{
	Area{Code: "130100000000", Name: "石家庄市"},
	Area{Code: "130200000000", Name: "唐山市"},
	Area{Code: "130300000000", Name: "秦皇岛市"},
	Area{Code: "130400000000", Name: "邯郸市"},
	Area{Code: "130500000000", Name: "邢台市"},
	Area{Code: "130600000000", Name: "保定市"},
	Area{Code: "130700000000", Name: "张家口市"},
	Area{Code: "130800000000", Name: "承德市"},
	Area{Code: "130900000000", Name: "沧州市"},
	Area{Code: "131000000000", Name: "廊坊市"},
	Area{Code: "131100000000", Name: "衡水市"},
}

// TangShanCounties 唐山辖属区县，http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2017/13/1302.html
var testDataTangShanCounties = []Area{
	Area{Code: "130201000000", Name: "市辖区"},
	Area{Code: "130202000000", Name: "路南区"},
	Area{Code: "130203000000", Name: "路北区"},
	Area{Code: "130204000000", Name: "古冶区"},
	Area{Code: "130205000000", Name: "开平区"},
	Area{Code: "130207000000", Name: "丰南区"},
	Area{Code: "130208000000", Name: "丰润区"},
	Area{Code: "130209000000", Name: "曹妃甸区"},
	Area{Code: "130223000000", Name: "滦县"},
	Area{Code: "130224000000", Name: "滦南县"},
	Area{Code: "130225000000", Name: "乐亭县"},
	Area{Code: "130227000000", Name: "迁西县"},
	Area{Code: "130229000000", Name: "玉田县"},
	Area{Code: "130271000000", Name: "唐山市芦台经济技术开发区"},
	Area{Code: "130272000000", Name: "唐山市汉沽管理区"},
	Area{Code: "130273000000", Name: "唐山高新技术产业开发区"},
	Area{Code: "130274000000", Name: "河北唐山海港经济开发区"},
	Area{Code: "130281000000", Name: "遵化市"},
	Area{Code: "130283000000", Name: "迁安市"},
}

var testDataBeiJingCityCode = "110100000000"   // 北京市辖区
var testDataHaiDianCountyCode = "110108000000" // 海淀Code
var testDataHaiDianTowns = []Area{
	Area{Code: "110108001000", Name: "万寿路街道办事处"},
	Area{Code: "110108002000", Name: "永定路街道办事处"},
	Area{Code: "110108003000", Name: "羊坊店街道办事处"},
	Area{Code: "110108004000", Name: "甘家口街道办事处"},
	Area{Code: "110108005000", Name: "八里庄街道办事处"},
	Area{Code: "110108006000", Name: "紫竹院街道办事处"},
	Area{Code: "110108007000", Name: "北下关街道办事处"},
	Area{Code: "110108008000", Name: "北太平庄街道办事处"},
	Area{Code: "110108010000", Name: "学院路街道办事处"},
	Area{Code: "110108011000", Name: "中关村街道办事处"},
	Area{Code: "110108012000", Name: "海淀街道办事处"},
	Area{Code: "110108013000", Name: "青龙桥街道办事处"},
	Area{Code: "110108014000", Name: "清华园街道办事处"},
	Area{Code: "110108015000", Name: "燕园街道办事处"},
	Area{Code: "110108016000", Name: "香山街道办事处"},
	Area{Code: "110108017000", Name: "清河街道办事处"},
	Area{Code: "110108018000", Name: "花园路街道办事处"},
	Area{Code: "110108019000", Name: "西三旗街道办事处"},
	Area{Code: "110108020000", Name: "马连洼街道办事处"},
	Area{Code: "110108021000", Name: "田村路街道办事处"},
	Area{Code: "110108022000", Name: "上地街道办事处"},
	Area{Code: "110108023000", Name: "万柳地区办事处"},
	Area{Code: "110108024000", Name: "东升地区办事处"},
	Area{Code: "110108025000", Name: "曙光街道办事处"},
	Area{Code: "110108026000", Name: "温泉地区办事处"},
	Area{Code: "110108027000", Name: "四季青地区办事处"},
	Area{Code: "110108028000", Name: "西北旺地区办事处"},
	Area{Code: "110108029000", Name: "苏家坨地区办事处"},
	Area{Code: "110108030000", Name: "上庄地区办事处"},
}

// 110108019000	西三旗街道办事处
var testDataXiSanQiTownCode = "110108019000"
var testDataXiSanQiTownVillages = []Area{
	Area{Code: "110108019003", Name: "永泰园第一社区居委会"},
	Area{Code: "110108019004", Name: "清缘里社区居委会"},
	Area{Code: "110108019005", Name: "建材西里社区居委会"},
	Area{Code: "110108019006", Name: "机械学院联合社区居委会"},
	Area{Code: "110108019009", Name: "永泰西里社区居委会"},
	Area{Code: "110108019010", Name: "宝盛里社区居委会"},
	Area{Code: "110108019012", Name: "建材东里社区居委会"},
	Area{Code: "110108019013", Name: "清缘东里社区居委会"},
	Area{Code: "110108019014", Name: "悦秀园社区居委会"},
	Area{Code: "110108019015", Name: "电科院社区居委会"},
	Area{Code: "110108019016", Name: "沁春家园社区居委会"},
	Area{Code: "110108019018", Name: "育新花园社区居委会"},
	Area{Code: "110108019021", Name: "冶金研究院社区居委会"},
	Area{Code: "110108019022", Name: "北新集团社区居委会"},
	Area{Code: "110108019023", Name: "9511工厂联合社区居委会"},
	Area{Code: "110108019024", Name: "永泰庄社区居委会"},
	Area{Code: "110108019025", Name: "清润家园社区居委会"},
	Area{Code: "110108019026", Name: "永泰园第二社区居委会"},
	Area{Code: "110108019027", Name: "建材城联合社区居委会"},
	Area{Code: "110108019028", Name: "小营联合社区居委会"},
	Area{Code: "110108019029", Name: "怡清园社区居委会"},
	Area{Code: "110108019030", Name: "枫丹丽舍社区居委会"},
	Area{Code: "110108019032", Name: "知本时代社区居委会"},
	Area{Code: "110108019033", Name: "清景园社区居委会"},
	Area{Code: "110108019034", Name: "清缘西里社区居委会"},
	Area{Code: "110108019035", Name: "富力桃园社区居委会"},
	Area{Code: "110108019036", Name: "永泰东里社区居委会"},
}

func dropDb() {
	db := connTryCreateTb()
	defer db.Close()

	db.DropTable(&Area{}, &orm.DropTableOptions{
		IfExists: true,
	})
}

func TestParseRelativeHref2BaseURL(t *testing.T) {
	refferURL := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2017/11/1101.html"
	href := "01/110108.html"
	relativeHref := parseRelativeHref2BaseURL(href, refferURL)

	shouldBeURL := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2017/11/01/110108.html"
	if BaseURL+relativeHref != shouldBeURL {
		t.Errorf("没解析对\nExp: %s , \nGot: %s", shouldBeURL, BaseURL+relativeHref)
	}

}

func TestFetchProvinces(t *testing.T) {
	// t.Skip("跳过抓取省份测试")
	dropDb()

	fetchAllProvinces()

	provinces := findByLevelAndParent(LVProvince, ignoreCode)

	count := 0
	for _, prvc := range provinces {
		_, exists := testDataProvinceNamesMap[prvc.Name]
		if exists {
			count++
		}
	}

	if count != len(testDataProvinceNames) {
		t.Errorf("抓取到的省份量不对, Exp: [%d], Got: [%d]", len(testDataProvinceNames), count)
		t.Log("Db provinces:", provinces)
	}
}

func TestFetchCities(t *testing.T) {
	dropDb()
	fetchAllProvinces()
	fetchAllCities()

	db := connTryCreateTb()
	defer db.Close()

	cityCount, _ := db.Model(&Area{}).
		Where("level = ?", LVCity).
		Count()

	if cityCount < 300 {
		t.Errorf("至少有300多个市吧，Got: [%d]", cityCount)
	}

	hebeiCities := findByLevelAndParent(LVCity, testDataHeiBeiSheng.Code)
	if len(testDataHeBeiCities) != len(hebeiCities) {
		t.Errorf("抓取到的河北省的数据似乎不对:\n Exp: %v \n Got: %v", testDataHeBeiCities, hebeiCities)
	}

}

func TestFetchCountyOfCity(t *testing.T) {
	dropDb()

	fetchAllProvinces()
	fetchAllCities()

	db := connTryCreateTb()
	defer db.Close()
	cities := findByLevelAndParent(LVCity, testDataHeiBeiSheng.Code)

	var tangShanCity *Area
	for _, city := range cities {
		if city.Name == "唐山市" {
			tangShanCity = &city
			break
		}
	}
	if tangShanCity == nil {
		t.Error("数据库里找不到唐山")
	}
	fetchCountiesOfCity(*tangShanCity, db)

	tangShanCounties := findByLevelAndParent(LVCounty, tangShanCity.Code)

	if len(tangShanCounties) != len(testDataTangShanCounties) {
		t.Errorf("抓取到的唐山辖属区县数量好像不对, EXP: [%d], Got: [%d]", len(testDataTangShanCounties), len(tangShanCounties))
	}
}

func TestFetchAllCounties(t *testing.T) {
	dropDb()

	fetchAllProvinces()
	fetchAllCities()

	fetchAllCounties()

	db := connTryCreateTb()
	defer db.Close()

	count, _ := db.Model(&Area{}).Count()

	if count < 3000 {
		t.Errorf("感觉全国怎么也得有个3000多个区县的吧, Exp: > 3000, Got: [%d]", count)
	}
}

func TestFetchTownsOfCountyFetchVillagesOfTown(t *testing.T) {
	dropDb()

	fetchAllProvinces()
	fetchAllCities()

	db := connTryCreateTb()
	defer db.Close()

	beiJinCity := &Area{Code: testDataBeiJingCityCode}
	err := db.Select(beiJinCity)
	panicIf(err)
	fetchCountiesOfCity(*beiJinCity, db)

	haiDianCounty := &Area{Code: testDataHaiDianCountyCode}
	err = db.Select(haiDianCounty)
	panicIf(err)

	fetchTownsOfCounty(*haiDianCounty, db)

	townsOfHaidianCounty := findByLevelAndParent(LVTown, testDataHaiDianCountyCode)

	if len(testDataHaiDianTowns) != len(townsOfHaidianCounty) {
		t.Errorf("嗯，抓取到的海淀的towns似乎不对, Exp: [%d] Got: [%d] \n : %v", len(testDataHaiDianTowns), len(townsOfHaidianCounty), townsOfHaidianCounty)
	}

	xiSanQiTown := &Area{Code: testDataXiSanQiTownCode}
	err = db.Select(xiSanQiTown)
	panicIf(err)
	fetchVillagesOfTown(*xiSanQiTown, db)

	villagesOfXiSanQiTown := findByLevelAndParent(LVVillage, xiSanQiTown.Code)

	expectedLen, actualLen := len(testDataXiSanQiTownVillages), len(villagesOfXiSanQiTown)
	if expectedLen != actualLen {
		t.Errorf("似乎没抓到对的数据，Exp: [%d], Act: [%d]", expectedLen, actualLen)
	}

}

// TestFetchAllTowns FBI warning 这个测试用例啊，最好别跑，有点费时间
func TestFetchAllTowns(t *testing.T) {
	dropDb()
	fetchAllProvinces()
	fetchAllCities()
	fetchAllCounties()
	fetchAllTowns()

	db := connTryCreateTb()
	defer db.Close()

	count, _ := db.Model(&Area{}).
		Where("level = ?", LVTown).
		Count()

	if count < 40000 {
		t.Errorf("全国街道/乡镇差不多应该超过40000的吧？ Exp: [40000], Got: [%d]", count)
	}

	fmt.Println("db pool stats:", db.PoolStats())
}
