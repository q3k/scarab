load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@build_bazel_rules_nodejs//:defs.bzl", "rollup_bundle")

load("@io_bazel_rules_go//extras:embed_data.bzl", "go_embed_data")

rollup_bundle(
    name = "bundle",
    entry_point = "app.js",
    srcs = [
        "remote.js",
        "components.js",
        "//proto/common:common_js_proto",
    ],
    deps = [
        "@npm//vue",
        "@npm//vue-router",
        "@npm//grpc-web",
        "@npm//google-protobuf",
    ],
)

go_embed_data(
    name = "js",
    srcs = [":bundle.es2015.js"],
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
