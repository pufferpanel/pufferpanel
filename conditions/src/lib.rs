use rhai::{Dynamic, Engine};
use wasm_bindgen::prelude::*;

#[wasm_bindgen]
pub fn resolve_if(script: &str) -> Dynamic {
    let engine = Engine::new();
    let ast = engine.compile(script);


    let result = engine.eval::<bool>(script);
    return match result {
        Ok(v) => v,
        Err(e) => {
            e
        }
    };
}