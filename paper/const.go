package paper

type PageStyle string

var (
	A3                PageStyle = "A3"
	A3Landscape       PageStyle = "A3.landscape"
	A4                PageStyle = "A4"
	A4Landscape       PageStyle = "A4.landscape"
	L3One             PageStyle = "L3-1 "
	L3Two             PageStyle = "L3-2 "
	L3Thr             PageStyle = "L3-3 "
	A5                PageStyle = "A5"
	A5Landscape       PageStyle = "A5.landscape"
	Letter            PageStyle = "letter"
	LetterLandscape   PageStyle = "letter.landscape "
	Legal             PageStyle = "legal"
	LegalLandscape    PageStyle = "legal.landscape"
	Document          PageStyle = "document"
	DocumentLandscape PageStyle = "document.landscape"
	Ratio                       = 0.0393700787
)

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
