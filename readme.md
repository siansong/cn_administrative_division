

## 使用

装postgresql

在main.go中搜Password,改成自己的数据库账号密码

`go run main.go`

QA
---

Q:Why this project?<br>
A:之前找到一个，nodejs写的(https://github.com/modood/Administrative-divisions-of-China)，不过有点不满意，原本想用ts改进的，想想干脆用golang造个轮子

Q:https://github.com/modood/Administrative-divisions-of-China 有什么不满<br>
A:主要是速度有点慢，另外，觉得数据结构也可以改进下

Q:


TODO
----

1. 数据基于国家统计局网站抓取，较民政部的即时性不足，有改进的方案
2. <del>[dep]支持其他数据库</del>，划掉了，喜欢的话装个postgresql，或者fork了自己玩，完全没有难度的吧「应该」
3. 增加测试
4. self parentCode ref FK
5. db复用，不用每次都Close()
6. 增加retry机制，某些页面偶尔的会！=200
7. 寻求容错机制，即便众多过程中有失败，也能保证过最终数据的完整性
8. connection reset by peer 😂


其他
---

china_administrative_division >> cad