use wasmtime::component::*;
use wasmtime::{Config, Engine, Store, Caller, Module};


bindgen!();

struct MyState {
  name: String,
}

// Imports into the world, like the `name` import for this world, are satisfied
// through traits.
impl HelloWorldImports for MyState {
  // Note the `Result` return value here where `Ok` is returned back to
  // the component and `Err` will raise a trap.
  fn name(&mut self) -> wasmtime::Result<String> {
    Ok(self.name.clone())
  }
}

fn main() -> wasmtime::Result<()> {
  // Configure an `Engine` and compile the `Component` that is being run for
  // the application.
  let mut config = Config::new();
  config.wasm_component_model(true);
  let engine = Engine::new(&config)?;
  let component = Component::from_file(&engine, "./your-component.wasm")?;

  // Instantiation of bindings always happens through a `Linker`.
  // Configuration of the linker is done through a generated `add_to_linker`
  // method on the bindings structure.
  //
  // Note that the closure provided here is a projection from `T` in
  // `Store<T>` to `&mut U` where `U` implements the `HelloWorldImports`
  // trait. In this case the `T`, `MyState`, is stored directly in the
  // structure so no projection is necessary here.
  let mut linker = Linker::new(&engine);
  HelloWorld::add_to_linker(&mut linker, |state: &mut MyState| state)?;

  // As with the core wasm API of Wasmtime instantiation occurs within a
  // `Store`. The bindings structure contains an `instantiate` method which
  // takes the store, component, and linker. This returns the `bindings`
  // structure which is an instance of `HelloWorld` and supports typed access
  // to the exports of the component.
  let mut store = Store::new(
    &engine,
    MyState {
      name: "me".to_string(),
    },
  );
  let (bindings, _) = HelloWorld::instantiate(&mut store, &component, &linker)?;

  // Here our `greet` function doesn't take any parameters for the component,
  // but in the Wasmtime embedding API the first argument is always a `Store`.
  bindings.call_greet(&mut store)?;
  Ok(())
}


// #[tokio::main]
// async fn main() {
//   // let mut config = wasmtime::Config::new();
//   // config.async_support(true).wasm_component_model(true);
//   // let engine = Engine::new(&config)?;
//   // let component = Component::from_file(&engine, &file)?;
//   // let mut linker = Linker::new(&engine);
//   // let mut store = Store::new(&engine, 4);
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//   // Modules can be compiled through either the text or binary format
//   // let engine = Engine::default();
//   // let wat = r#"
//   //       (module
//   //           (import "host" "host_func" (func $host_hello (param i32)))
//   //
//   //           (func (export "hello")
//   //               i32.const 3
//   //               call $host_hello)
//   //       )
//   //   "#;
//   // let module = Module::new(&engine, wat).unwrap();
//   //
//   // // All wasm objects operate within the context of a "store". Each
//   // // `Store` has a type parameter to store host-specific data, which in
//   // // this case we're using `4` for.
//   // let mut store = Store::new(&engine, 4);
//   // let host_func = wasmtime::Func::wrap(&mut store, |caller: Caller<'_, u32>, param: i32| {
//   //   println!("Got {} from WebAssembly", param);
//   //   println!("my host state is: {}", caller.data());
//   // });
//   //
//   // // Instantiation of a module requires specifying its imports and then
//   // // afterwards we can fetch exports by name, as well as asserting the
//   // // type signature of the function with `get_typed_func`.
//   // let instance = wasmtime::Instance::new(&mut store, &module, &[host_func.into()]).unwrap();
//   // let hello = instance.get_typed_func::<(), ()>(&mut store, "hello").unwrap();
//   //
//   // // And finally we can call the wasm!
//   // hello.call(&mut store, ()).unwrap();
//
//
//
//
//
//
//
//
//
//
//
//
//   // Configure an `Engine` and compile the `Component` that is being run for
//   // the application.
//   let mut config = Config::new();
//   config.wasm_component_model(true);
//   let engine = Engine::new(&config)?;
//   let component = Component::from_file(&engine, "./your-component.wasm")?;
//
//   // Instantiation of bindings always happens through a `Linker`.
//   // Configuration of the linker is done through a generated `add_to_linker`
//   // method on the bindings structure.
//   //
//   // Note that the closure provided here is a projection from `T` in
//   // `Store<T>` to `&mut U` where `U` implements the `HelloWorldImports`
//   // trait. In this case the `T`, `MyState`, is stored directly in the
//   // structure so no projection is necessary here.
//   let mut linker = Linker::new(&engine);
//   HelloWorld::add_to_linker(&mut linker, |state: &mut MyState| state)?;
//
//   // As with the core wasm API of Wasmtime instantiation occurs within a
//   // `Store`. The bindings structure contains an `instantiate` method which
//   // takes the store, component, and linker. This returns the `bindings`
//   // structure which is an instance of `HelloWorld` and supports typed access
//   // to the exports of the component.
//   let mut store = Store::new(
//     &engine,
//     MyState {
//       name: "me".to_string(),
//     },
//   );
//   let (bindings, _) = HelloWorld::instantiate(&mut store, &component, &linker)?;
//
//   // Here our `greet` function doesn't take any parameters for the component,
//   // but in the Wasmtime embedding API the first argument is always a `Store`.
//   bindings.call_greet(&mut store)?;
//   Ok(())
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//   //
//   // let result = run(None).await?;
//   // println!("{:?}", result);
//   // Ok(())
// }
