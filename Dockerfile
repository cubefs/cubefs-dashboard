FROM golang:alpine AS backend

COPY ./backend /backend
COPY ./depends /depends

WORKDIR /backend

RUN go build \
    -ldflags "-X github.com/cubefs/cubefs/proto.Version=$(git describe --abbrev=0 --tags 2>/dev/null) \
    -X github.com/cubefs/cubefs/proto.CommitID=$(git rev-parse HEAD 2>/dev/null) \
    -X github.com/cubefs/cubefs/proto.BranchName=$(git rev-parse --abbrev-ref HEAD 2>/dev/null) \
    -X 'github.com/cubefs/cubefs/proto.BuildTime=$(date +%Y-%m-%d\ %H:%M)'" \
    -o /bin/cfs-gui \
    *.go

FROM node:14 AS frontend

COPY ./frontend /frontend

WORKDIR /frontend

RUN npm install && npm run build

FROM alpine:latest

COPY --from=backend /bin/cfs-gui /app/cfs-gui
COPY --from=frontend /frontend/dist /app/dist

WORKDIR /app

CMD ["/app/cfs-gui", "-c", "/app/config.yml"]
