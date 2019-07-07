load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/q3k/scarab
gazelle(name = "gazelle")

go_library(
    name = "go_default_library",
    srcs = [
        "generate.go",
        "http.go",
        "scarab.go",
        "storage.go",
        "storage_leveldb.go",
    ],
    importpath = "github.com/q3k/scarab",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/common:go_default_library",
        "//proto/storage:go_default_library",
        "//templates:go_default_library",
        "//js:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_improbable_eng_grpc_web//go/grpcweb:go_default_library",
        "@com_github_syndtr_goleveldb//leveldb:go_default_library",
        "@com_github_syndtr_goleveldb//leveldb/util:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

alias(
    name = "tsconfig.json",
    actual = "//js:tsconfig.json",
    visibility = ["//visibility:public"],
)

load("@build_bazel_rules_nodejs//:defs.bzl", "nodejs_binary")

nodejs_binary(
    name = "rollup",
    data = [
        "@npm//is-builtin-module",
        "@npm//rollup",
        "@npm//rollup-plugin-alias",
        "@npm//rollup-plugin-amd",
        "@npm//rollup-plugin-commonjs",
        "@npm//rollup-plugin-json",
        "@npm//rollup-plugin-node-resolve",
        "@npm//rollup-plugin-replace",
        "@npm//rollup-plugin-sourcemaps",
        "@npm//rollup-plugin-vue",
    ],
    entry_point = "@npm//:node_modules/rollup/bin/rollup",
    install_source_map_support = False,
    visibility = ["//visibility:public"],
)
