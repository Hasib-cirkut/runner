### Runner is a web based code runner. 
Stack: Go(gin), Vue3

Supported Languages: Javascript (till now. more coming)

### Run locally
- Make sure you have go and docker installed

#### Backend
- Make sure docker is up and running
- run `go run main.go` to start the backend server.
- Backend server will be listening on port `:8080` by default

##### Developer exprience
To test routes, this project has swagger installed

After running the server, 
visit http://localhost:port/swagger/index.html#/ 

- To update the swagger documentation, run
```bash
swag init
```