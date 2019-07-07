load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@build_bazel_rules_nodejs//:defs.bzl", "rollup_bundle")

#load("@build_stack_rules_proto//:compile.bzl", "proto_compile")
load("@io_bazel_rules_go//extras:embed_data.bzl", "go_embed_data")

#closure_grpc_web_library(
#    name = "common_grpc_web_proto",
#    deps = ["//proto/common:common_proto"],
#)

rollup_bundle(
    name = "bundle",
    entry_point = "app.js",
    srcs = [
        "//proto/common:common_js_proto",
    ],
    deps = [
        "@npm//vue",
        "@npm//grpc-web",
        "@npm//google-protobuf",
    ],
)

go_embed_data(
    name = "js",
    srcs = [":bundle"],
    package = "js",
    flatten = True,
)

# keep
go_library(
    name = "go_default_library",
    srcs = [":js"],
    importpath = "github.com/q3k/scarab/js",
    visibility = ["//visibility:public"],
)
