package disk

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestAllDiskTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//ReadRec...

			dsk, err := New("./testdata", ".json")
			if err != nil {
				t.Fatal(err)
			}

			rec, err := dsk.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := "{\"id\":3,\"first_name\":\"Bill\",\"last_name\":\"Shakespeare\",\"age\":18}\n"
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
	}

	for i, fn := range tests {
		diskTestSetup()
		t.Run(strconv.Itoa(i), fn)
		diskTestTeardown()
	}
}

func diskTestSetup() {
	cmd := exec.Command("cp", "./testdata/contacts.bak", "./testdata/contacts.json")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func diskTestTeardown() {
	if err := os.Remove("./testdata/contacts.json"); err != nil {
		panic(err)
	}
}
