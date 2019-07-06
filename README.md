# mux-chi
  基于go-chi/chi框架定制而成，可快速构建http api(restful api)应用。

# golang版本
  建议用golang1.10+以上版本

# 日志插件
  基于 https://github.com/daheige/thinkgo/blob/master/common/log.go 实现每天流动式日志，可将日志记录到文件或输出到终端。

# 性能监控
    采用net/http/pprof包
        浏览器访问http://localhost:2338/debug/pprof，就可以查看
    在命令终端查看：
        安装graphviz
            $ apt install graphviz
        查看profile
            go tool pprof http://localhost:2338/debug/pprof/profile?seconds=60
            (pprof) top 10 --cum --sum
            (pprof) web  #web页面查看cpu使用情况

            每一列的含义：
            flat：给定函数上运行耗时
            flat%：同上的 CPU 运行耗时总比例
            sum%：给定函数累积使用 CPU 总比例
            cum：当前函数加上它之上的调用运行总耗时
            cum%：同上的 CPU 运行耗时总比例

        它会收集30s的性能profile,可以用go tool查看
            go tool pprof profile /home/heige/pprof/pprof.go-api.samples.cpu.002.pb.gz
        查看heap和goroutine
            查看活动对象的内存分配情况
            go tool pprof http://localhost:2338/debug/pprof/heap
            go tool pprof http://localhost:2338/debug/pprof/goroutine
        
        prometheus性能监控
        http://localhost:2338/metrics

# wrk工具压力测试
    https://github.com/wg/wrk
    
    ubuntu系统安装如下
    1、安装wrk
        # 安装 make 工具
        sudo apt-get install make git
        
        # 安装 gcc编译环境
        sudo apt-get install build-essential
        sudo mkdir /web/
        sudo chown -R $USER /web/
        cd /web/
        git clone https://github.com/wg/wrk.git
        # 开始编译
        cd /web/wrk
        make
    2、wrk压力测试
        $ wrk -c 100 -t 8 -d 2m http://localhost:1338/index
        Running 2m test @ http://localhost:1338/index
        8 threads and 100 connections
        Thread Stats   Avg      Stdev     Max   +/- Stdev
            Latency    17.45ms   31.76ms 633.02ms   95.52%
            Req/Sec     0.95k   180.65     1.64k    72.94%
        882466 requests in 2.00m, 148.96MB read
        Socket errors: connect 0, read 0, write 0, timeout 96
        Requests/sec:   7351.26
        Transfer/sec:      1.24MB

# 第三方包
  redisgo
  gorm
  go-chi/chi https://github.com/go-chi/chi
