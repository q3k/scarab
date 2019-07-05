workspace(
    name = "scarab",
    managed_directories = {"@npm": ["node_modules"]},
)

# JS / TS / Node

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "58ccb7a1a1cf9c6f4e36f79e7839f226eba4a53c7537782a86a94e1ca05ad269",
    urls = ["https://github.com/q3k/rules_nodejs/archive/ed9bda56fd661e7dfe1e5513d1232babb5c72ac5.zip"],
    strip_prefix = "rules_nodejs-ed9bda56fd661e7dfe1e5513d1232babb5c72ac5",
)

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories", "yarn_install")

node_repositories(package_json = ["//js:package.json"])

yarn_install(
    name = "npm",
    package_json = "//js:package.json",
    yarn_lock = "//js:yarn.lock",
)

load("@npm//:install_bazel_dependencies.bzl", "install_bazel_dependencies")

install_bazel_dependencies()

load("@npm_bazel_typescript//:index.bzl", "ts_setup_workspace")

ts_setup_workspace()

# rules_proto (for JS proto)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    urls = ["https://github.com/stackb/rules_proto/archive/b93b544f851fdcd3fc5c3d47aee3b7ca158a8841.tar.gz"],
    sha256 = "c62f0b442e82a6152fcd5b1c0b7c4028233a9e314078952b6b04253421d56d61",
    strip_prefix = "rules_proto-b93b544f851fdcd3fc5c3d47aee3b7ca158a8841",
)

load("@build_stack_rules_proto//github.com/grpc/grpc-web:deps.bzl", "commonjs_grpc_compile")

commonjs_grpc_compile()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)

http_archive(
    name = "grpc_grpc_web",
    urls = ["https://github.com/grpc/grpc-web/archive/a3174cb1d0dedfc130bb5b8133e60c78e78ac722.zip"],
    sha256 = "c6c680c2e01516e618ef3a467788374628b6cf08c20f222244afbdd23d1b6222",
    strip_prefix = "grpc-web-a3174cb1d0dedfc130bb5b8133e60c78e78ac722",
)

# Go

http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/0.18.6/rules_go-0.18.6.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/0.18.6/rules_go-0.18.6.tar.gz",
    ],
    sha256 = "f04d2373bcaf8aa09bccb08a98a57e721306c8f6043a2a0ee610fd6853dcde3d",
)

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"],
    sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_fsnotify_fsnotify",
    importpath = "github.com/fsnotify/fsnotify",
    tag = "v1.4.7",
)

go_repository(
    name = "com_github_gogo_protobuf",
    importpath = "github.com/gogo/protobuf",
    tag = "v1.2.1",
)

go_repository(
    name = "com_github_golang_glog",
    commit = "23def4e6c14b",
    importpath = "github.com/golang/glog",
)

go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    tag = "v1.3.1",
)

go_repository(
    name = "com_github_golang_snappy",
    commit = "2e65f85255db",
    importpath = "github.com/golang/snappy",
)

go_repository(
    name = "com_github_hpcloud_tail",
    importpath = "github.com/hpcloud/tail",
    tag = "v1.0.0",
)

go_repository(
    name = "com_github_kevinburke_go_bindata",
    importpath = "github.com/kevinburke/go-bindata",
    tag = "v3.13.0",
)

go_repository(
    name = "com_github_kisielk_errcheck",
    importpath = "github.com/kisielk/errcheck",
    tag = "v1.1.0",
)

go_repository(
    name = "com_github_kisielk_gotool",
    importpath = "github.com/kisielk/gotool",
    tag = "v1.0.0",
)

go_repository(
    name = "com_github_onsi_ginkgo",
    importpath = "github.com/onsi/ginkgo",
    tag = "v1.7.0",
)

go_repository(
    name = "com_github_onsi_gomega",
    importpath = "github.com/onsi/gomega",
    tag = "v1.4.3",
)

go_repository(
    name = "com_github_syndtr_goleveldb",
    importpath = "github.com/syndtr/goleveldb",
    tag = "v1.0.0",
)

go_repository(
    name = "in_gopkg_check_v1",
    commit = "20d25e280405",
    importpath = "gopkg.in/check.v1",
)

go_repository(
    name = "in_gopkg_fsnotify_v1",
    importpath = "gopkg.in/fsnotify.v1",
    tag = "v1.4.7",
)

go_repository(
    name = "in_gopkg_tomb_v1",
    commit = "dd632973f1e7",
    importpath = "gopkg.in/tomb.v1",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    tag = "v2.2.1",
)

go_repository(
    name = "org_golang_x_net",
    commit = "161cd47e91fd",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_sync",
    commit = "1d60e4601c6f",
    importpath = "golang.org/x/sync",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "d0be0721c37e",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    tag = "v0.3.0",
)

go_repository(
    name = "org_golang_x_tools",
    commit = "07fd8470d635",
    importpath = "golang.org/x/tools",
)

go_repository(
    name = "com_github_improbable_eng_grpc_web",
    commit = "5b222b4af3764c826506204bb4925dd1a704a414",
    importpath = "github.com/improbable-eng/grpc-web",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "ae1634f6a98965ded3b8789c626cb4e0bd78c3de",
    importpath = "github.com/gorilla/websocket",
)

go_repository(
    name = "com_github_rs_cors",
    commit = "33ffc0734c60052d8b609f1e2d8ebbadc6cd80d8",
    importpath = "github.com/rs/cors",
)
