go_library(
    name = "zip",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = [
        "//:all",
        "//tools/please_pex/pex",
    ],
    deps = [
        "//third_party/go:logging",
        "//third_party/go:testify",
    ],
)

go_test(
    name = "preamble_test",
    srcs = ["preamble_test.go"],
    data = [
        "test_data_4",
    ],
    deps = [
        ":zip",
    ],
)

go_test(
    name = "writer_test",
    srcs = ["writer_test.go"],
    data = [
        "test_data",
        "test_data_2",
        "test_data_3",
    ],
    deps = [
        ":zip",
    ],
)
