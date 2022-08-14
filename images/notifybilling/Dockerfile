FROM public.ecr.aws/lambda/provided:al2-arm64 as build
WORKDIR /go/src/project
RUN yum update golang & yum install -y golang
ADD go.mod go.sum main.go ./
RUN GOOS=linux GOARCH=arm64 go build -o /main ./main.go

FROM public.ecr.aws/lambda/provided:al2-arm64
COPY --from=build /main /main
RUN chmod 755 /main
ENTRYPOINT [ "/main" ]
