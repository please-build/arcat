package(default_visibility = ["PUBLIC"])

go_toolchain(
    name = "toolchain",
    architectures = [
        "darwin_amd64",
        "darwin_arm64",
        "freebsd_amd64",
        "linux_amd64",
        "linux_arm64",
    ],
    version = "1.17.5",
    tags = [
        "osusergo",
        "netgo",
    ],
)

go_module(
    name = "logging",
    licences = ["BSD-3-Clause"],
    module = "gopkg.in/op/go-logging.v1",
    version = "v1.0.0-20160211212156-b2cb9fa56473",
    visibility = ["PUBLIC"],
)

go_module(
    name = "go-flags",
    licences = ["BSD-3-Clause"],
    module = "github.com/thought-machine/go-flags",
    version = "v1.6.0",
    visibility = ["PUBLIC"],
    deps = [":xsys"],
)

go_module(
    name = "cli-init",
    install = ["..."],
    module = "github.com/peterebden/go-cli-init/v5",
    version = "v5.2.0",
    deps = [
        ":go-flags",
        ":humanize",
        ":logging",
        ":xcrypto",
    ],
)

go_module(
    name = "xsys",
    install = ["..."],
    module = "golang.org/x/sys",
    version = "v0.0.0-20210823070655-63515b42dcdf",
)

go_module(
    name = "testify",
    install = [
        "assert",
        "require",
    ],
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "v1.7.0",
    visibility = ["PUBLIC"],
    deps = [
        ":difflib",
        ":spew",
        ":yaml",
    ],
)

go_module(
    name = "difflib",
    install = ["difflib"],
    licences = ["BSD-3-Clause"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.0",
    visibility = ["PUBLIC"],
)

go_module(
    name = "yaml",
    licences = ["MIT"],
    module = "gopkg.in/yaml.v3",
    version = "v3.0.0-20210107192922-496545a6307b",
    visibility = ["PUBLIC"],
)

go_module(
    name = "spew",
    install = ["spew"],
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.1",
    visibility = ["PUBLIC"],
)

go_module(
    name = "ar",
    licences = ["MIT"],
    module = "github.com/peterebden/ar",
    version = "v0.0.0-20181115090543-a0ae3a11a518",
    visibility = ["PUBLIC"],
)

go_module(
    name = "ulikunitz_xz",
    install = [
        "internal/...",
        ".",
        "lzma",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/ulikunitz/xz",
    version = "v0.5.10",
    visibility = ["PUBLIC"],
)

go_module(
    name = "xi2_xz",
    licences = ["Public Domain"],
    module = "github.com/xi2/xz",
    version = "v0.0.0-20171230120015-48954b6210f8",
    visibility = ["PUBLIC"],
)

go_module(
    name = "zstd",
    install = [
        "zstd/...",
        "s2",
        "snappy",
        ".",
        "fse",
        "huff0",
        "internal/snapref",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/klauspost/compress",
    version = "v1.13.6",
    visibility = ["PUBLIC"],
)

go_module(
    name = "humanize",
    licences = ["MIT"],
    module = "github.com/dustin/go-humanize",
    version = "v1.0.0",
    visibility = ["PUBLIC"],
)

go_module(
    name = "xcrypto",
    install = [
        "cast5",
        "curve25519",
        "curve25519/internal/field",
        "ed25519",
        "ssh/terminal",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/crypto",
    version = "v0.0.0-20210920023735-84f357641f63",
    visibility = ["PUBLIC"],
    deps = [":xterm"],
)

go_module(
    name = "xterm",
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/term",
    version = "v0.0.0-20210615171337-6886f2dfbf5b",
    visibility = ["PUBLIC"],
    deps = [":xsys"],
)
