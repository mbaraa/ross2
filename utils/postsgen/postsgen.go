package postsgen

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/mbaraa/dsc_logo_generator/logogen"
	"github.com/mbaraa/ross2/models"
	"github.com/ungerik/go-cairo"
)

/*
	This is where the fun continues!
*/

type Point2 struct {
	X, Y float64
}

type TextFieldProps struct {
	Position Point2
	XWidth   float64
	FontSize float64
}

type TeamPostsGenerator struct {
	teams         []models.Team
	teamNameProps TextFieldProps
	membersProps  []TextFieldProps
	baseImage     image.Image
}

func NewTeamsPostsGenerator(teams []models.Team, teamNameProps TextFieldProps, membersNamesProps []TextFieldProps,
	baseImage string) *TeamPostsGenerator {

	img, err := base64.StdEncoding.DecodeString(baseImage)
	if err != nil {
		return nil
	}
	r := bytes.NewReader(img)
	pngImage, err := png.Decode(r)
	if err != nil {
		return nil
	}

	teamNameProps.Position = movePoint(teamNameProps.Position, Point2{0, .08})
	for i := range membersNamesProps {
		membersNamesProps[i].Position = movePoint(membersNamesProps[i].Position, Point2{0, .04})
	}

	return &TeamPostsGenerator{
		teams:         teams,
		teamNameProps: teamNameProps,
		membersProps:  membersNamesProps,
		baseImage:     pngImage,
	}
}

func (p *TeamPostsGenerator) GenerateToB64Images() ([]string, error) {
	return p.generatePosts()
}

func (p *TeamPostsGenerator) GenerateToZipFile() (*os.File, error) {
	return p.generateZipFile()
}

func (p *TeamPostsGenerator) GenerateToZipFileBytes() (zipBytes []byte, err error) {
	zipFile, err := p.generateZipFile()
	if err != nil {
		return nil, err
	}

	_, err = zipFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	_, err = base64.
		NewEncoder(base64.StdEncoding, zipFile).
		Write(zipBytes)

	_ = zipFile.Close()

	return
}

func (p *TeamPostsGenerator) generateZipFile() (*os.File, error) {
	posts, err := p.generatePosts()
	if err != nil {
		return nil, err
	}

	pz, err := NewPostsZipper(posts)
	if err != nil {
		return nil, err
	}

	return pz.MakeZipFile()
}

func (p *TeamPostsGenerator) generatePosts() (postsInB64 []string, err error) {
	for _, team := range p.teams {
		post, err := p.generatePost(team)
		if err != nil {
			return nil, err
		}
		postsInB64 = append(postsInB64, post)
	}
	return
}

// NewText returns a new logogen.Text instance w/o an error,
// ie when a font error is encountered the Default font is used
func NewText(content string, fgColor color.RGBA64, fontSize float64, fontName string) logogen.Text {
	font, err := os.ReadFile(fmt.Sprintf("./res/%s.ttf", fontName))
	if err != nil {
		font, _ = os.ReadFile("./res/Default.ttf")
	}
	txt0, _ := logogen.NewText(content, fgColor, fontSize, font)
	return *txt0
}

// getCenterStartOfElement return the coordinate of the first point of the child element
// that will allow it to appear in the middle of the parent
func getCenterStartOfElement(childLength float64, parentLength float64) float64 {
	return math.Abs((parentLength - childLength) / 2.0)
}

// don't ask about the ratios, they just work :)
func (p *TeamPostsGenerator) generatePost(team models.Team) (string, error) {
	sur := cairo.NewSurfaceFromImage(p.baseImage)
	defer sur.Finish()

	sur.SelectFontFace("Product Sans", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)

	t := NewText(team.Name, color.RGBA64{R: 255, G: 255, B: 255, A: 255}, p.teamNameProps.FontSize, "Product Sans")
	sur.SetSourceRGBA(t.GetColorRGBA())

	xLength, newSize := t.GetXLengthUsingParent(p.teamNameProps.XWidth, .85)
	sur.SetFontSize(newSize)

	sur.MoveTo(p.teamNameProps.Position.X+
		(getCenterStartOfElement(xLength*1.1, p.teamNameProps.XWidth)),
		p.teamNameProps.Position.Y)
	sur.ShowText(t.GetContent())

	sur.SelectFontFace("Product Sans", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)

	for propsIndex, membersIndex, membersCount, propsCount := 0, 0, len(team.Members), len(p.membersProps); membersIndex < membersCount; propsIndex, membersIndex = (propsIndex+1)%propsCount, membersIndex+1 {
		t.SetContent(team.Members[membersIndex].Name)
		t.SetFontSize(p.membersProps[propsIndex].FontSize)

		xLength, newSize = t.GetXLengthUsingParent(p.membersProps[propsIndex].XWidth, .8)
		sur.SetFontSize(newSize)

		sur.MoveTo(p.membersProps[propsIndex].Position.X+
			getCenterStartOfElement(xLength*1.1, p.membersProps[propsIndex].XWidth),
			p.membersProps[propsIndex].Position.Y)
		sur.ShowText(t.GetContent())
	}

	stream, err := sur.WriteToPNGStream()
	if err != cairo.STATUS_SUCCESS {
		return "", errors.New(err.String())
	}

	return base64.StdEncoding.EncodeToString(stream), nil
}

func movePoint(p, percentage Point2) Point2 {
	p.X += p.X * percentage.X
	p.Y += p.Y * percentage.Y
	return p
}
