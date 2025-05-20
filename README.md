## ERP同步工具桌面端

#### 编译图标资源

```shell
go:generate rsrc -ico assets/images/favicon.ico -manifest assets/app.manifest -o main.syso
```

#### 编译程序

```shell
go build "-ldflags=-H windowsgui" -o app.exe app
```

#### V3版本说明

```text
    - 支持断点继续任务
    - UI重构优化，更换UI框架（解决卡顿、布局混乱、不好拓展的问题、可跨平台）
    - 提升比对差异效率（json->结构体)
    - 使用本地序列化(持久化)和内存缓存替代redis，提升效率，减少内存和网络负载
    - 项目结构重构，并优化代码可读性
    - UI和业务逻辑分离，可直接编译无UI版本运行到linux
    - 使用odbc+gorm替代sqlx，更加统一
    - 无需安装和管理员权限，单文件exe直接可运行
    - 使用异步并发批处理任务，同时执行多行数据，提升效率
    - 重构同步任务的运行机制，不同类型的无关任务不会被阻塞，之前版本客户资料同步全部完成前不会同步订单
```