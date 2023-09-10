package interval

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestSearchTree_Insert_UpdateValue(t *testing.T) {
	st := NewSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st.root)

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

func TestMultiValueSearchTree_Insert(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st.root)

	vals := []string{"value1", "value2", "value3", "value4"}
	start, end := 17, 19

	err := st.Insert(start, end, vals...)
	if err != nil {
		t.Fatalf("MultiValueSearchTree.Insert: got unexpected error: %v", err)
	}

	got, ok := st.Find(start, end)
	if !ok {
		t.Fatalf("st.Find(%v, %v): got no interval value; want %v", start, end, vals)
	}

	if !reflect.DeepEqual(got, vals) {
		t.Errorf("st.Find(%v, %v): got unexpected value %q; want %q", start, end, got, vals)
	}

	val := "another value"
	st.Insert(start, end, val)

	got, _ = st.Find(start, end)

	if want := append(vals, val); !reflect.DeepEqual(got, want) {
		t.Errorf("st.Find(%v, %v): got unexpected value %v; want %v", start, end, got, want)
	}
}

func TestMultiValueSearchTree_Insert_Error(t *testing.T) {
	t.Run("InvalidRange", func(t *testing.T) {
		st := NewMultiValueSearchTree[int](timeCmp)

		start, end := time.Now(), time.Now().Add(-(1 * time.Hour))
		err := st.Insert(start, end, 0)

		var wantErr InvalidIntervalError
		if !errors.As(err, &wantErr) {
			t.Errorf("st.Insert(%v, %v, 0): got error type %T; want it to be %T", start, end, err, wantErr)
		}
	})

	t.Run("EmptyValueList", func(t *testing.T) {
		st := NewMultiValueSearchTree[int](timeCmp)

		start, end := time.Now(), time.Now().Add(time.Hour)
		err := st.Insert(start, end)

		var wantErr EmptyValueListError
		if !errors.As(err, &wantErr) {
			t.Errorf("st.Insert(%v, %v): got error type %T; want it to be %T", start, end, err, wantErr)
		}
	})
}

func TestMultiValueSearchTree_Upsert(t *testing.T) {
	st := NewMultiValueSearchTree[string](func(x, y int) int { return x - y })
	defer mustBeValidTree(t, st.root)

	vals := []string{"value1", "value2", "value3", "value4"}
	start, end := 17, 19

	err := st.Upsert(start, end, vals...)
	if err != nil {
		t.Fatalf("MultiValueSearchTree.Upsert: got unexpected error: %v", err)
	}

	got, ok := st.Find(start, end)
	if !ok {
		t.Fatalf("st.Find(%v, %v): got no interval value; want %v", start, end, vals)
	}

	if !reflect.DeepEqual(got, vals) {
		t.Errorf("st.Find(%v, %v): got unexpected value %q; want %q", start, end, got, vals)
	}

	val := "another value"
	st.Upsert(start, end, val)

	got, _ = st.Find(start, end)

	if want := []string{val}; !reflect.DeepEqual(got, want) {
		t.Errorf("st.Find(%v, %v): got unexpected value %v; want %v", start, end, got, want)
	}
}

func TestMultiValueSearchTree_Upsert_Error(t *testing.T) {
	t.Run("InvalidRange", func(t *testing.T) {
		st := NewMultiValueSearchTree[int](timeCmp)

		start, end := time.Now(), time.Now().Add(-(1 * time.Hour))
		err := st.Upsert(start, end, 0)

		var wantErr InvalidIntervalError
		if !errors.As(err, &wantErr) {
			t.Errorf("st.Upsert(%v, %v, 0): got error type %T; want it to be %T", start, end, err, wantErr)
		}
	})

	t.Run("EmptyValueList", func(t *testing.T) {
		st := NewMultiValueSearchTree[int](timeCmp)

		start, end := time.Now(), time.Now().Add(time.Hour)
		err := st.Upsert(start, end)

		var wantErr EmptyValueListError
		if !errors.As(err, &wantErr) {
			t.Errorf("st.Upsert(%v, %v): got error type %T; want it to be %T", start, end, err, wantErr)
		}
	})
}
