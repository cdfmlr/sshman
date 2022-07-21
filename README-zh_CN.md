# sshman

[English](README.md) | 机翻中文

Sshman，是一个更真实的例子，说明 crud 如何帮助你快速轻松地建立一个CRUD REST API。

从字面上看，sshman的意思是 "SSH 管理器"，这表明它是一个管理 SSH 资源的系统。具体来说，这个问题的定义是：

> 我们是一个组织，希望通过向人们提供ssh访问权来与社区分享我们的计算资源。我们正在开发一个系统，以帮助我们管理我们的用户账户、可用的主机和ssh密钥。
>
> 即使是这样一个简单的系统，我们也想让它易于使用，比如说，对于麻瓜来说。所以我们需要把它做成一个设计良好的漂亮的网络应用程序，当然，要有花哨的技术，比如前端和后端分离，JWT 认证等等（毕竟是2022年）。
>
> 这里，sshman 就是后端的 CRUD API 服务。

好了，话不多说，让我们开始吧。这里是项目的布局：

```
sshman
├── README.md       现在正在阅读的冗长的文档
├── auth.go         JWT 认证和用户角色认证
├── config.go       应用程序的配置
├── config.yaml     是一个配置文件的例子
├── controller.go   用户 API 的控制器
├── db.go           连接到数据库，并初始化它
├── go.mod/go.sum   你的友好邻居 go 模块
├── log.go          初始化全局记录器
├── main.go         是程序的入口
├── models.go       模型：主机(Host)、会话(Session)和用户(User)
├── router.go       使用自定义的中间件的路由器。
└── test.http       对 API 的 HTTP 请求的例子
```

在你转向代码之前，让我们再看看关于项目设计的一些细节。

首先，你所需要的就是模型! 非常直接，我们需要管理的条目是:

- Hosts: SSH 主机 (监听 SSH 服务在 ip:port)
- Sessions: SSH sessions (user:password@host)
- User：SSH 用户，一个用户可以有很多会话，而一个会话可以被多个用户使用。

Host:

| 字段       | 类型       | 描述        |
|----------|----------|-----------|
| id       | integer  | 主机的 ID    |
| hostname | string   | 主机名       |
| ip       | 字符串      | ip 地址     |
| port     | 整数       | SSH 服务的端口 |

Session:

| 字段         | 类型     | 描述         |
|------------|--------|------------|
| id         | 整数     | session ID |
| host       | Host   | 连接的主机      |
| username   | string | 连接到主机的用户名  |
| privateKey | 字符串    | 连接到主机的私钥   |

User:

| 字段        | 类型        | 描述                                       |
|-----------|-----------|------------------------------------------|
| id        | integer   | 用户的ID                                    |
| name      | string    | 用户的名字                                    |
| 电子邮件      |  字符串      | 用户的电子邮件                                  |
| role      | enum      | 用户的角色：用户或管理员，用户只能查询他/她自己的会话，管理员可以访问任何API |
| sessions  | []Session | 用户的可用SSH会话                               |

我们需要为这些实体中的每一个创建模型（models.go）。并使用 crud 为每个实体创建 CRUD APIs（router.go）。

```
GET /hosts               # 列出所有主机
GET /hosts/:id           # 通过id获得一个主机

GET /users               # 列出所有用户
GET /users/:id           # 通过id获得一个用户

GET /users/:id/sessions  # 列出一个用户的所有会话

POST /hosts              # 创建一个主机
POST /users              # 创建一个用户
POST /users/:id/sessions # 为一个用户创建或添加会话

PUT /hosts/:id           # 更新一个主机
PUT /users/:id           # 更新一个用户
PUT /sessions/:id        # 更新一个会话

DELETE /hosts/:id        # 删除一个主机
DELETE /users/:id        # 删除一个用户
DELETE /sessions/:id     # 删除一个会话

DELETE /users/:id/sessions/:id   # 从用户身上删除一个会话
```

但是，上面的 API 设计真的很糟糕。
由于我们不希望每个人都能访问所有的 API，我们需要添加一些认证。为了简单起见，我们把路由器分成两部分。

- `/admin/*`：只有管理员可以访问这些 API，对所有实体进行完整的 CRUD。
- `/user/*` ：对于用户来说，只有他/她自己的用户和会话（及其嵌套的主机）的只读 API。

我们将使用 JWT 认证来保护所有的 API。我们将添加一个中间件来检查用户的角色，防止用户访问管理员的 API（auth.go）。

而令人高兴的是，你的一些英俊的同事已经实现了一个用户登录和 JWT 认证模块（这是另一个服务）。所以我们决定使用他们的登录系统（也就是说，我们不会写一个登录接口），并简单地窃取他们的 JWT 密钥来验证请求中的 token（auth.go）。

BTW，为了向你展示如何使用 crud 包的其他低级部分，我们添加了一个 controller.go 来为一般用户手工制作 API（`/user/*`）。
这样做是不必要的，（一个简单的中间件就可以做到），但这是一个很好的例子，说明如何在你的实际应用中使用 crud。

> 该文档通过www.DeepL.com/Translator（免费版）翻译
