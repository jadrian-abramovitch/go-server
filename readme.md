# Simple go server using gin, just messing around

## Commands for testing
### SETUP:
> Install sqlite3
### Build server: 
> go build 
### Run server:
>./go-server
### Build [Tailwind](https://tailwindcss.com/docs/installation):
>npx tailwindcss -i ./src/main.css -o ./src/output.css --watch

### curl commands
* > curl localhost:8080/courses
* > curl localhost:8080/course/1
* > curl localhost:8080/newCourse -d @payload.json
* > curl localhost:8080/course/3 -d @updated-payload.json -X "PATCH"
* > curl localhost:8080/course/2 -X "DELETE" 