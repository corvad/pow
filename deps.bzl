load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=",
        version = "v1.5.3",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        importpath = "github.com/grpc-ecosystem/grpc-gateway/v2",
        sum = "h1:RtRsiaGvWxcwd8y0BiRYN6oBRLZJLOeWNMmBWY5kqxg=",
        version = "v2.19.0",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:gVBHo4v9F7DoYwYl0pZLBhEihJJwrjz3Qhjs3vYqJgw=",
        version = "v0.0.0-20240125205218-1f4bbc51befe",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:bXHQvfJcDKVfmjtzPgbF9sHVT9kgRDSLx/Y9fJCaM2Q=",
        version = "v0.0.0-20240125205218-1f4bbc51befe",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:Z5Iec2pjwb+LEOqzpB2MR12/eKFhDPhuqW91O+4bwUk=",
        version = "v1.60.1",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:kf48UfkdszAmrW0rYOYGPoU+x1ToPOQzNbwMBwVmgg=",
        version = "v1.32.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:aCL9BSgETF1k+blQaYUBx9hJ9LOGP3gAVemcZlf1Kpo=",
        version = "v0.20.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:xWw16ngr6ZMtmxDyKyIgsE93KNKz5HKmMa3b8ALHidU=",
        version = "v0.16.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:ScX5w1eTa3QqT8oi6+ziP7dTV1S2+ALU0bI+0zXKWiQ=",
        version = "v0.14.0",
    )
