# zrpc
### ZZINNO自用rpc
一个轻量级的rpc方案，基于net/rpc封装。



# 设计思路
1. 所有服务端平级，通过点对点通信完成信息的传输
2. 每个端必定带有自己的名字来界定信息的给出和返回
4. 每个端可以处理的函数自行注册和管理，比如消费者注册个say函数，生产者给hello，消费者返回world，如果没得，就必须明确拒绝
5. 生产者可以保存多个配置连接，并根据需要进行调用

# 使用注意
推荐使用github.com/vmihailenco/msgpack
对传输前后的[]byte进行处理


# 性能
因为是轻量级的，中小型企业用用就很方便,但是百万并发大厂出门左转https://rpcx.io

```
测试平台：
                 ..                    *************** 
               .PLTJ.                  --------------- 
              <><><><>                 OS: CentOS Linux 7 (Core) x86_64 
     KKSSV' 4KKK LJ KKKL.'VSSKK        Host: 2288H V5 Purley 
     KKV' 4KKKKK LJ KKKKAL 'VKK        Kernel: 3.10.0-1062.18.1.el7.x86_64 
     V' ' 'VKKKK LJ KKKKV' ' 'V        Uptime: *************************** 
     .4MA.' 'VKK LJ KKV' '.4Mb.        Packages: ********** 
   . KKKKKA.' 'V LJ V' '.4KKKKK .      Shell: zsh 5.0.2 
 .4D KKKKKKKA.'' LJ ''.4KKKKKKK FA.    Theme: Adwaita [GTK2/3] 
<QDD ++++++++++++  ++++++++++++ GFD>   Icons: Adwaita [GTK2/3] 
 'VD KKKKKKKK'.. LJ ..'KKKKKKKK FV     Terminal: /dev/pts/1 
   ' VKKKKK'. .4 LJ K. .'KKKKKV '      CPU: Intel Xeon Silver 4210 (40) @ 3.200GHz 
      'VK'. .4KK LJ KKA. .'KV'         GPU: Intelligent Management system chip w/VGA support] 
     A. . .4KKKK LJ KKKKA. . .4        Memory: ******MiB / 63845MiB 
     KKA. 'KKKKK LJ KKKKK' .4KK
     KKSSA. VKKK LJ KKKV .4SSKK                                
              <><><><>                                         
               'MKKM'
                 ''

====================================
2020/12/27 12:52:35 2020-12-27 12:52:35.6447072 +0800 CST m=+0.015991001 [INFO] [[Function Index From  hello : [ s]]]
request count: 500000
success rate: 100 %
max cost: 9397 ms
avg cost: 3718 ms
all cost: 9937 ms
ERRNUM: 0
rps: 50316 request/s
2020/12/27 12:52:45 2020-12-27 12:52:45.6015023 +0800 CST m=+9.972786101 [INFO] [[Function Index From  hello : [ s]]]
request count: 500000
success rate: 100 %
max cost: 7816 ms
avg cost: 3999 ms
all cost: 7965 ms
ERRNUM: 0
rps: 62774 request/s
2020/12/27 12:52:53 2020-12-27 12:52:53.5909341 +0800 CST m=+17.962217901 [INFO] [[Function Index From  hello : [ s]]]
request count: 500000
success rate: 100 %
max cost: 6364 ms
avg cost: 3268 ms
all cost: 6683 ms
ERRNUM: 0
rps: 74816 request/s
2020/12/27 12:53:00 2020-12-27 12:53:00.2947751 +0800 CST m=+24.666058901 [INFO] [[Function Index From  hello : [ s]]]
request count: 500000
success rate: 100 %
max cost: 7219 ms
avg cost: 3974 ms
all cost: 7791 ms
ERRNUM: 0
rps: 64176 request/s
```
