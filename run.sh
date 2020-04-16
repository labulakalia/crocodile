nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
nohup ./bin/crocodile client -c core.toml &
./bin/crocodile server -c core.toml



2020-04-17T01:40:00.065+0800    error   schedule/schedule2.go:687       task run failed {"taskid": "256572987300384768", "error": "errgroup: panic recovered: runtime error: invalid memory address or nil pointer dereference\ngoroutine 265 [running]:\ngithub.com/labulaka521/crocodile/common/errgroup.(*Group).do.func1(0xc001b93f80, 0xc001bf3340)\n\t/Users/labulakalia/workerspace/golang/crocodile/common/errgroup/errgroup.go:55 +0x146\npanic(0x4b49960, 0x66a2470)\n\t/usr/local/Cellar/go/1.13.1/libexec/src/runtime/panic.go:679 +0x1b2\ngoogle.golang.org/grpc.(*ClientConn).Target(...)\n\t/Users/labulakalia/go/pkg/mod/google.golang.org/grpc@v1.25.1/clientconn.go:765\ngithub.com/labulaka521/crocodile/core/schedule.(*task2).runTask(0xc0002d8a90, 0x4e3d460, 0xc001c09000, 0xc001b9cae0, 0x12, 0x4030
201, 0x0, 0x0)\n\t/Users/labulakalia/workerspace/golang/crocodile/core/schedule/schedule2.go:791 +0x203e\ngithub.com/labulaka521/crocodile/core/schedule.(*task2).StartRun.func5(0x4e3d460, 0xc001c09000, 0xc0000ca9c0, 0x0)\n\t/Users/labulakalia/workerspace/golang/crocodile/core/schedule/schedule2.go:679 +0x5f\ngithub.com/labulaka521/crocodile/common/errgroup.(*Group).do(0xc001bf3340, 0xc001c1e680)\n\t/Users/labulakalia/workerspace/golang/crocodile/common/errgroup/errgroup.go:68 +0xaf\ngithub.com/labulaka521/crocodile/common/errgroup.(*Group).GOMAXPROCS.func1.1(0xc001bf3340)\n\t/Users/labulakalia/workerspace/golang/crocodile/common/errgroup/errgroup.go:85 +0x57\ncreated by github.com/labulaka521/crocodile/common/errgroup.(*Group).GOMAXPROCS.func1\n\t/Users/labulakalia/workerspace/golang/crocodile/common/errgroup/errgroup.go:83 +0x87\n"}
2020-04-17T01:40:00.079+0800    info    model/log.go:21 start savelog   {"tasklog": {"name":"testtask","runby_taskid":"256572987300384768","start_time":1587058800011,"start_timestr":"","end_time":1587058800072,"end_timestr":"","total_runtime":61,"sta