package globe

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	outputImages = flag.Bool("images", false, "output images produced in test")
)

func AssertPNGMD5(t *testing.T, g *Globe, expect string) {
	m := g.Image(1024)
	h := md5.New()
	var w io.Writer = h
	if *outputImages {
		filename := fmt.Sprintf("%s.png", t.Name())
		f, err := os.Create(filename)
		require.NoError(t, err)
		defer f.Close()
		w = io.MultiWriter(h, f)
	}
	err := png.Encode(w, m)
	require.NoError(t, err)
	assert.Equal(t, expect, hex.EncodeToString(h.Sum(nil)))
}

func TestGraticule(t *testing.T) {
	g := New()
	g.DrawGraticule(10.0)
	AssertPNGMD5(t, g, "5861947654cabd808bc3c75ab8018576")
}

func TestGraticuleCenterOn(t *testing.T) {
	g := New()
	g.DrawGraticule(10.0)
	g.CenterOn(60, 5)
	AssertPNGMD5(t, g, "04c8271e8d48ea580a6d340d9be7a261")
}

func TestDrawDots(t *testing.T) {
	lat, lng := -89.5, -179.5
	g := New()
	for i := 0; i < 180; i++ {
		g.DrawDot(lat, lng, 0.1)
		lat += 1.0
		lng += 2.0
	}
	g.CenterOn(0, 0)
	AssertPNGMD5(t, g, "200588de765c11b6b4136b8df36a698c")
}

func TestDrawLand(t *testing.T) {
	g := New()
	g.DrawLandBoundaries()
	g.CenterOn(51.453349, -2.588323)
	AssertPNGMD5(t, g, "d7667370d3395cb7d987cee444e59b17")
}

func TestDrawCountries(t *testing.T) {
	g := New()
	g.DrawCountryBoundaries()
	g.CenterOn(40.645423, -73.903879)
	AssertPNGMD5(t, g, "139333c753ce3af5077a89b3c4bd8cdb")
}

func TestLine(t *testing.T) {
	g := New()
	g.DrawLine(51.453349, -2.588323, 40.645423, -73.903879)
	g.CenterOn(30, -37)
	AssertPNGMD5(t, g, "a2c7d6a4de4c652a5538544950fcde92")
}

func TestRect(t *testing.T) {
	g := New()
	g.DrawRect(41.897209, 12.500285, 55.782693, 37.615993)
	g.CenterOn(48, 25)
	AssertPNGMD5(t, g, "a1f7d08345c78e484612e211f8dc5e4b")
}

func TestCartestian(t *testing.T) {
	x, y, z := cartestian(42, -163)
	assert.Equal(t, -0.7106729309733519, x)
	assert.Equal(t, -0.2172745194807066, y)
	assert.Equal(t, -0.6691306063588582, z)
}

func TestHaversine(t *testing.T) {
	cases := []struct {
		Lat1, Lng1 float64
		Lat2, Lng2 float64
		Distance   float64
	}{
		{61.76543033984557, 37.67770367266306, -7.156018382728888, 59.24161915865653, 7886.581201923186},
		{21.940661070124392, -27.130501054344364, -43.39006318567746, -156.37067308170856, 14798.28914867192},
		{-23.464218263492626, -145.09097319078558, 38.849724010954006, 5.4765462607435325, 16755.646528079917},
		{-13.808977659204544, -102.86500587034503, -3.5672654028449813, -65.49905724108126, 4254.824289188308},
		{-24.44353568488961, -78.10770557503974, -34.25616705840568, 64.47048333127782, 12385.613939186398},
		{-16.156214450511584, -106.85272440696376, 46.46742324464643, 25.442379385568103, 14489.882829355967},
		{-23.94345515581749, -74.4788719606111, -35.93270748905276, 90.92629279858028, 13182.635916739864},
		{23.168612332699823, 131.52060468056197, -70.62954835525073, 8.575310178000308, 13619.095429569754},
		{12.38659860174893, -123.00182001175405, -57.25569219003512, 171.08698278980825, 9783.195280980919},
		{-61.85524721381324, 34.13109516590251, -23.387935711306838, 69.12885144712033, 5025.638176612954},
	}
	for _, c := range cases {
		d := haversine(c.Lat1, c.Lng1, c.Lat2, c.Lng2)
		assert.InDelta(t, c.Distance, d, 0.00001)
	}
}

func TestIntermediate(t *testing.T) {
	cases := []struct {
		Lat1, Lng1 float64
		Lat2, Lng2 float64
		Fraction   float64
		Lat3, Lng3 float64
	}{
		{61.76543033984557, 37.67770367266306, 59.579238015627766, 30.987111615865974, 0.4246374970712657, 60.878036362889496, 34.725267292961924},
		{-60.31092387052814, 67.25630623215943, -59.81972184198988, 68.70037420777328, 0.30091186058528707, -60.16477620153165, 67.69537164211135},
		{38.849724010954006, 5.4765462607435325, 39.53243555910779, 9.805266232872793, 0.31805817433032985, 39.08419448208641, 6.8441006009550795},
		{-25.71733560446036, -11.19965583512763, -27.200820830482833, -4.581209185196006, 0.21855305259276428, -26.067400191907335, -9.76790788460991},
		{-16.156214450511584, -106.85272440696376, -23.127177330958297, -110.46750918167837, 0.29311424455385804, -18.206742475084216, -107.880715918046},
		{30.341049367778716, -73.05027711973504, 32.139284152631184, -64.19154811423957, 0.6967191657466347, 31.658637661505345, -66.91338783046506},
		{-70.62954835525073, 8.575310178000308, -67.20824427099608, 20.46555773756403, 0.9752416188605784, -67.30215247794456, 20.213952362608747},
		{10.930445007176502, -151.3966955854061, 16.719163720796363, -149.04192853179413, 0.30152268100656, 12.67816445562967, -150.69863258670347},
		{4.715016392832226, -117.6241542542261, 2.305834742279346, -118.31051998611643, 0.4231522015718281, 3.6956282925427137, -117.91503132483825},
		{-29.532620146295827, 11.010857526253858, -30.707453241094573, 19.098736201034995, 0.3618054804803169, -30.01478484795552, 13.914487010595325},
	}
	for _, c := range cases {
		lat, lng := intermediate(c.Lat1, c.Lng1, c.Lat2, c.Lng2, c.Fraction)
		assert.InDelta(t, c.Lat3, lat, 0.00001, "latitude error")
		assert.InDelta(t, c.Lng3, lng, 0.00001, "longitude error")
	}
}

func TestDestination(t *testing.T) {
	cases := []struct {
		Lat1, Lng1 float64
		Bearing    float64
		Distance   float64
		Lat2, Lng2 float64
	}{
		{61.76543033984557, 37.67770367266306, 239.24161915865656, 437.7141871869802, 59.579238015627766, 30.987111615865974},
		{21.940661070124392, -27.130501054344364, 23.62932691829144, 156.51925473279124, 23.229104720261365, -26.516583832765832},
		{-23.464218263492626, -145.09097319078558, 185.47654626074353, 813.6399609900968, -30.7458989905092, -145.90134837100578},
		{-13.808977659204544, -102.86500587034503, 114.50094275891874, 468.8898449024232, -15.52459863254281, -98.88297813500867},
		{-24.44353568488961, -78.10770557503974, 244.47048333127785, 218.5530525927643, -25.277919578585674, -80.06910413078825},
		{-16.156214450511584, -106.85272440696376, 205.44237938556813, 862.4914374478864, -23.127177330958297, -110.46750918167837},
		{-23.94345515581749, -74.4788719606111, 270.9262927985803, 206.5826619136986, -23.900059713658717, -76.51076302260704},
		{23.168612332699823, 131.52060468056197, 188.5753101780003, 28.303083325889997, 22.91691676265927, 131.47939880676725},
		{12.38659860174893, -123.00182001175405, 351.0869827898082, 79.45362337387198, 13.092489873895186, -123.1154794135781},
		{-61.85524721381324, 34.13109516590251, 249.12885144712033, 301.52268100656, -62.713070661498655, 28.597790472373454},
	}
	for _, c := range cases {
		lat, lng := destination(c.Lat1, c.Lng1, c.Distance, c.Bearing)
		assert.InDelta(t, c.Lat2, lat, 0.00001, "latitude error")
		assert.InDelta(t, c.Lng2, lng, 0.00001, "longitude error")
	}
}
