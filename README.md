## ERP同步工具桌面端

#### 编译图标资源

```shell
go:generate rsrc -ico assets/images/favicon.ico -manifest assets/app.manifest -o main.syso
```

#### 编译程序

```shell
go build "-ldflags=-H windowsgui" -o app.exe app
```