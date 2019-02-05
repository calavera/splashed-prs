FROM golang:1.11 AS builder
WORKDIR /go/src/github.com/calavera/splashed_prs
ADD . /go/src/github.com/calavera/splashed_prs
RUN make build

FROM scratch
LABEL "com.github.actions.name" "Splashed PRs"
LABEL "com.github.actions.description" "Add random beautiful photos from Unsplash to your Pull Requests"
LABEL "com.github.actions.icon" "camera"
LABEL "com.github.actions.color" "purple"

COPY --from=builder /go/src/github.com/calavera/splashed_prs/dist/splashed_prs /
ENTRYPOINT ["/splashed_prs"]
