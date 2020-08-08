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

			rec, err := dsk.ReadRec("episodes", 3)
			if err != nil {
				t.Fatal(err)
			}

			recStr := string(rec)

			recShouldBe := "{\"id\":3,\"season\":6,\"episode\":9,\"film\":\"The Skydivers\",\"shorts\":[\"Why Study Industrial Arts?\"], \"year_film_released\":1963,\"date_episode_aired\":\"1994-08-27T00:00:00Z\",\"host\":\"Mike\"}\n"

			if recShouldBe != recStr {
				t.Errorf("want %v; got %v", recShouldBe, recStr)
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
	cmd := exec.Command("cp", "./testdata/episodes_default.txt", "./testdata/episodes.json")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func diskTestTeardown() {
	if err := os.Remove("./testdata/episodes.json"); err != nil {
		panic(err)
	}
}
