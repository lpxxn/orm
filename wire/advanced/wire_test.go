package advanced

import (
	"fmt"

	alternateinjectorsyntax "github.com/lpxxn/orm/wire/advanced/alternate_injector_syntax"
	argsanderror "github.com/lpxxn/orm/wire/advanced/args_and_error"
	"github.com/lpxxn/orm/wire/advanced/bindinginterfaces"
	"github.com/lpxxn/orm/wire/advanced/bindingstruct"
	"github.com/lpxxn/orm/wire/advanced/bindingvalues"
	"github.com/lpxxn/orm/wire/advanced/cleanupfunctions"
	"github.com/lpxxn/orm/wire/advanced/providersets"
	"github.com/lpxxn/orm/wire/advanced/structfields"
	"github.com/lpxxn/orm/wire/advanced/structproviders"
)

func main() {
	// inject args and return error
	{
		e, err := argsanderror.InitializeEvent("Hello World!")
		if err != nil {
			fmt.Println(err)
		}
		e.Start()
	}

	// Provider Sets
	{
		e, err := providersets.InitializeEvent("Provider Sets")
		if err != nil {
			fmt.Println(err)
		}
		e.Start()
	}

	// Struct Providers
	{
		e, err := structproviders.InitializeEvent("Struct Providers", 1)
		if err != nil {
			fmt.Println(err)
		}
		e.Start()
	}

	// Struct fields
	{
		m := structfields.InitializeMessage("Struct fields", 1)
		fmt.Println(m)
	}

	// Binding Values
	{
		m := bindingvalues.InitializeMessage()
		fmt.Printf("%+v\n", m)
	}

	// Binding Interfaces
	{
		w := bindinginterfaces.InitializeWriter()
		bindinginterfaces.Write(w, "Binding Interfaces")
	}

	// Binding struct to interface
	{
		msg := &bindingstruct.Message{
			Content: "content",
			Code:    1,
		}
		err := bindingstruct.RunStore(msg)
		if err != nil {
			fmt.Println(err)
		}
		err = bindingstruct.WireRunStore(msg)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Cleanup functions
	{
		f, cleanup, err := cleanupfunctions.InitializeFile("testdata/demo.txt")
		if err != nil {
			fmt.Println(err)
		}
		content, err := cleanupfunctions.ReadFile(f)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(content)
		cleanup()
	}

	// Multi cleanup functions
	// {
	// 	app, cleanup, err := multi.InitializeApp()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	defer cleanup()
	// 	n, err := app.Log.Write([]byte("Hello World!"))
	// 	fmt.Println(n, err)
	// }

	// Alternate Injector Syntax
	{
		m := alternateinjectorsyntax.InitializeMessage("Alternate Injector Syntax")
		fmt.Println(m)
	}
}
