package pokeAPI

// import (
// 	"testing"
// )

// func TestInitLocations(t *testing.T) {
// 	l := InitLocations()
// 	if l == nil {
// 		t.Errorf("InitLocations() returned nil")
// 	} else {
// 		if l.Count < 1 {
// 			t.Errorf("InitLocations() returned locations count < 1")
// 		}
// 	}
// }

// func TestGetLocations(t *testing.T) {
// 	l := GetLocations()
// 	if l == nil {
// 		t.Errorf("GetLocations() returned nil")
// 	} else {
// 		if l.Count < 1 {
// 			t.Errorf("GetLocations() returned locations count < 1")
// 		}
// 	}
// }

// func TestGetNext(t *testing.T) {
// 	l := GetLocations()
// 	if l == nil {
// 		t.Errorf("GetLocations() returned nil")
// 	} else {
// 		if l.Count < 1 {
// 			t.Errorf("GetLocations() returned locations count < 1")
// 		}
// 	}
// 	next := l.Next
// 	if next == nil {
// 		t.Errorf("l.next is nil")
// 	} else {
// 		count := l.Count

// 		l.getNewData(*next)
// 		if count < l.Count {
// 			t.Errorf("GetNext() did not get the next set of locations")
// 		}
// 	}
// }

// func TestGetPrevious(t *testing.T) {
// 	l := GetLocations()
// 	if l == nil {
// 		t.Errorf("GetLocations() returned nil")
// 	} else {
// 		if l.Count < 1 {
// 			t.Errorf("GetLocations() returned locations count < 1")
// 		}
// 	}
// 	previous := l.Previous
// 	if previous == nil {
// 		t.Errorf("l.previous is nil")
// 	} else {
// 		count := l.Count
// 		l.getNewData(*previous)
// 		if count < l.Count {
// 			t.Errorf("GetPrevious() did not get the previous set of locations")
// 		}
// 	}
// }

// func TestPrintLocations(t *testing.T) {
// 	l := GetLocations()
// 	if l == nil {
// 		t.Errorf("GetLocations() returned nil")
// 	} else {
// 		if l.Count < 1 {
// 			t.Errorf("GetLocations() returned locations count < 1")
// 		}
// 	}
// 	err := l.PrintLocations()
// 	if err != nil {
// 		t.Errorf("PrintLocations() returned error: %v", err)
// 	}
// }
