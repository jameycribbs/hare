package disk

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestAllTableFileTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//readRec...

			filePtr, err := os.OpenFile("./testdata/episodes.json", os.O_RDWR, 0660)
			if err != nil {
				t.Fatal(err)
			}

			tf, err := NewTableFile("episodes", filePtr)
			if err != nil {
				t.Fatal(err)
			}

			rec, err := tf.readRec(3)
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
		tableFileTestSetup()
		t.Run(strconv.Itoa(i), fn)
		tableFileTestTeardown()
	}
}

func tableFileTestSetup() {
	cmd := exec.Command("cp", "./testdata/episodes_default.txt", "./testdata/episodes.json")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func tableFileTestTeardown() {
	if err := os.Remove("./testdata/episodes.json"); err != nil {
		panic(err)
	}
}
