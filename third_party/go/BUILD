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
    tags = [
        "netgo",
        "osusergo",
    ],
    version = "1.27.1",
)

go_repo(
    name = "logging",
    install = ["."],
    licences = ["BSD-3-Clause"],
    module = "gopkg.in/op/go-logging.v1",
    version = "v1.0.0-20160211212156-b2cb9fa56473",
)

go_repo(
    name = "go-flags",
    install = ["."],
    licences = ["BSD-3-Clause"],
    module = "github.com/thought-machine/go-flags",
    version = "v1.7.0",
)

go_repo(
    name = "cli-init",
    install = ["..."],
    licences = ["Apache-2.0"],
    module = "github.com/peterebden/go-cli-init/v5",
    version = "v5.2.1",
)

go_repo(
    name = "xsys",
    install = ["..."],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/sys",
    version = "v0.48.0",
)

go_repo(
    name = "testify",
    install = [
        "assert",
        "assert/yaml",
        "require",
    ],
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "v1.12.1",
)

go_repo(
    name = "difflib",
    install = ["difflib"],
    licences = ["BSD-3-Clause"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.0",
)

go_repo(
    name = "yaml",
    install = ["."],
    licences = ["MIT"],
    module = "go.yaml.in/yaml/v3",
    version = "v3.0.5",
)

go_repo(
    name = "spew",
    install = ["spew"],
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.1",
)

go_repo(
    name = "ar",
    install = ["."],
    licences = ["MIT"],
    module = "github.com/please-build/ar",
    version = "v0.0.0-20251128102243-20fe5956df94",
)

go_repo(
    name = "ulikunitz_xz",
    install = [
        ".",
        "lzma",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/ulikunitz/xz",
    version = "v0.5.17",
)

go_repo(
    name = "xi2_xz",
    install = ["."],
    licences = ["Public Domain"],
    module = "github.com/xi2/xz",
    version = "v0.0.0-20171230120015-48954b6210f8",
)

go_repo(
    name = "zstd",
    install = [
        "zstd/...",
        "s2",
        "snappy",
        ".",
        "fse",
        "huff0",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/klauspost/compress",
    version = "v1.20.0",
)

go_repo(
    name = "humanize",
    install = ["."],
    licences = ["MIT"],
    module = "github.com/dustin/go-humanize",
    version = "v1.1.0",
)

go_repo(
    name = "xcrypto",
    install = [
        "cast5",
        "curve25519",
        "ed25519",
        "ssh/terminal",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/crypto",
    version = "v0.57.0",
)

go_repo(
    name = "xterm",
    install = ["."],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/term",
    version = "v0.46.0",
)
