package crypt

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xvzf/htw-crypto-project/pkg/crypt"
	"github.com/xvzf/htw-crypto-project/pkg/image"
)

// {{{ taken from https://poestories.com/read/goldbug
const blindText = `MANY years ago, I contracted an intimacy with a Mr. William Legrand. He was of an ancient Huguenot family, and had once been wealthy; but a series of misfortunes had reduced him to want. To avoid the mortification consequent upon his disasters, he left New Orleans, the city of his forefathers, and took up his residence at Sullivan's Island, near Charleston, South Carolina.
This Island is a very singular one. It consists of little else than the sea sand, and is about three miles long. Its breadth at no point exceeds a quarter of a mile. It is separated from the main land by a scarcely perceptible creek, oozing its way through a wilderness of reeds and slime, a favorite resort of the marsh-hen. The vegetation, as might be supposed, is scant, or at least dwarfish. No trees of any magnitude are to be seen. Near the western extremity, where Fort Moultrie stands, and where are some miserable frame buildings, tenanted, during summer, by the fugitives from Charleston dust and fever, may be found, indeed, the bristly palmetto; but the whole island, with the exception of this western point, and a line of hard, white beach on the seacoast, is covered with a dense undergrowth of the sweet myrtle, so much prized by the horticulturists of England. The shrub here often attains the height of fifteen or twenty feet, and forms an almost impenetrable coppice, burthening the air with its fragrance.
In the inmost recesses of this coppice, not far from the eastern or more remote end of the island, Legrand had built himself a small hut, which he occupied when I first, by mere accident, made his acquaintance. This soon ripened into friendship --for there was much in the recluse to excite interest and esteem. I found him well educated, with unusual powers of mind, but infected with misanthropy, and subject to perverse moods of alternate enthusiasm and melancholy. He had with him many books, but rarely employed them. His chief amusements were gunning and fishing, or sauntering along the beach and through the myrtles, in quest of shells or entomological specimens;-his collection of the latter might have been envied by a Swammerdamm. In these excursions he was usually accompanied by an old negro, called Jupiter, who had been manumitted before the reverses of the family, but who could be induced, neither by threats nor by promises, to abandon what he considered his right of attendance upon the footsteps of his young "Massa Will." It is not improbable that the relatives of Legrand, conceiving him to be somewhat unsettled in intellect, had contrived to instil this obstinacy into Jupiter, with a view to the supervision and guardianship of the wanderer.
The winters in the latitude of Sullivan's Island are seldom very severe, and in the fall of the year it is a rare event indeed when a fire is considered necessary. About the middle of October, 18--, there occurred, however, a day of remarkable chilliness. Just before sunset I scrambled my way through the evergreens to the hut of my friend, whom I had not visited for several weeks --my residence being, at that time, in Charleston, a distance of nine miles from the Island, while the facilities of passage and re-passage were very far behind those of the present day. Upon reaching the hut I rapped, as was my custom, and getting no reply, sought for the key where I knew it was secreted, unlocked the door and went in. A fine fire was blazing upon the hearth. It was a novelty, and by no means an ungrateful one. I threw off an overcoat, took an arm-chair by the crackling logs, and awaited patiently the arrival of my hosts.
Soon after dark they arrived, and gave me a most cordial welcome. Jupiter, grinning from ear to ear, bustled about to prepare some marsh-hens for supper. Legrand was in one of his fits --how else shall I term them? --of enthusiasm. He had found an unknown bivalve, forming a new genus, and, more than this, he had hunted down and secured, with Jupiter's assistance, a scarabaeus which he believed to be totally new, but in respect to which he wished to have my opinion on the morrow.
`

// }}}

var (
	cipher          *crypt.Container
	blindText256Enc []crypt.PixelPosition
)

func init() {

	var err error

	i := image.Mock()
	for !image.CheckAccept(i) {
		i = image.Mock()
	}

	// Build new crypt engine
	cipher, err = crypt.New(i)
	if err != nil {
		log.Fatal(err)
	}

	blindText256Enc = make([]crypt.PixelPosition, len(blindText)*256)
	for c := 0; c < 256; c++ {
		enc, _ := cipher.Encrypt(blindText)
		blindText256Enc = append(blindText256Enc, enc...)
	}
}

func TestLoad(t *testing.T) {
	var sum int
	var a *Analyse

	a = Load(blindText256Enc)
	sum = 0

	assert.NotEmpty(t, a.Frequency)

	for _, v := range a.Frequency {
		sum += v
	}

	assert.Equal(t, len(blindText256Enc), sum)
	assert.Equal(t, a.Total, sum)
}

func TestAnalysis_ExtractGroups(t *testing.T) {

	// Build blindText set
	blindTextSet := make(map[uint8]bool)
	for _, v := range blindText {
		if _, ok := blindTextSet[uint8(v)]; ok {
			continue
		}
		blindTextSet[uint8(v)] = true
	}

	// a := Load(blindText256Enc)
	// TODO
	// groups := a.ExtractGroups()
	// assert.LessOrEqual(t, len(blindTextSet), len(groups))
}
