# QuantBot

最近研究量化交易，学习了很好的一个项目：[Samaritan](https://github.com/miaolz123/samaritan)

可惜这个项目已经很久很久没有更新過了，另外项目中也有一些BUG，其中最致命的BUG就是在实现javascript并发任务这个功能的时候没有考虑资源冲突的处理，导致程序无法正常工作，所以对这些部分进行了一些修改，使其可以正常工作，并且更新了文档，另外原有的一些交易所接口也因为各种原因失效了，所以这里也重新更新了部分交易所接口，比如火币，比特儿国际，币安，OKEX等，并且更新了文档，然后给项目重新改了个更直观的名字。另外每个交易所的交易对只是选取了几个大币种的，如果需要添加新的交易对，可以修改源代码进行添加。还有就是某些交易所需要搭梯子，测试的时候请自行准备梯子，并修改对应交易所接口源码。

这里我写了个简单的搬砖演示程序：[代码](https://github.com/phonegapX/trader-sample) [博客](http://phonegap.me/post/52.html)

[更新的文档](http://www.quantbot.org/#/)

版本更新到v0.0.2，新增了中币交易所的接口。  
版本更新到v0.0.3，新增了BigONE交易所的接口，当前热门的 “交易挖矿+持币分红” 交易所。(2018-07-09)

## 编译

直接运行build.sh进行编译，需要用到xgo和Docker，可以同时编译出多个平台下的可执行文件，如果不需要多平台支持，可以直接用go编译，推荐使用LiteIDE。

## 包依赖问题

通过使用[glide](https://github.com/Masterminds/glide)工具可以解决包的依赖。安装glide后执行下列命令。

```shell
$ cd QuantBot
$ glide install
```

## 支持的交易所

| 交易所 | 货币类型 |
| -------- | ----- |
| zb | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `LTC/USDT`, `QTUM/USDT` |
| okex | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `ONT/USDT`, `QTUM/USDT`, `ONT/ETH` |
| 火币网 | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `ONT/USDT`, `QTUM/USDT` |
| 比特儿国际 | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `ONT/USDT`, `QTUM/USDT` |
| 币安 | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `ONT/USDT`, `QTUM/USDT` |
| poloniex | `ETH/BTC`, `XMR/BTC`, `BTC/USDT`, `LTC/BTC`, `ETC/BTC`, `XRP/BTC`, `ETH/USDT`, `ETC/ETH`, ... |
| okex 期货 | `BTC.WEEK/USD`, `BTC.WEEK2/USD`, `BTC.MONTH3/USD`, `LTC.WEEK/USD`, ... |
| BigONE | `BTC/USDT`, `ONE/USDT`, `EOS/USDT`, `ETH/USDT`, `BCH/USDT`, `EOS/ETH` |
