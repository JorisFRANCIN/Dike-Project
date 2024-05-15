# Golang: gin-gonic Back-end

To launch the API from root of the project:

```
sudo systemctl start docker
sudo systemctl enable docker

docker compose up --build
```

To clean the MongoDB:

```
docker compose down
docker volume rm -f $(docker volume ls -q)
```

### This is the architecture:

```
// ├── about.json
// ├── controllers
// │   ├── about.go
// │   ├── service_controller.go
// │   └── user_controller.go
// ├── Dockerfile
// ├── Dockerfile.mongo
// ├── docs
// │   ├── docs.go
// │   ├── swagger.json
// │   └── swagger.yaml
// ├── go.code-workspace
// ├── go.mod
// ├── go.sum
// ├── helpers
// │   ├── authHelper.go
// │   └── tokenHelper.go
// ├── main.go
// ├── middleware
// │   ├── auth_middleware.go
// │   └── cors_middleware.go
// ├── models
// │   ├── cron_job.go
// │   ├── service.go
// │   └── user.go
// ├── README.md
// ├── routes
// │   ├── authRouter.go
// │   └── userRouter.go
// └── utils
//     ├── database.go
//     └── logger.go
```

* The "models" directory contains the struct includes fields such as ID, Username and Password.

* The "utils" package contains helper functions for managing the database connection, error handling, and logging.

* The "middleware" package includes the auth_middleware & cors_middleware to protect certain routes that require authentication.

* The "controllers" package contains the UserController, which handles CRUD operations for User items. It interacts with the database and uses the AuthMiddleware to protect certain routes.

You can start building the User API with the project structure in place. You'll create routes for user registration, user login, and CRUD operations for managing User items.

Here's a high-level overview of the project implementation:

    [x] Define the User struct in the "models" package.

    [x] Implement database functions in the "utils/database.go" to connect to SQLite and perform CRUD operations for User items.

    [x] Create the AuthMiddleware in the "middleware/auth_middleware.go" to handle JWT-based authentication.

    [x] Build the UserController in the "controllers/user_controller.go" to handle CRUD operations for User items. The controller will use the database functions, authentication middleware, and error-handling utilities.

    [x] Define routes in the "main.go" file to handle user registration, user login, and CRUD operations for User items. Use the UserController methods as route handlers and apply the AuthMiddleware to protect certain routes.

    [ ] Implement a persistent authentification

    [ ] A method of identifying users via OAuth2 (eg Google / Apple / Facebook / Github.). In this case, the client processes itself this identification and warns the application server if successful

    [x] Import API to Swagger.io

Possibly implementable to the code while handling:

    [x] The application server should answer the call http://localhost:8080/about.json.
