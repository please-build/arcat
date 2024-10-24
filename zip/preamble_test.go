package zip

import (
	"crypto/sha256"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreamble(t *testing.T) {
	for _, tc := range []struct{
		ZipFile          string
		PreambleChecksum string
	}{
		{
			ZipFile:          "kitten_preamble.zip",
			PreambleChecksum: "cd41527591eeac94f556fefa7505e54d9651e783cff0b8498ae0d65b55ffd416", // kitten.png
		},
		{
			ZipFile:          "shebang_oneline_preamble.zip",
			PreambleChecksum: "378f13f639d3ec2cbc18878b2d244adfaf4f9969464dd05aff0d5d9ed998d592", // shebang_oneline.txt
		},
		{
			ZipFile:          "shebang_twolines_preamble.zip",
			PreambleChecksum: "46533b2dfa35ad537d3561ebee0c7af8941bc65363c1b188e1be6eaf79e9138c", // shebang_twolines.txt
		},
		{
			ZipFile:          "zip_preamble.zip",
			PreambleChecksum: "038a57f3f807fa91bdd30239b9711fccf0d782fe2f036e03211852237e94d24c", // another.zip
		},
		{
			ZipFile:          "arcat_preamble.zip",
			PreambleChecksum: "46533b2dfa35ad537d3561ebee0c7af8941bc65363c1b188e1be6eaf79e9138c", // shebang_twolines.txt
		},
		{
			ZipFile:          "empty.zip",
			PreambleChecksum: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", // empty string
		},
	} {
		t.Run(tc.ZipFile, func(t *testing.T) {
			r := require.New(t)

			pr, err := Preamble("zip/test_data_4/" + tc.ZipFile)
			r.NoError(err)
			h := sha256.New()
			_, err = io.Copy(h, pr)
			r.NoError(err)
			preambleChecksum := fmt.Sprintf("%x", h.Sum(nil))
			r.Equal(tc.PreambleChecksum, preambleChecksum)
			err = pr.Close()
			r.NoError(err)
		})
	}
}
