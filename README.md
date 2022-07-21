# sshman

Sshman, is a more real world example of how crud can help you build a
CRUD REST API fast and easily.

Literally, sshman means "ssh manager", which indicates that it is a system
for managing ssh resources. To be specific, the question is defined as:

> We are an organization that wants to share our computing resources with
> the community by offering people ssh accesses.
> And we are developing a system that will help us manage
> our user accounts, available hosts and ssh keys.
>
> Even such a simple system, we want to make it easy to use, say,
> for muggles. So we need to make it a well-designed beautiful
> web application, of course, with the fancy technologies like
> separation of front-end and backend, JWT authentication, etc.
> (It's 2022 after all).
>
> Here, sshman is the backend CRUD API service.

OK, talk is cheap, let's get started. Here is the layout of the project:

```
sshman
├── README.md       the verbose documentation you are reading now
├── auth.go         JWT authentication and User Role authentication
├── config.go       configuration for the application
├── config.yaml     a example configuration file
├── controller.go   controllers for users' APIs
├── db.go           connect to the database, and initialize it
├── go.mod/go.sum   your friendly neighborhood go module
├── log.go          initialize a global logger
├── main.go         the program's entry point
├── models.go       models: Host, Session and User
├── router.go       the router with custom middlewares
└── test.http       example HTTP requests for APIs
```

Before you turn to the code, let's take a look at some more details about the
design of the project.

Firstly, All you need is models!
Very straightforward, the entries we need to manage are:

- **Hosts**: ssh hosts (listening ssh service at ip:port)
- **Sessions**: ssh sessions (user:password@host)
- **Users**: ssh users, a user is allowed to have many sessions, while a session can be used by multiple users

Host:

| Field    | Type        | Description             |
|----------|-------------|-------------------------|
| id       | `integer`   | The user id             |
| hostname | `string`    | The hostname            |
| ip       | `string`    | The ip address          |
| port     | `integer`   | The port to SSH service |

Session:

| Field      | Type       | Description                         |
|------------|------------|-------------------------------------|
| id         | `integer`  | The user id                         |
| host       | `Host`     | The host to connect                 |
| username   | `string`   | The username to connect to host     |
| privateKey | `string`   | The private key to connect to host  |

User:

| Field     | Type        | Description                                                                                    |
|-----------|-------------|------------------------------------------------------------------------------------------------|
| id        | `integer`   | The user id                                                                                    |
| name      | `string`    | The user's name                                                                                |
| email     | `string`    | The user's email                                                                               |
| role      | `enum`      | The user's role: user or admin, user can only query his/her own sessions, admin access any API |
| sessions  | `[]Session` | available SSH sessions for user                                                                |

We need to create models for each of these entities (models.go).
And use crud to create the CRUD APIs for each of them (router.go):

```http
GET /hosts                # list all hosts
GET /hosts/:id            # get a host by id

GET /users                # list all users
GET /users/:id            # get a user by id

Get /users/:id/sessions   # list all sessions for a user

POST /hosts               # create a host
POST /users               # create a user
POST /users/:id/sessions  # create or add a session for a user

PUT /hosts/:id            # update a host
PUT /users/:id            # update a user
PUT /sessions/:id         # update a session

DELETE /hosts/:id         # delete a host
DELETE /users/:id         # delete a user
DELETE /sessions/:id      # delete a session

DELETE /users/:id/sessions/:id  # remove a session from user
```

But, the API designs above really sucks.

Since we do not want everyone to be able to access all the APIs, we need to
add some authentication. To make it simple, we split the router into two parts:

- `/admin/*`: only admin can access these APIs, full CRUDs for all entities;
- `/user/*` : for user, only read-only APIs for his/her own User and Sessions (with the nested Host);

We will use JWT authentication to protect all the APIs.
And we will add a middleware to check the user's role, preventing user from
accessing the admin APIs (auth.go).

And happily, some of your handsome colleagues have already implemented a user login
and JWT authentication module (which is another service). So we decide to
use their login system (i.e. we will not write a login interface) and simply
steal their JWT secret to verify tokens from requests (auth.go).

BTW, to show you how to use other low level parts of crud package, we add a
`controller.go` to hand-craft APIs for users (`/user/*`). It is unnecessary to
do so, (a simple middleware makes it), but it is a good example of how to use
crud in your real world applications.

 
