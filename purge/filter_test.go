package purge

import "testing"

func TestEqInt64(t *testing.T) {
	if !EqInt64(int64(1), int64(1)) {
		t.Error("not equal")
	}
}


func TestNeInt64(t *testing.T) {
	if !NeInt64(int64(1), int64(2)) {
		t.Error("equal?")
	}
}

func TestParseFilter(t *testing.T) {
	f1 := "abc>20"
	f2 := "cdfg<=12m"
	f3 := "any!=some"
	ft1, err := parseFilter(f1)
	if err != nil {
		t.Errorf("parse filter error: %s", err)
	}
	if ft1.Field != "abc" || ft1.Comparator != GT || ft1.Value != "20" {
		t.Errorf("wrong parse %s", f1)
	}
	ft2, err := parseFilter(f2)
	if err != nil {
		t.Errorf("parse filter error: %s", err)
	}
	if ft2.Field != "cdfg" || ft2.Comparator != LTE || ft2.Value != "12m" {
		t.Errorf("wrong parse %s", f2)
	}
	ft3, err := parseFilter(f3)
	if err != nil {
		t.Errorf("parse filter error: %s", err)
	}
	if ft3.Field != "any" || ft3.Comparator != NE || ft3.Value != "some" {
		t.Errorf("wrong parse %s", f3)
	}
}

func TestParseSize(t *testing.T) {
	var k int64 = 1024
	m := k * 1024
	g := m * 1024
	s1 := "500m"
	s2 := "500M"
	s3 := "623k"
	s4 := "623K"
	s5 := "1G"
	s6 := "1g"
	size, err := parseSize(s1)
	if err != nil {
		t.Errorf("parse size error: %s, err: %s", s1, err)
	}
	if size != 500 * m {
		t.Errorf("size error: %d, expected: %d", size, 500*m)
	}

	size, err = parseSize(s2)
	if err != nil {
		t.Errorf("parse size error: %s, err: %s", s2, err)
	}
	if size != 500 * m {
		t.Errorf("size error: %d, expected: %d", size, 500*m)
	}

	size, err = parseSize(s3)
	if err != nil {
		t.Errorf("parse size error: %s, err: %s", s3, err)
	}
	if size != 623 * k {
		t.Errorf("size error: %d, expected: %d", size, 623*k)
	}

	size, err = parseSize(s4)
	if err != nil {
		t.Errorf("parse size error: %s, err: %s", s4, err)
	}
	if size != 623 * k {
		t.Errorf("size error: %d, expected: %d", size, 623*k)
	}


	size, err = parseSize(s5)
	if err != nil {
		t.Errorf("parse size error: %s, err: %s", s5, err)
	}
	if size != 1 * g {
		t.Errorf("size error: %d, expected: %d", size, 1*g)
	}

	size, err = parseSize(s6)
	if err != nil {
		t.Errorf("parse size error: %s, err: %s", s6, err)
	}
	if size != 1 * g {
		t.Errorf("size error: %d, expected: %d", size, 1*g)
	}
}
