load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//extras:embed_data.bzl", "go_embed_data")

go_embed_data(
    name = "templates",
    srcs = glob(["*.html"]),
    package = "templates",
)

# keep
go_library(
    name = "go_default_library",
    srcs = [":templates"],
    importpath = "github.com/q3k/scarab/templates",
    visibility = ["//visibility:public"],
)
