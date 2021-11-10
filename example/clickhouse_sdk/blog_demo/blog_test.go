/**
 * @Author: vincent
 * @Description:
 * @File:  TestBlog
 * @Version: 1.0.0
 * @Date: 2021/11/9 14:54
 */

package blog_demo

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"

	"github.com/yanyiwu/gojieba"

	"github.com/ClickHouse/clickhouse-go"
)

var blogs = []string{
	`我的11.11必囤榜单NO.1，猜猜看是什么？作为“11.11京东美妆OLAY品牌选品官”，我首先挑中的是美白精华@OLAY 抗糖小白瓶！它有革命性抗糖成分革糖素搭配烟酰胺，三重抗糖减黄气，能让肌肤白到发光~ 想要趁秋冬变白减黄气的宝宝们，别错过#京选大牌 美丽价到#福利，到@京东美妆 搜OLAY就能抢我的同款抗糖小白瓶N玉兰油（OLAY）抗糖小白瓶精华液50ml面部精华... ，趁11.11优惠多多赶紧囤货吧！`,
	`“当时风花纷纷应是好时辰
  何故要端坐风月里孤身”
心无旁骛应该是很甜的啊！！
可是“太适合牺牲”？？感觉又是虐的？
海市和方诸的故事到底是怎么又甜又虐？
来看看@电视剧斛珠夫人 就知道啦！！！`,
	`过安检所有物品都按要求放好，却总是通不过，问题到底出在哪里？点击视频揭晓答案。 `,
	`有这样一群人，烈火中逆行，险境里救援，洪水中坚守，他们有一个共同的名字—“火焰蓝”。
没有生来英勇，只是选择无畏，119消防宣传日，一起#点亮火焰蓝#，致敬最美逆行者。 ​​​​
`,
	`【17岁少女肠道布满寄生虫！医生提醒：吃饭别犯这个错！】最近广东17岁的女孩小青经常感到腹痛腹泻，去医院一查竟发现她的肠道里布满寄生虫，还活蹦乱跳的。医生称：随便一个镜头下都看到有五六条虫，有的甚至咬住肠粘膜…原来这道家常菜成了她患病的元凶↓↓医生提醒，这些菜不能这样吃`,
	`湖人这球赢的真不容易，瓜哥今天的发挥太重要了，全场29分进了7个三分，加时赛这个三分也非常重要，这赛季要赢球居然真的得靠瓜哥的发挥…… ​​​​`,
	`改版无数次的美沫油膏上架咯！！O抽奖详情
❄️油敷法太适合这个季节了
这款油膏它能真正的做到油肤
但不会弄到满脸是油 糊也根本不存在❌
✔️疏通毛孔 平衡油脂分泌
✔️提亮去黄超级明显 🧏🏻‍♀️小脸都嫩爆了
如果早上起来满脸油光的媳妇更加要试试🙋🏻‍♀️第二天肉眼可见 你一定会来感谢我🙏
一瓶多用可当面膜、面霜、按摩膏的性价比之王`,
	`我在想一个大事
有点不计后果
但是还是要征求同事们的同意
你们等我一下`,
	`XX片女星珍妮李曾叱咤美国XX市场，以亮丽外型、姣好身材闯出一片天，还因此封为“XX女帝”，却在2009年突然宣布放弃女优事业，转往模特儿界发展。不过，她自此渐渐销声匿迹，曾住在豪华公寓的她，近年被发现栖身于下水道中。

珍妮李出身美国田纳西州，本名史蒂芬妮‧萨多拉，在XX网站上有超过破亿流量，更有4万5000人订阅。作品无数的她，当年有着“XX女帝”的地位，更有不少杂志邀约拍摄，珍妮李住在豪华的顶层公寓，过着时髦又奢华的生活。 然而，珍妮李12年前突然宣布从XX界引退，从此不再拍片，并转型当起模特儿，之后却逐渐消声匿迹，消失在荧光幕前。直到2019年，外媒找到了她，当年的“XX女帝”，如今却是住在美国拉斯维加斯的一座下水道。

据报道，珍妮李的“房间”中，仅有一张床、一个柜子、一盏灯，还有简单几件衣物，和过去居住的豪宅呈现强烈对比。她当时虽未透露住在地下水道的原因，但表示自己过得很快乐，住在地下水道的日子，也没有外界想像得那么糟，和邻居们也都相处愉快：“大家也都彼此尊重，人都非常好，我交到更多真诚的朋友。” `,
	`#美国音乐节踩踏事件歌手及组织方被起诉# 综合多家美国媒体报道，截至7日，至少两名休斯敦“天文世界”音乐节踩踏事件受伤者提起诉讼，指控音乐节主创者、美国当红说唱歌手特拉维斯·斯科特等人及音乐节组织方轻忽观众安全，以致酿成至少8人死亡的踩踏惨剧。#美国一音乐节踩踏致人死亡演出仍继续# （cr.环球资讯）`,
}

func TestCreateTable(t *testing.T) {
	// tcp://127.0.0.1:9000?debug=true，里面的debug=true，会打印如下的调试信息
	//[clickhouse]host(s)=127.0.0.1:9000, database=default, username=default
	//[clickhouse][dial] secure=false, skip_verify=false, strategy=random, ident=1, server=0 -> 127.0.0.1:9000
	//[clickhouse][connect=1][hello] -> Golang SQLDriver 1.1.54213
	//[clickhouse][connect=1][hello] <- ClickHouse 21.10.54449 (UTC)
	//[clickhouse][connect=1]-> ping
	//[clickhouse][connect=1][process] <- pong

	// step 1: connect
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}

	// step 2: create table
	// todo: 需要加自增id
	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS blog(
		    uid UInt64,
		    uname String,
		    mid UInt64,
		    content String,
		    words Array(String),
		    create_day Date,
		    create_time DateTime
		) engine=Memory
	`)

	if err != nil {
		log.Fatal(err)
	}

	use_hmm := true
	x := gojieba.NewJieba()
	defer x.Free()

	// step 3: insert with transaction

	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO blog (uid, uname, mid, content, words, create_day, create_time) VALUES (?, ?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()
	iStart := 10000
	cnt := len(blogs)
	iEnd := iStart + cnt
	for i := iStart; i < iEnd; i++ {
		randDate := randDateUpToNow()
		uid := i
		uname := "mengxin_" + strconv.Itoa(i)
		mid := rand.Uint64()
		content := blogs[i%cnt]
		words := x.Cut(content, use_hmm)
		if _, err := stmt.Exec(
			uid,
			uname,
			mid,
			content,
			clickhouse.Array(words),
			//clickhouse.Array([]int16{1, 2, 3}),
			randDate,
			randDate,
		); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func TestBlogCount(t *testing.T) {
	fmt.Println(len(blogs))
}
