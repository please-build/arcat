go_binary(
    name = "arcat",
    srcs = ["main.go"],
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go:cli-init",
        "//third_party/go:logging",
        "//ar",
        "//tar",
        "//unzip",
        "//zip",
    ],
)

genrule(
    name = "version",
    srcs = ["VERSION"],
    outs = ["version.build_defs"],
    cmd = "echo VERSION = \\\"$(cat $SRCS)\\\" > $OUT",
    visibility = ["//package:all"],
)
