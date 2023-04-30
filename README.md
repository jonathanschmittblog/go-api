# Example API

This is an example of a basic Golang API.

To start the docker compose services, run

```bash
$ docker-compose up -d
```

To start the application, run
```bash
# Start application
$ cd cmd
$ go run main.go

# If the application don't start properly, run those commands and try again
$ go mod tidy
$ go mod vendor
```

### Troubleshooting

While running the project on a Windows environment it can occur an error like this:

```/bin/bash bad interpreter: no such file or directory```

This is caused by the CRLF line break. The postgres container uses the file ```/database/postgres-init-user-db.sh``` to automatically starts the database on the container initialization. As this is a bash file that will be executed inside the container (who is running a linux image) it needs to be settled with LF line break. This can be easy switched by opening this file inside VSCode and toggling to LF in the bottom right corner bar.