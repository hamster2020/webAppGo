# webAppGo
###This is an simple web application implemented in solely in Go.

###The features included in the application include:
* **Authentication**
  * password hashing with the *golang.org/x/crypto/bcrypt* package
  * `password` hash stored in `user` db table
* **Logins**
  * timeouts clients that reach specified number of failed login attempts to a given `username` within a specified timeframe
  * timeouts clients that reach specified number of failed login attempts to multiple `username` coming from a same ip address within a specified timeframe
  * `ip`, `username`, `timestamp`, and `attempt` are stored in login db table
* **Session management**
  * sessions handled with secure cookies
  * `sessionid`, `userid`, and `timestamp` stored in `session` db table
* **Pages**
  * allows **logged in** client to create, read, and update `pages`
  * a `page` is saved and retrieved from both the `pages` db table and a local cache
  * `title`, `body`, and `timestamp` stored in `pages` db table
* **Files**
  * allows logged in client to upload, download, and list all files.
  * files stored in environment defined file path
* **User Accounts**
  * Allows a given client to create a `user` by signing up, modifying the `user` information if logged in, and then delete a given `user` account if logged in.
  * `userid`, `username`, `firstname`, `lastname`, `email`, and hash of `password` are stored in `users` db table
* **JSON API**
  * allows any client to create, read, update, and delete a specified `page` (CRUD) through the `/page/` endpoint
  * allows any client to list all `pages` through the `/pages` endpoint


### The application was also designed with the following features:
* **Dependency Injection**
  This application is broken up into several sub-packages, where the top-level sub-package contains several interfaces that are passed into the web and api sub-packages via *dependency injection*. This keeps the code flexible, decoupled, and easy to test different components in isolation.
* **Toggle-able Databases**
  The db is handled with dependency injection by using a `Datastore` interface that is passed down into the `web` and `api` sub-packages. Two other sub-packages, `sqlite` and `postgres`, both implement the `Datastore` interface, which allows the developer to easily toggle between two different databases.   
* **Environment Defined Variables**
  The environment variable is a struct which contains the different interfaces such as the `Datastore` and `Cachestore`, as well as other environment-wide variables such as the file path and logging settings.
* **Custom Logging**
  A custom logger was built around Go's `log` package, but also includes leveled logging, detailed prefixing, and simple controls to specify the path to vase to a file or merely print to the console.
* **Testing**
  Go's OOTB testing was implemented in both the `web` and `api` packages in this project. Since the db and cache were decoupled through *dependency injection*, testing in these sub-packages were possible with completely different sources, created just for testing, called `mockDB` and `mockCache`.
