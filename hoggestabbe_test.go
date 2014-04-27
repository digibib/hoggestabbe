package main

import (
	"encoding/xml"
	"strings"
	"testing"
)

var marcInput string = `*0010010463
*008871001                a          0 nob r
*015  $a29$bBibliofilID
*019  $bl
*0823 $a242
*090  $c242$dKar
*100 0$aKarlén, Barbro$d1954-$jsv.$310008600
*24010$aI begynnelsen skapade Gud
*24510$aI begynnelsen skapte Gud$bdikt og salmer i norsk gjendiktning   ved Inger Hagerup ; prosaen overs. av Gunnel Malmström ; illus. av Vanja Hübinette-Simonson
*260  $aOslo$bAschehoug$c1968
*300  $a97 s.$bill.
*574  $aOriginaltittel: I begynnelsen skapade Gud
*655  $aAndaktsbøker$310008700
*70010$aHagerup, Inger$d1905-1985$jn.$310008800
*850  $aDEICHM$sn`

var xmlOutput string = `<record>
  <leader>00573nam a22     1  4500</leader>
  <controlfield tag="001">0010463</controlfield>
  <controlfield tag="008">871001                a          0 nob r</controlfield>
  <datafield tag="015" ind1=" " ind2=" ">
    <subfield code="a">29</subfield>
    <subfield code="b">BibliofilID</subfield>
  </datafield>
  <datafield tag="019" ind1=" " ind2=" ">
    <subfield code="b">l</subfield>
  </datafield>
  <datafield tag="082" ind1="3" ind2=" ">
    <subfield code="a">242</subfield>
  </datafield>
  <datafield tag="090" ind1=" " ind2=" ">
    <subfield code="c">242</subfield>
    <subfield code="d">Kar</subfield>
  </datafield>
  <datafield tag="100" ind1=" " ind2="0">
    <subfield code="a">Karlén, Barbro</subfield>
    <subfield code="d">1954-</subfield>
    <subfield code="j">sv.</subfield>
    <subfield code="3">10008600</subfield>
  </datafield>
  <datafield tag="240" ind1="1" ind2="0">
    <subfield code="a">I begynnelsen skapade Gud</subfield>
  </datafield>
  <datafield tag="245" ind1="1" ind2="0">
    <subfield code="a">I begynnelsen skapte Gud</subfield>
    <subfield code="b">dikt og salmer i norsk gjendiktning   ved Inger Hagerup ; prosaen overs. av Gunnel Malmström ; illus. av Vanja Hübinette-Simonson</subfield>
  </datafield>
  <datafield tag="260" ind1=" " ind2=" ">
    <subfield code="a">Oslo</subfield>
    <subfield code="b">Aschehoug</subfield>
    <subfield code="c">1968</subfield>
  </datafield>
  <datafield tag="300" ind1=" " ind2=" ">
    <subfield code="a">97 s.</subfield>
    <subfield code="b">ill.</subfield>
  </datafield>
  <datafield tag="574" ind1=" " ind2=" ">
    <subfield code="a">Originaltittel: I begynnelsen skapade Gud</subfield>
  </datafield>
  <datafield tag="655" ind1=" " ind2=" ">
    <subfield code="a">Andaktsbøker</subfield>
    <subfield code="3">10008700</subfield>
  </datafield>
  <datafield tag="700" ind1="1" ind2="0">
    <subfield code="a">Hagerup, Inger</subfield>
    <subfield code="d">1905-1985</subfield>
    <subfield code="j">n.</subfield>
    <subfield code="3">10008800</subfield>
  </datafield>
  <datafield tag="850" ind1=" " ind2=" ">
    <subfield code="a">DEICHM</subfield>
    <subfield code="s">n</subfield>
  </datafield>
</record>`

func TestConversion(t *testing.T) {
	record, err := Parse(strings.Split(marcInput, "\n"))
	if err != nil {
		t.Fatal(err)
	}

	output, err := xml.MarshalIndent(record, "", "  ")
	if string(output) != xmlOutput {
		t.Errorf("wanted:\n%v, \nbut got:\n%v", xmlOutput, string(output))
	}
}
