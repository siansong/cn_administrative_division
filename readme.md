

## 使用

装postgresql

clone repo

```bash
go get -u "github.com/PuerkitoBio/goquery"
go get -u "github.com/djimenez/iconv-go"
go get -u "github.com/go-pg/pg"
go get -u "github.com/go-pg/pg/orm"
```

在main.go中搜Password,改成自己的数据库账号密码，表会自动建的

`go run main.go`

这样可以抓到街道/乡镇一级，如果要抓社区/村级别，那稍微看下代码吧，数据量有点大所以默认就不抓了


## Performance

基于我配置：
```bash
MacBook Pro 2017
2.8 GHz Intel Core i7
$sysctl -a | grep ".cpu."
hw.ncpu: 8
hw.physicalcpu: 4
hw.physicalcpu_max: 4
hw.logicalcpu: 8

16 GB 2133 MHz LPDDR3
```
到乡镇级别,约`47107`rows,大概60s左右(看网络状况吧，我快的时候跑到50s+)


TODO
----

1. 数据基于国家统计局网站抓取，较民政部的即时性不足，有改进的方案
2. <del>[dep]支持其他数据库</del>，划掉了，喜欢的话装个postgresql，或者fork了自己玩，完全没有难度的吧「应该」
3. 增加测试
4. self parentCode ref FK
5. db复用，不用每次都Close()
6. <del>[done]增加retry机制，某些页面偶尔的会！=200</del>
7. 寻求容错机制，即便众多过程中有失败，也能保证过最终数据的完整性
8. connection reset by peer 😂


## License

This repo is released under the [WTFPL](http://www.wtfpl.net/) – Do What the Fuck You Want to Public License.
