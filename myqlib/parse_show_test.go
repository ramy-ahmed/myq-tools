package myqlib

import (
	"testing"
	"time"
	// "fmt"
)

func TestSingleSample(t *testing.T) {
	l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysqladmin.single", ""}
	samples, err := l.getStatus()
	if err != nil {
		t.Error(err)
	}

	// Check some types on some known metrics to verify autodetection
	sample := <-samples
	typeTests := map[string]string{
		"connections":                "int64",
		"compression":                "string",
		"wsrep_local_send_queue_avg": "float64",
		"binlog_snapshot_file":       "string",
	}

	for varname, expectedtype := range typeTests {
		i, ierr := sample.getInt(varname)
		if ierr == nil {
			if expectedtype != "int64" {
				t.Fatal("Found integer, expected", expectedtype, "for", varname, "value: `", i, "`")
			} else {
				continue
			}
		}

		f, ferr := sample.getFloat(varname)
		if ferr == nil {
			if expectedtype != "float64" {
				t.Fatal("Found float, expected", expectedtype, "for", varname, "value: `", f, "`")
			} else {
				continue
			}
		}

		s, serr := sample.getString(varname)
		if serr == nil {
			if expectedtype != "string" {
				t.Fatal("Found string, expected", expectedtype, "for", varname, "value: `", s, "`")
			} else {
				continue
			}
		}
	}
}

func TestTwoSamples(t *testing.T) {
	l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysqladmin.two", ""}
	samples, err := l.getStatus()

	if err != nil {
		t.Error(err)
	}

	checksamples(t, samples, 2)
}

func TestManySamples(t *testing.T) {
	if testing.Short() {
		return
	}

	l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysqladmin.lots", ""}
	samples, err := l.getStatus()

	if err != nil {
		t.Error(err)
	}

	checksamples(t, samples, 220)
}

func TestSingleBatchSample(t *testing.T) {
	l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysql.single", ""}
	samples, err := l.getStatus()
	if err != nil {
		t.Error(err)
	}

	// Check some types on some known metrics to verify autodetection
	sample := <-samples
	typeTests := map[string]string{
		"connections":                "int64",
		"compression":                "string",
		"wsrep_local_send_queue_avg": "float64",
		"binlog_snapshot_file":       "string",
	}

	for varname, expectedtype := range typeTests {
		i, ierr := sample.getInt(varname)
		if ierr == nil {
			if expectedtype != "int64" {
				t.Fatal("Found integer, expected", expectedtype, "for", varname, "value: `", i, "`")
			} else {
				continue
			}
		}

		f, ferr := sample.getFloat(varname)
		if ferr == nil {
			if expectedtype != "float64" {
				t.Fatal("Found float, expected", expectedtype, "for", varname, "value: `", f, "`")
			} else {
				continue
			}
		}

		s, serr := sample.getString(varname)
		if serr == nil {
			if expectedtype != "string" {
				t.Fatal("Found string, expected", expectedtype, "for", varname, "value: `", s, "`")
			} else {
				continue
			}
		}
	}
}

func TestTwoBatchSamples(t *testing.T) {
	l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysql.two", ""}
	samples, err := l.getStatus()

	if err != nil {
		t.Error(err)
	}

	checksamples(t, samples, 2)
}

func TestManyBatchSamples(t *testing.T) {
	if testing.Short() {
		return
	}

	l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysql.lots", ""}
	samples, err := l.getStatus()

	if err != nil {
		t.Error(err)
	}

	checksamples(t, samples, 215)
}

func checksamples(t *testing.T, samples chan MyqSample, expected int) {
	i := 0
	var prev MyqSample
	for sample := range samples {
		t.Log("New MyqSample", i, len(sample), sample["uptime"])
		if prev != nil {
			t.Log("\tPrev", i, len(prev), prev["uptime"])

			if prev["uptime"] == sample["uptime"] {
				t.Fatal("previous has same uptime")
			}
		}

		if len(prev) > 0 && len(prev) > len(sample) {
			t.Log(prev["uptime"], "(previous) had", len(prev), "keys.  Current current has", len(sample))
			for pkey := range prev {
				_, ok := (sample)[pkey]
				if !ok {
					t.Log("Missing", pkey, "from current sample")
				}
			}
			t.Fatal("")
		}
		prev = sample
		i++
	}

	if i != expected {
		t.Errorf("Got unexpected number of samples: %d", i)
	}
}

func TestTokuSample(t *testing.T) {
	l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysql.toku", ""}
	samples, err := l.getStatus()

	if err != nil {
		t.Error(err)
	}

	checksamples(t, samples, 2)
}


func BenchmarkParseStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysqladmin.single", ""}
		samples, err := l.getStatus()

		if err != nil {
			b.Error(err)
		}
		<-samples
	}
}

func BenchmarkParseStatusBatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysql.single", ""}
		samples, err := l.getStatus()

		if err != nil {
			b.Error(err)
		}
		<-samples
	}
}

func BenchmarkParseVariablesBatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := FileLoader{loaderInterval(1 * time.Second), "../testdata/variables", ""}
		samples, err := l.getStatus()

		if err != nil {
			b.Error(err)
		}
		<-samples
	}
}

func BenchmarkParseVariablesTabular(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := FileLoader{loaderInterval(1 * time.Second), "../testdata/variables.tab", ""}
		samples, err := l.getStatus()

		if err != nil {
			b.Error(err)
		}
		<-samples
	}
}

func BenchmarkParseManyBatchSamples(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysql.lots", ""}
		samples, err := l.getStatus()

		if err != nil {
			b.Error(err)
		}
		for j := 0; j <= 61; j++ {
			<-samples
		}
	}
}

func BenchmarkParseManySamples(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := FileLoader{loaderInterval(1 * time.Second), "../testdata/mysqladmin.lots", ""}
		samples, err := l.getStatus()

		if err != nil {
			b.Error(err)
		}
		for j := 0; j <= 220; j++ {
			<-samples
		}
	}
}

func BenchmarkParseManySamplesLongInterval(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := FileLoader{loaderInterval(1 * time.Minute), "../testdata/mysqladmin.lots", ""}
		samples, err := l.getStatus()

		if err != nil {
			b.Error(err)
		}
		for j := 0; j <= 220; j++ {
			<-samples
		}
	}
}
