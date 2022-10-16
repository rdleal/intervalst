package interval

import (
	"testing"
	"time"
)

func TestSearchTree_Insert_UpdateValue(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })

	st.Insert(17, 19, "value")
	st.Insert(17, 19, "another value")

	start, end := 17, 19

	want := "another value"

	got, ok := st.Find(start, end)
	if !ok {
		t.Errorf("st.Find(%v, %v): got not interval; want %v", start, end, want)
	}

	if got != want {
		t.Errorf("st.Find(%v, %v): got unexpected value %v; want %v", start, end, got, want)
	}

	if err := st.Delete(start, end); err != nil {
		t.Fatalf("st.Delete(%v, %v): got unexpected error %v", start, end, err)
	}

	got, ok = st.Find(start, end)
	if ok {
		t.Errorf("st.Find(%v, %v): got unexpected value %v", start, end, got)
	}
}

func TestSearchTree_Insert_Error(t *testing.T) {
	t.Run("InvalidRange", func(t *testing.T) {
		st := NewSearchTree[int](timeCmp)

		start, end := time.Now(), time.Now().Add(-(1 * time.Hour))
		err := st.Insert(start, end, 0)
		if err == nil {
			t.Errorf("st.Insert(%v, %v): got nil error", start, end)
		}
	})
}
