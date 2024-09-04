package nullable

import (
	"encoding/json"
	"fmt"
	"os"
)

func ExampleNullable() {
	type person struct {
		Name    Nullable[string] `json:"name"`
		Age     Nullable[int]    `json:"age"`
		Address Nullable[string] `json:"address"`
	}

	var p person
	if err := json.Unmarshal([]byte(`{"age": null, "address": "Baile Ghrífín"}`), &p); err != nil {
		fmt.Fprintln(os.Stderr, "failed to unmarshal JSON:", err)
		os.Exit(1)
	}
	fmt.Println(p.Name.IsValid(), p.Name.IsNonNull())
	fmt.Println(p.Age.IsValid(), p.Age.IsNonNull())
	fmt.Println(p.Address.IsValid(), p.Address.IsNonNull())
	// Output:
	// false false
	// true false
	// true true
}
