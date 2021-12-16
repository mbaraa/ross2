package postsgen

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"strings"

	"github.com/mbaraa/dsc_logo_generator/logogen"
	"github.com/mbaraa/ross2/models"
	"github.com/ungerik/go-cairo"
)

/*
	This is where the fun continues!
*/

type Point2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type TextFieldProps struct {
	Position Point2  `json:"startPosition"`
	XWidth   float64 `json:"width"`
	FontSize float64 `json:"fontSize"`
}

type FieldsProps struct {
	TeamNameProps     TextFieldProps
	TeamOrderProps    TextFieldProps
	MembersNamesProps []TextFieldProps
}

type Builder struct {
	teams          []models.Team
	teamNameProps  TextFieldProps
	membersProps   []TextFieldProps
	teamOrderProps TextFieldProps
	b64Image       string
	baseImage      image.Image
}

func NewTeamsPostsGeneratorBuilder() *Builder { return new(Builder) }

func (b *Builder) Teams(t []models.Team) *Builder {
	b.teams = t[:]
	return b
}

func (b *Builder) TeamNameProps(props TextFieldProps) *Builder {
	b.teamNameProps = props
	return b
}

func (b *Builder) MembersNamesProps(props []TextFieldProps) *Builder {
	b.membersProps = props[:]
	return b
}

func (b *Builder) TeamOrderProps(props TextFieldProps) *Builder {
	b.teamOrderProps = props
	return b
}

func (b *Builder) B64Image(img string) *Builder {
	b.b64Image = img
	return b
}

func (b *Builder) setBaseImage() error {
	img, err := base64.StdEncoding.DecodeString(b.b64Image)
	if err != nil {
		return err
	}
	r := bytes.NewReader(img)
	b.baseImage, err = png.Decode(r)

	return err
}

func (b *Builder) shiftPoints() {
	b.teamNameProps.Position = movePoint(b.teamNameProps.Position, Point2{0, .08})
	for i := range b.membersProps {
		b.membersProps[i].Position = movePoint(b.membersProps[i].Position, Point2{0, .04})
	}
}

func (b *Builder) verify() error {
	errSb := new(strings.Builder)
	if b.teams == nil || len(b.teams) == 0 {
		errSb.WriteString("Posts Generator Builder: missing teams!\n")
	}
	if b.teamNameProps == (TextFieldProps{}) {
		errSb.WriteString("Posts Generator Builder: missing team name props!\n")
	}
	if b.teamOrderProps == (TextFieldProps{}) {
		errSb.WriteString("Posts Generator Builder: missing team order props!\n")
	}
	if b.membersProps == nil || len(b.membersProps) == 0 {
		errSb.WriteString("Posts Generator Builder: missing members props!\n")
	}
	if b.b64Image == "" {
		errSb.WriteString("Posts Generator Builder: missing b64 base image!\n")
	}

	if errSb.Len() > 0 {
		fmt.Println(errSb.String())
		return errors.New(errSb.String())
	}

	return nil
}

func (b *Builder) GetTeamsPostsGenerator() (*TeamPostsGenerator, error) {
	err := b.verify()
	if err != nil {
		return nil, err
	}

	err = b.setBaseImage()
	if err != nil {
		return nil, err
	}
	b.shiftPoints()

	return NewTeamsPostsGenerator(b), nil
}

/****************************************************************/

type TeamPostsGenerator struct {
	teams        []models.Team
	fieldsProps  FieldsProps
	baseImage    image.Image
	workingImage *cairo.Surface
}

func NewTeamsPostsGenerator(b *Builder) *TeamPostsGenerator {
	return &TeamPostsGenerator{
		teams: b.teams,
		fieldsProps: FieldsProps{
			TeamNameProps:     b.teamNameProps,
			TeamOrderProps:    b.teamOrderProps,
			MembersNamesProps: b.membersProps,
		},
		baseImage: b.baseImage,
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
		return
	}

	_, err = zipFile.Seek(0, 0)
	if err != nil {
		return
	}

	b, err := io.ReadAll(zipFile)
	_ = zipFile.Close()
	_ = os.Remove(zipFile.Name())
	zipBytes = []byte(base64.StdEncoding.EncodeToString(b))

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
	for i, team := range p.teams {
		post, err := p.generatePost(team, i+1)
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
func (p *TeamPostsGenerator) generatePost(team models.Team, teamOrder int) (string, error) {
	p.workingImage = cairo.NewSurfaceFromImage(p.baseImage)

	p.showTeamOrder(teamOrder)
	p.showTeamName(team)
	p.showMembers(team.Members)

	return p.exportImageToB64()
}

func (p *TeamPostsGenerator) showTeamOrder(teamOrder int) {
	p.workingImage.SelectFontFace("Product Sans", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	t := NewText(fmt.Sprintf("Team %d", teamOrder), color.RGBA64{R: 255, G: 255, B: 255, A: 255}, p.fieldsProps.TeamOrderProps.FontSize, "Product Sans")
	p.workingImage.SetSourceRGBA(t.GetColorRGBA())

	xLength, newSize := t.GetXLengthUsingParent(p.fieldsProps.TeamOrderProps.XWidth, .85)
	p.workingImage.SetFontSize(newSize)

	p.workingImage.MoveTo(p.fieldsProps.TeamOrderProps.Position.X+
		(getCenterStartOfElement(xLength*1.1, p.fieldsProps.TeamOrderProps.XWidth)),
		p.fieldsProps.TeamOrderProps.Position.Y)
	p.workingImage.ShowText(t.GetContent())
}

func (p *TeamPostsGenerator) showTeamName(team models.Team) {
	p.workingImage.SelectFontFace("Product Sans", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	t := NewText(team.Name, color.RGBA64{R: 255, G: 255, B: 255, A: 255}, p.fieldsProps.TeamNameProps.FontSize, "Product Sans")
	t.SetFontSize(p.fieldsProps.TeamNameProps.FontSize)

	xLength, newSize := t.GetXLengthUsingParent(p.fieldsProps.TeamNameProps.XWidth, .8)
	p.workingImage.SetFontSize(newSize)

	p.workingImage.MoveTo(p.fieldsProps.TeamNameProps.Position.X+
		getCenterStartOfElement(xLength*1.1, p.fieldsProps.TeamNameProps.XWidth),
		p.fieldsProps.TeamNameProps.Position.Y)
	p.workingImage.ShowText(t.GetContent())
}

func (p *TeamPostsGenerator) showMembers(members []models.Contestant) {
	t := NewText("", color.RGBA64{R: 255, G: 255, B: 255, A: 255}, 0, "Product Sans")
	p.workingImage.SelectFontFace("Product Sans", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)

	for propsIndex, membersIndex, membersCount, propsCount := 0, 0, len(members), len(p.fieldsProps.MembersNamesProps); membersIndex < membersCount; propsIndex, membersIndex = (propsIndex+1)%propsCount, membersIndex+1 {

		t.SetContent(members[membersIndex].Name)
		t.SetFontSize(p.fieldsProps.MembersNamesProps[propsIndex].FontSize)

		xLength, newSize := t.GetXLengthUsingParent(p.fieldsProps.MembersNamesProps[propsIndex].XWidth, .8)
		p.workingImage.SetFontSize(newSize)

		p.workingImage.MoveTo(p.fieldsProps.MembersNamesProps[propsIndex].Position.X+
			getCenterStartOfElement(xLength*1.1, p.fieldsProps.MembersNamesProps[propsIndex].XWidth),
			p.fieldsProps.MembersNamesProps[propsIndex].Position.Y)
		p.workingImage.ShowText(t.GetContent())
	}
}

func (p *TeamPostsGenerator) exportImageToB64() (string, error) {
	stream, err := p.workingImage.WriteToPNGStream()
	if err != cairo.STATUS_SUCCESS {
		return "", errors.New(err.String())
	}

	p.workingImage.Finish()

	return base64.StdEncoding.EncodeToString(stream), nil
}

func movePoint(p, percentage Point2) Point2 {
	p.X += p.X * percentage.X
	p.Y += p.Y * percentage.Y
	return p
}

/****************************************************************/
// sample post

var (
	ThreeMembersPostSamplePostBuilder = NewTeamsPostsGeneratorBuilder().
		TeamOrderProps(TextFieldProps{
			Position: Point2{X: 165, Y: 82},
			XWidth:   233,
			FontSize: 40,
		}).
		TeamNameProps(TextFieldProps{
			Position: Point2{X: 111, Y: 140},
			XWidth:   369,
			FontSize: 30.5,
		}).
		MembersNamesProps([]TextFieldProps{
			{
				Position: Point2{X: 73, Y: 210},
				XWidth:   447,
				FontSize: 25,
			},
			{
				Position: Point2{X: 73, Y: 272},
				XWidth:   447,
				FontSize: 25,
			},
			{
				Position: Point2{X: 73, Y: 333},
				XWidth:   447,
				FontSize: 25,
			},
		}).
		B64Image("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ve7z7u72uf17pzm/H2RDcqF/v7tYAAAPw2lUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4KPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgNC40LjAtRXhpdjIiPgogPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgeG1sbnM6eG1wTU09Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9tbS8iCiAgICB4bWxuczpzdEV2dD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL3NUeXBlL1Jlc291cmNlRXZlbnQjIgogICAgeG1sbnM6ZGM9Imh0dHA6Ly9wdXJsLm9yZy9kYy9lbGVtZW50cy8xLjEvIgogICAgeG1sbnM6R0lNUD0iaHR0cDovL3d3dy5naW1wLm9yZy94bXAvIgogICAgeG1sbnM6dGlmZj0iaHR0cDovL25zLmFkb2JlLmNvbS90aWZmLzEuMC8iCiAgICB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iCiAgIHhtcE1NOkRvY3VtZW50SUQ9ImdpbXA6ZG9jaWQ6Z2ltcDpjY2M5ZDlhYS0yYWVhLTRjYmYtODgxNy1mMmNiODVlMzcxZjQiCiAgIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6Nzc0YThiZDgtNjMyZi00NWFiLWE4MjktNzk4NGRmMjcxY2JlIgogICB4bXBNTTpPcmlnaW5hbERvY3VtZW50SUQ9InhtcC5kaWQ6ZmUxOTVlNDYtNzUyMS00ZWIzLWIwNGMtMjY4Njg2YmFjNTZkIgogICBkYzpGb3JtYXQ9ImltYWdlL3BuZyIKICAgR0lNUDpBUEk9IjIuMCIKICAgR0lNUDpQbGF0Zm9ybT0iTGludXgiCiAgIEdJTVA6VGltZVN0YW1wPSIxNjM5NjU4Njc2MTg3Nzg2IgogICBHSU1QOlZlcnNpb249IjIuMTAuMjgiCiAgIHRpZmY6T3JpZW50YXRpb249IjEiCiAgIHhtcDpDcmVhdG9yVG9vbD0iR0lNUCAyLjEwIj4KICAgPHhtcE1NOkhpc3Rvcnk+CiAgICA8cmRmOlNlcT4KICAgICA8cmRmOmxpCiAgICAgIHN0RXZ0OmFjdGlvbj0ic2F2ZWQiCiAgICAgIHN0RXZ0OmNoYW5nZWQ9Ii8iCiAgICAgIHN0RXZ0Omluc3RhbmNlSUQ9InhtcC5paWQ6MWRkMjZjMjktZDQyZC00YWE1LWJhMDgtMzgxNDUzNGIyYmJmIgogICAgICBzdEV2dDpzb2Z0d2FyZUFnZW50PSJHaW1wIDIuMTAgKExpbnV4KSIKICAgICAgc3RFdnQ6d2hlbj0iMjAyMS0xMi0wMlQwNDoyMTowMSswMjowMCIvPgogICAgIDxyZGY6bGkKICAgICAgc3RFdnQ6YWN0aW9uPSJzYXZlZCIKICAgICAgc3RFdnQ6Y2hhbmdlZD0iLyIKICAgICAgc3RFdnQ6aW5zdGFuY2VJRD0ieG1wLmlpZDo4ZTE5ZGIyOC0yY2EyLTQ4YzUtODVmZi1lMzhhNjhhNzZhYjIiCiAgICAgIHN0RXZ0OnNvZnR3YXJlQWdlbnQ9IkdpbXAgMi4xMCAoTGludXgpIgogICAgICBzdEV2dDp3aGVuPSIyMDIxLTEyLTAyVDA2OjQ3OjAxKzAyOjAwIi8+CiAgICAgPHJkZjpsaQogICAgICBzdEV2dDphY3Rpb249InNhdmVkIgogICAgICBzdEV2dDpjaGFuZ2VkPSIvIgogICAgICBzdEV2dDppbnN0YW5jZUlEPSJ4bXAuaWlkOjE2NWM5NTNkLWFiMWItNGJiMy1iYWNhLWIzNWFjOTM0NmNhZSIKICAgICAgc3RFdnQ6c29mdHdhcmVBZ2VudD0iR2ltcCAyLjEwIChMaW51eCkiCiAgICAgIHN0RXZ0OndoZW49IjIwMjEtMTItMTZUMTM6MTg6MzgrMDI6MDAiLz4KICAgICA8cmRmOmxpCiAgICAgIHN0RXZ0OmFjdGlvbj0ic2F2ZWQiCiAgICAgIHN0RXZ0OmNoYW5nZWQ9Ii8iCiAgICAgIHN0RXZ0Omluc3RhbmNlSUQ9InhtcC5paWQ6NTc3YmFmMmEtMGE2Mi00ZTMzLTg5OGItNWQ1ZDY2NWZmYmNjIgogICAgICBzdEV2dDpzb2Z0d2FyZUFnZW50PSJHaW1wIDIuMTAgKExpbnV4KSIKICAgICAgc3RFdnQ6d2hlbj0iMjAyMS0xMi0xNlQxNDo0NDozNiswMjowMCIvPgogICAgPC9yZGY6U2VxPgogICA8L3htcE1NOkhpc3Rvcnk+CiAgPC9yZGY6RGVzY3JpcHRpb24+CiA8L3JkZjpSREY+CjwveDp4bXBtZXRhPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgCjw/eHBhY2tldCBlbmQ9InciPz6JyZPKAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH5QwQDCwkE5zsFAAAEZZJREFUeNrt3W9zHMWBwOHu+buSrIBTiWRDwqWKIhVexHcUVL7/J7gUVCqV+DgqxI5tsLEOYku7M9Pd92KwI4xxbLBs7eh5XqjWYlntzPbOb3p2djdeu/aHAAAXQGUVAKB5AKB5AKB5AKB5AKB5AKB5AKB5AKB5AKB5AGgeAGgeAGgeAGgeAGgeAGgeAGgeAGgeAGgeAGgeAJoHAJoHAJoHAJoHAJoHAJoHAJoHAJoHAJoHAJoHgOYBgOYBgOYBgOYBgOYBgOYBgOYBgOYBgOYBgOYBoHkAoHkAoHkAoHkAoHkAoHkAoHkAoHkAoHkAoHkAaB4AaB4AaB4AaB4AaB4AaB4AaB4AaB4AaB4AaB4AmgcAmgcAmgcAmgcAmgcAmgcAmgcAmgcAmgcAmgeA5gGA5gGA5gGA5gGA5gGA5gGA5gGA5gGA5gGA5gGgeQCgeQCgeQCgeQCgeQCgeQCgeQCgeQCgeQCgeQBoHgBoHgBoHgBo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ld59992jo6M7d+60bbu7u7ter6vKOTsA35pneKWUvb29lNKtW7e+/vrrn//855q3raZp+vvf/75er3d2dnzGNMBpdV1vNpu5effv379+/Xop5cqVK0te5MPDtxe8eAcHB3fu3Jmm6dKlS1VVVVXldBKAWSl5ngwcHx9/+umnX3311fvvv7/sRV7+sb6c82effXbr1q0QQtM0RjnAbBiGvb29nPPNmzdv3Lixv7+/+EWO1679YfELefv27bt3vzw8vHLt2rW+79u2ffyybc65ruumac74bQwAr01K6fHXiFZVNf8MIezu7t64ceOvf/3r/fv333zzzatXry5+VVyI5s3+9re/bTab995779KlS2+++ebe3t5cvpTS/AY+Twxgkeq6nr89O+c8TdN6vd5sNuM4/vnPf/7mm2+qqvrtb397QVbFBWrenL1vvvk6xmq1Wu3s7PR9P78lZX4zgycGsEjDMDw6oSGklMZxHIZhmqaU0v7+/ttvv31xVsXFen3rN7/5zXzhiy++OD4+fvjwYUop52yeByx7nhceHdisqqppmq7rVqvVRTiYeaGb99jh4aGnAcBF4z3aAGgeAGgeAGgeAGgeAGgeAGgeAGgeAGgeAGgeAJoHAJoHAJoHAJoHAJoHAJoHAJoHAJoHAJoHAJoHgOYBgOYBgOYBgOYBgOYBgOYBgOYBgOYBgOYBgOYBcJH9P8D6ix/p3RrMAAAAAElFTkSuQmCC")
)
