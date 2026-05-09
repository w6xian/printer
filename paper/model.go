package paper

type Paper struct {
	Width  float64
	Height float64
}

func GetPageStyle(name string) PageStyle {
	switch name {
	case "A3":
		return A3
	case "A3.landscape":
		return A3Landscape
	case "A4":
		return A4
	case "A4.landscape":
		return A4Landscape
	case "L3-1":
		return L3One
	case "L3-2":
		return L3Two
	case "L3-3":
		return L3Thr
	case "A5":
		return A5
	case "A5.landscape":
		return A5Landscape
	case "letter":
		return Letter
	case "letter.landscape":
		return LetterLandscape
	case "legal":
		return Legal
	case "legal.landscape":
		return LegalLandscape
	case "document":
		return Document
	case "document.landscape":
		return DocumentLandscape
	default:
		return A4
	}
}

/**
1mm == 0.0393700787
A3                  { width: 297mm; height: 419mm }
A3.landscape        { width: 420mm; height: 296mm }
A4                  { width: 210mm; height: 296mm }
A4.landscape        { width: 297mm; height: 209mm }
L3-1                { width: 214mm; height: 93mm }
L3-2                { width: 214mm; height: 140mm }
L3-3                { width: 214mm; height: 280mm }
A5                  { width: 148mm; height: 209mm }
A5.landscape        { width: 210mm; height: 147mm }
letter              { width: 216mm; height: 279mm }
letter.landscape    { width: 280mm; height: 215mm }
legal               { width: 216mm; height: 356mm }
legal.landscape     { width: 357mm; height: 215mm }
document            { width: 120mm; height: 240mm }
document.landscape  { width: 240mm; height: 120mm }
**/

func (p *Paper) FixWithPageStyle(name PageStyle) *Paper {
	switch name {
	case A3:
		p.Width = Ratio * 297
		p.Height = Ratio * 419
	case A3Landscape:
		p.Width = Ratio * 420
		p.Height = Ratio * 296
	case A4:
		p.Width = Ratio * 210
		p.Height = Ratio * 296
	case A4Landscape:
		p.Width = Ratio * 297
		p.Height = Ratio * 209
	case L3One:
		p.Width = Ratio * 214
		p.Height = Ratio * 93
	case L3Two:
		p.Width = Ratio * 214
		p.Height = Ratio * 140
	case L3Thr:
		p.Width = Ratio * 214
		p.Height = Ratio * 280
	case A5:
		p.Width = Ratio * 148
		p.Height = Ratio * 209
	case A5Landscape:
		p.Width = Ratio * 210
		p.Height = Ratio * 147
	case Letter:
		p.Width = Ratio * 216
		p.Height = Ratio * 279
	case LetterLandscape:
		p.Width = Ratio * 280
		p.Height = Ratio * 215
	case Legal:
		p.Width = Ratio * 216
		p.Height = Ratio * 356
	case LegalLandscape:
		p.Width = Ratio * 357
		p.Height = Ratio * 215
	case Document:
		p.Width = Ratio * 120
		p.Height = Ratio * 240
	case DocumentLandscape:
		p.Width = Ratio * 240
		p.Height = Ratio * 120
	}
	return p
}

func NewPager() *Paper {
	// default 与 pdf 一致
	return &Paper{
		Width:  8.5,
		Height: 11,
	}
}

func GetInchFromMM(mm float64) float64 {
	return Ratio * mm
}

type Padding struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

func NewPadding(p1 float64, pn ...float64) *Padding {
	p := &Padding{Left: 0, Top: 0, Right: 0, Bottom: 0}
	l := len(pn)
	if l == 0 {
		p.Top = p1
		p.Bottom = p1
		p.Left = p1
		p.Right = p1
	} else if l == 1 {
		p.Top = p1
		p.Bottom = p1
		p.Left = pn[0]
		p.Right = pn[0]
	} else if l == 2 {
		p.Top = p1
		p.Left = pn[0]
		p.Right = pn[1]
		p.Bottom = pn[0]
	} else {
		p.Top = p1
		p.Right = pn[0]
		p.Bottom = pn[1]
		p.Left = pn[2]
	}
	return p
}
