go_library(
    name = "unzip",
    srcs = ["unzip.go"],
    visibility = ["//:all"],
    deps = [
        "//third_party/go:xi2_xz",
        "//third_party/go:zstd",
    ],
)

go_test(
    name = "unzip_test",
    srcs = ["unzip_test.go"],
    data = ["test_data"],
    deps = [
        ":unzip",
        "//third_party/go:testify",
    ],
)
